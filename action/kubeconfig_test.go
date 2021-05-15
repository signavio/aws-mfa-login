package action

import (
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	eksTypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"k8s.io/client-go/tools/clientcmd"
	"log"
	"os"
	"testing"
)

func TestUpdate(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping testing in short mode")
	}
	initViper()
	clusters := &Clusters{}
	clusters.InitConfig()
	updater := &KubeConfigUpdater{
		Profile:  "mfa",
		Clusters: clusters,
	}
	_ = os.Setenv("KUBECONFIG", "/tmp/kubeconfig/kubeconfig")
	updater.Init()
	updater.SetupClusters()
	updater.ListClusters()
}

func TestWriteKubconfig(t *testing.T) {
	_ = os.Setenv("KUBECONFIG", "/tmp/kubeconfig/kubeconfig-test")
	clusters := []ClusterConfig{
		{
			Name:      "cluster1",
			Alias:     "alias1",
			AccountID: "12345678",
			Role:      "Role1",
			Region:    "eu-central-1",
		},
		{
			Name:      "cluster1",
			Alias:     "alias2",
			AccountID: "12345678",
			Role:      "Role2",
			Region:    "eu-central-1",
		},
	}
	for _, cluster := range clusters {
		clusterOutput := &eks.DescribeClusterOutput{
			Cluster: &eksTypes.Cluster{
				Arn: aws.String(fmt.Sprintf("arn:aws:eks:%s:%s:cluster/%s", cluster.Region, cluster.AccountID, cluster.Name)),
				CertificateAuthority: &eksTypes.Certificate{
					// cert-data
					Data: aws.String("Y2VydC1kYXRhCg=="),
				},
				Name:     aws.String(cluster.Name),
				Endpoint: aws.String("https://123456789.gr7.eu-central-1.eks.amazonaws.com"),
			},
		}
		cluster.writeKubeconfig(clusterOutput)
	}

	kubeConfigPath, _ := findKubeConfig()
	kubeConfig, err := clientcmd.LoadFromFile(kubeConfigPath)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(kubeConfig.Clusters), "Expecting only one cluster")
	assert.Equal(t, 2, len(kubeConfig.AuthInfos), "Expecting two different users")
	assert.Equal(t, 2, len(kubeConfig.Contexts), "Expecting two different contexts")

}

func initViper() {
	home, err := homedir.Dir()
	if err != nil {
		log.Fatal(err)
	}
	viper.AddConfigPath(home)
	viper.SetConfigName(".aws-mfa")
	viper.AutomaticEnv()
	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
}
