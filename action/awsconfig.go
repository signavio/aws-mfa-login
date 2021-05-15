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

type Clusters struct {
	ClusterConfigs []ClusterConfig
	states         map[State]int
}

type ClusterConfig struct {
	Name      string `yaml:"name"`
	Alias     string `yaml:"alias"`
	AccountID string `yaml:"accountId"`
	Role      string `yaml:"role"`
	Region    string `yaml:"region"`
}

func (clusters *Clusters) InitConfig() {
	var conf []ClusterConfig
	err := viper.UnmarshalKey("clusters", &conf)
	if err != nil {
		log.Fatal(err)
	}
	clusters.ClusterConfigs = conf
}

func (clusters *Clusters) PrintConfig() {
	yamlCluster, err := yaml.Marshal(clusters.ClusterConfigs)
	if err != nil {
		log.Fatalf("unable to marshal config to YAML: %v", err)
	}
	fmt.Printf("%s#####\n", yamlCluster)
}

func (clusters *Clusters) WriteAll(filePath string) error {
	if filePath == "" {
		filePath = getAwsConfigFilePath()
	}
	clusters.states = make(map[State]int)
	for _, cluster := range clusters.ClusterConfigs {
		state, err := cluster.Write(filePath)
		if err != nil {
			return err
		}
		clusters.states[state] += 1
	}
	fmt.Printf("Updated aws credentials in %s\n", filePath)
	fmt.Printf("%d sections updated and %d sections created\n\n", clusters.states[Updated], clusters.states[Created])
	return nil
}

func (c *ClusterConfig) Write(filePath string) (State, error) {
	file, err := ini.Load(filePath)
	if err != nil {
		return Error, err
	}
	state := Error
	_, err = file.GetSection(c.Alias)
	if err == nil {
		state = Updated
	} else {
		state = Created
	}
	section, err := file.NewSection(c.Alias)
	if err != nil {
		return Error, err
	}

	arn := fmt.Sprintf("arn:aws:iam::%s:role/%s", c.AccountID, c.Role)
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

func PrintAwsConfig(filePath string) {
	if filePath == "" {
		filePath = getAwsConfigFilePath()
	}
	PrintFile(filePath)
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

func PrintConfigWithoutClusterConfig() {
	fmt.Println("Current configuration located in ~/.aws-mfa.yaml\n#####")
	for _, key := range viper.AllKeys() {
		if key != "clusters" {
			fmt.Printf("%v: %v\n", key, viper.Get(key))
		}
	}
	fmt.Print("#####\n")
}
