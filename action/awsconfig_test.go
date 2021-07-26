package action

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/go-ini/ini"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	conf := []ClusterConfig{
		{
			Name:      "testname",
			Alias:     "testalias",
			AccountID: "123456",
			Role:      "dev",
		},
		{
			Name:      "testname2",
			Alias:     "testalias2",
			AccountID: "123456",
			Role:      "admin",
		},
		{
			Name:        "testname3",
			Alias:       "testalias2",
			AccountID:   "123456",
			Role:        "admin",
			Destination: "altProfile",
		},
	}
	viper.Set("clusters", conf)
	viper.Set("destination", "test-mfa")

	clusters := &Clusters{}
	clusters.InitConfig()
	clusters.PrintConfig()

	file, err := ioutil.TempFile(os.TempDir(), "aws-test")
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			log.Fatal(err)
		}
		err = os.Remove(file.Name())
		if err != nil {
			log.Fatal(err)
		}
	}()

	err = clusters.WriteAll(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	result, err := ini.Load(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	PrintFile(file.Name())

	foundSection, err := result.GetSection(conf[0].Alias)
	if err != nil {
		log.Fatal(err)
	}
	foundRole, err := foundSection.GetKey("role_arn")
	if err != nil {
		log.Fatal(err)
	}
	foundProfile, err := foundSection.GetKey("source_profile")
	if err != nil {
		log.Fatal(err)
	}

	// one section DEFAULT will be added by default
	assert.ElementsMatch(t, result.SectionStrings(), []string{"DEFAULT", conf[0].Alias, conf[1].Alias})
	assert.Equal(t, foundRole.Value(), getArn(conf[0].AccountID, conf[0].Role))
	assert.Equal(t, foundProfile.Value(), viper.GetString("destination"))
	assert.Equal(t, clusters.states[Created], 2)
	assert.Equal(t, clusters.states[Updated], 0)
	assert.Equal(t, clusters.states[Deleted], 0)

	modified := &ClusterConfig{
		Name:      "changed",
		Alias:     conf[1].Alias,
		AccountID: "123",
		Role:      "changed",
	}

	// we modify a section. the number of sections should not change but the content of the section.

	state, err := modified.Write(file.Name())
	if err != nil {
		log.Fatal(err)
	}
	result, err = ini.Load(file.Name())
	if err != nil {
		log.Fatal(err)
	}

	foundSection, err = result.GetSection(conf[1].Alias)
	if err != nil {
		log.Fatal(err)
	}
	foundRole, err = foundSection.GetKey("role_arn")
	if err != nil {
		log.Fatal(err)
	}
	foundProfile, err = foundSection.GetKey("source_profile")
	if err != nil {
		log.Fatal(err)
	}

	assert.Len(t, result.SectionStrings(), 3)
	assert.Equal(t, foundRole.Value(), getArn(modified.AccountID, modified.Role))
	assert.Equal(t, foundProfile.Value(), viper.GetString("destination"))
	assert.Equal(t, state, Updated)

	PrintFile(file.Name())
}

func getArn(accountId string, role string) string {
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", accountId, role)
}
