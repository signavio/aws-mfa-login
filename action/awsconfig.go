package action

import (
	"fmt"
	"github.com/ghodss/yaml"
	"github.com/go-ini/ini"
	"github.com/spf13/viper"
	"log"
	"os"
)

type State int

const (
	Updated State = iota + 1
	Created
	Deleted
	Error
)

type Config struct {
	Clusters []ConfigCluster
}

type ConfigCluster struct {
	Name      string `yaml:"name"`
	Alias     string `yaml:"alias"`
	AccountID string `yaml:"accountId"`
	Role      string `yaml:"role"`
}

func GetClusterConfig() []ConfigCluster {
	var conf []ConfigCluster
	err := viper.UnmarshalKey("clusters", &conf)
	if err != nil {
		log.Fatal(err)
	}
	return conf
}

func PrintClusterConfig() {
	clusterConfig := viper.Get("clusters")
	yamlCluster, err := yaml.Marshal(clusterConfig)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	fmt.Printf("%s#####\n", yamlCluster)
	fmt.Printf("%d, %d, %d", Created, Updated, Deleted)
}

func PrintConfigWithoutClusterConfig() {
	fmt.Println("Current Config located in ~/.aws-mfa.yaml\n#####")
	for _, key := range viper.AllKeys() {
		if key != "clusters" {
			fmt.Printf("%v: %v\n", key, viper.Get(key))
		}
	}
	fmt.Print("#####\n")
}

func WriteAll(filePath string) (map[State]int, error) {
	if filePath == "" {
		filePath = getAwsConfigFilePath()
	}
	states := make(map[State]int)
	for _, cluster := range GetClusterConfig() {
		state, err := Write(&cluster, filePath)
		if err != nil {
			log.Fatal(err)
		}
		states[state] += 1
	}
	fmt.Printf("Updated aws credentials in %s\n", filePath)
	fmt.Printf("%d sections updated and %d sections created\n\n", states[Updated], states[Created])
	return states, nil
}

func Write(cluster *ConfigCluster, filePath string) (State, error) {
	file, err := ini.Load(filePath)
	if err != nil {
		return Error, err
	}
	state := Error
	_, err = file.GetSection(cluster.Alias)
	if err == nil {
		state = Updated
	} else {
		state = Created
	}
	section, err := file.NewSection(cluster.Alias)
	if err != nil {
		return Error, err
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", cluster.AccountID, cluster.Role)
	_, err = section.NewKey("role_arn", arn)
	if err != nil {
		return Error, err
	}
	_, err = section.NewKey("source_profile", viper.GetString("destination"))
	if err != nil {
		return Error, err
	}

	return state, file.SaveTo(filePath)
}

func getAwsConfigFilePath() string {
	home, err := os.UserHomeDir()
	path := fmt.Sprintf("%s/.aws/credentials", home)
	if err != nil {
		log.Fatalf("can not located aws config in %s %v", path, err)
		return ""
	}
	return path
}
