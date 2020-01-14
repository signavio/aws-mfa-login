/**
I dont like to use system calls.
But to use whole API, you would need to reimplement complete package to update kubeconfig file
https://github.com/aws/aws-cli/tree/develop/awscli/customizations/eks
We can refactor this later with API calls, for now we rely on aws cli implementation.
*/
package action

import (
	"fmt"
	"github.com/blang/semver"
	"log"
	"os/exec"
	"regexp"
)

const REQUIRED_MIN_AWS_VERSION = "1.16.308"

func runCommand(command string, args []string) string {
	cmd := exec.Command(command, args...)
	stdout, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatalf("%s%s", err, stdout)
	}
	fmt.Printf("%s", stdout)
	return string(stdout)
}

func ListClusters() {
	for _, cluster := range GetClusterConfig() {
		fmt.Println("####")
		fmt.Printf("Cluster: %s\nRegion: %s\n", cluster.Alias, cluster.Region)
		runCommand("aws", []string{
			"eks",
			"list-clusters",
			"--region", cluster.Region,
			"--profile", cluster.Alias,
		})
		fmt.Printf("####\n\n")
	}
}

func SetupClusters() {
	for _, cluster := range GetClusterConfig() {
		runCommand("aws", []string{
			"eks",
			"update-kubeconfig",
			"--region", cluster.Region,
			"--profile", cluster.Alias,
			"--alias", cluster.Alias,
			"--name", cluster.Name,
		})
	}
}

func CheckRequiredAwsVersion(versionString string) (bool, error) {
	version, err := semver.Make(versionString)
	if err != nil {
		//log.Fatalf("could not parse incoming version: %s", versionString)
		return false, err
	}
	minmumVersion, _ := semver.Make(REQUIRED_MIN_AWS_VERSION)
	if version.Compare(minmumVersion) < 0 {
		return false, fmt.Errorf("aws cli version must be greater than %s and you have %s\n\n", REQUIRED_MIN_AWS_VERSION, versionString)
	}
	fmt.Printf("Version is greater than %s\n", REQUIRED_MIN_AWS_VERSION)
	return true, nil
}

//func getAwsVersion() string {
//	version := runCommand("aws", []string{"--version"})
//	parsedVersion, err := parseAwsVersion(version)
//	if err != nil {
//		log.Fatal(parsedVersion)
//	}
//	return parsedVersion
//}

func parseAwsVersion(version string) (string, error) {
	regex := regexp.MustCompile(`^aws-cli/([\d]+\.[\d]+\.[\d]+).*`)
	if !regex.MatchString(version) {
		return "", &AwsVersionParseException{"Could not parse aws version from: " + version}
	}
	return regex.FindStringSubmatch(version)[1], nil
}
