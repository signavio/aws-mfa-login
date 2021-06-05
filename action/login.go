package action

import (
	"bufio"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/go-ini/ini"
	"github.com/spf13/viper"
	"log"
	"os"
	"strings"
	"time"
)

type CredUpdater struct {
	sourceProfile      string
	destinationProfile string
	iamClient          *iam.Client
	stsClient          *sts.Client
}

func UpdateSessionCredentials() {
	updater := &CredUpdater{
		sourceProfile:      viper.GetString("source"),
		destinationProfile: viper.GetString("destination"),
	}
	updater.init()
	username := updater.getUsername()
	serial := updater.getMfaSerial(username)
	code := updater.readCode()
	token := updater.getSessionToken(serial, code)
	updater.updateAwsConfig(token)
}

func (updater *CredUpdater) init() {
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithSharedConfigProfile(updater.sourceProfile),
		config.WithRegion("eu-central-1"),
	)
	if err != nil {
		log.Fatal(err)
	}
	updater.iamClient = iam.NewFromConfig(cfg)
	updater.stsClient = sts.NewFromConfig(cfg)
}

func (updater *CredUpdater) getUsername() string {
	callerOutput, err := updater.stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}
	username, err := ParseUsername(callerOutput)
	if err != nil {
		log.Fatalf(err.Error())
	}
	return username
}

func (updater *CredUpdater) getMfaSerial(username string) string {
	var mfaId string
	deviceInput := &iam.ListMFADevicesInput{UserName: aws.String(username)}
	devices, err := updater.iamClient.ListMFADevices(context.TODO(), deviceInput)
	if err != nil {
		log.Fatal(err)
	}
	if len(devices.MFADevices) == 0 {
		log.Fatal("user has no MFA Device activated")
	} else {
		mfaId = *devices.MFADevices[0].SerialNumber
		fmt.Printf("detected MFA device with serial number %s\n", mfaId)
	}
	return mfaId
}

func (updater *CredUpdater) readCode() string {
	fmt.Print("enter 6-digit MFA code: ")
	userInput := bufio.NewReader(os.Stdin)
	in, _, err := userInput.ReadLine()
	if err != nil {
		log.Fatal(err)
	}
	return string(in)
}

func (updater *CredUpdater) getSessionToken(serial string, code string) *sts.GetSessionTokenOutput {
	tokenInput := &sts.GetSessionTokenInput{
		DurationSeconds: nil,
		SerialNumber:    &serial,
		TokenCode:       &code,
	}
	token, err := updater.stsClient.GetSessionToken(context.TODO(), tokenInput)
	if err != nil {
		log.Fatal(err)
	}
	return token
}

func (updater *CredUpdater) updateAwsConfig(token *sts.GetSessionTokenOutput) {
	home, err := os.UserHomeDir()
	awsFilePath := fmt.Sprintf("%s/.aws/credentials", home)
	if err != nil {
		log.Fatal(err)
	}
	awsFile, err := ini.Load(awsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	section, err := awsFile.NewSection(viper.GetString("destination"))
	if err != nil {
		log.Fatal(err)
	}
	_, _ = section.NewKey("aws_access_key_id", *token.Credentials.AccessKeyId)
	_, _ = section.NewKey("aws_secret_access_key", *token.Credentials.SecretAccessKey)
	_, _ = section.NewKey("aws_session_token", *token.Credentials.SessionToken)

	err = awsFile.SaveTo(awsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	expires := int(token.Credentials.Expiration.Sub(time.Now()).Hours())

	fmt.Printf(`
Successfully updated access tokens for profile %s.
Access will be valid for %d hours. You can now use that profile.

export AWS_PROFILE=%s

`,
		viper.GetString("destination"),
		expires,
		viper.GetString("destination"),
	)
}

func ParseUsername(input *sts.GetCallerIdentityOutput) (string, error) {
	var arn string
	arn = *input.Arn
	if !strings.Contains(arn, "/") {
		return "", &ArnParseException{"arn does not have expected format arn:aws:iam::123456:user/someuser"}
	}
	arr := strings.Split(arn, "/")
	if len(arr) == 0 || arr[1] == "" {
		msg := fmt.Sprintf("Could not detect user name from %v", input)
		return "", &ArnParseException{msg}
	}
	return arr[1], nil
}
