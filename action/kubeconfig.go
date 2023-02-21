package action

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/aws/smithy-go"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

type KubeConfigUpdater struct {
	Profile   string
	Clusters  *Clusters
	stsClient *sts.Client
	awsConfig aws.Config
}

func (updater *KubeConfigUpdater) Init() {

	clusters := &Clusters{}
	clusters.InitConfig()
	updater.Clusters = clusters
	updater.Profile = viper.GetString("destination")

	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(updater.Profile),
		config.WithRegion("eu-central-1"),
	)
	updater.awsConfig = cfg
	if err != nil {
		log.Fatal(err)
	}
	updater.stsClient = sts.NewFromConfig(cfg)
}

func (c *ClusterConfig) assumeRole(stsClient *sts.Client) *eks.Client {
	roleArn := fmt.Sprintf("arn:aws:iam::%s:role/%s", c.AccountID, c.Role)
	appCreds := stscreds.NewAssumeRoleProvider(stsClient, roleArn)
	return eks.New(eks.Options{Credentials: appCreds, Region: c.Region})
}

func (c *ClusterConfig) getCluster(stsClient *sts.Client) (*eks.DescribeClusterOutput, error) {
	eksClient := c.assumeRole(stsClient)
	return eksClient.DescribeCluster(context.TODO(), &eks.DescribeClusterInput{
		Name: aws.String(c.Name),
	})
}

func (c *ClusterConfig) List(stsClient *sts.Client) {
	eksClient := c.assumeRole(stsClient)
	out, err := eksClient.ListClusters(context.TODO(), &eks.ListClustersInput{})
	if err != nil {
		fmt.Printf("Could not list clusters in region %s in account %s\n", c.Region, c.AccountID)
		return
	}
	fmt.Printf("Region: %s\n", c.Region)
	fmt.Printf("Clusters:\n")
	for _, cluster := range out.Clusters {
		fmt.Println(" - ", cluster)
	}
}

func (c *ClusterConfig) writeKubeconfig(clusterInfo *eks.DescribeClusterOutput) {
	kubeConfigPath, err := findKubeConfig()
	if err != nil {
		PrintError("Failed to detect kubeconfig %s\n", err.Error())
		return
	}
	err = createFileIfNotExist(kubeConfigPath)
	if err != nil {
		PrintError("Could not create initial kubeconfig file %s\n", err.Error())
		return
	}
	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)
	if err != nil {
		PrintError("Failed to unmarshall existing kubeconfig %s\n", err.Error())
		return
	}

	clusterArn := aws.ToString(clusterInfo.Cluster.Arn)
	// certificate is already base64 encoded but we a decoded string
	caData := base64Decode(aws.ToString(clusterInfo.Cluster.CertificateAuthority.Data))

	clusterConfig := clientcmdapi.NewCluster()
	clusterConfig.Server = aws.ToString(clusterInfo.Cluster.Endpoint)
	clusterConfig.CertificateAuthorityData = caData

	authConfig := clientcmdapi.NewAuthInfo()
	authConfig.Exec = &clientcmdapi.ExecConfig{
		Command: "aws",
		Args: []string{
			"--region",
			c.Region,
			"eks",
			"get-token",
			"--cluster-name",
			aws.ToString(clusterInfo.Cluster.Name),
		},
		Env: []clientcmdapi.ExecEnvVar{
			{
				Name:  "AWS_PROFILE",
				Value: c.Alias,
			},
		},
		APIVersion: "client.authentication.k8s.io/v1",
	}

	contextConfig := clientcmdapi.NewContext()
	contextConfig.Cluster = aws.ToString(clusterInfo.Cluster.Arn)
	contextConfig.AuthInfo = c.Alias

	kubeConfig.Clusters[clusterArn] = clusterConfig
	kubeConfig.AuthInfos[c.Alias] = authConfig
	kubeConfig.Contexts[c.Alias] = contextConfig

	configAccess := &clientcmd.ClientConfigLoadingRules{}
	configAccess.ExplicitPath = kubeConfigPath
	err = clientcmd.ModifyConfig(configAccess, *kubeConfig, true)
	if err != nil {
		PrintError("Failed to write kubeconfig %s\n", err.Error())
		return
	}
}

func (updater *KubeConfigUpdater) SetupClusters() {

	for _, cluster := range updater.Clusters.ClusterConfigs {
		clusterInfo, err := cluster.getCluster(updater.stsClient)
		if err != nil {
			var notFound *eksTypes.ResourceNotFoundException
			var accessDenied *smithy.OperationError
			if errors.As(err, &notFound) {
				PrintWarn("Skipping setup for cluster %s %s\n", cluster.Name, aws.ToString(notFound.Message))
				continue
			}
			if errors.As(err, &accessDenied) {
				PrintWarn("Skipping setup for cluster %s because not authorized\n", cluster.Name)
				continue
			}
			PrintError("Skipping setup for cluster %s %s\n", cluster.Name, err.Error())
			continue
		}
		cluster.writeKubeconfig(clusterInfo)
		PrintSuccess("Successfully setup kubeconfig for cluster %s\n", cluster.Name)

	}
	if len(updater.Clusters.ClusterConfigs) >= 1 {
		fmt.Printf("\nYou can now switch the context. Kubectx is recommenend, "+
			"although it's also possible with plain kubectl.\n"+
			"kubectx %s\nkubectl config use-context %s\n",
			updater.Clusters.ClusterConfigs[0].Alias,
			updater.Clusters.ClusterConfigs[0].Alias)
	}
}

func (updater *KubeConfigUpdater) ListClusters() {
	for _, cluster := range updater.Clusters.ClusterConfigs {
		cluster.List(updater.stsClient)
	}
}

func findKubeConfig() (string, error) {
	env := os.Getenv("KUBECONFIG")
	if env != "" {
		return env, nil
	}
	path, err := homedir.Expand("~/.kube/config")
	if err != nil {
		return "", err
	}
	return path, nil
}
