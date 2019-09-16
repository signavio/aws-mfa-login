package action

import (
	"bufio"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
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
	iamClient          *iam.IAM
	stsClient          *sts.STS
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
	config := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", updater.sourceProfile),
	}
	session, err := session.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}
	updater.iamClient = iam.New(session)
	updater.stsClient = sts.New(session)
}

func (updater *CredUpdater) getUsername() string {
	callerOutput, err := updater.stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}
	return parseUsername(callerOutput)
}

func (updater *CredUpdater) getMfaSerial(username string) string {
	var mfaId string
	deviceInput := &iam.ListMFADevicesInput{UserName: aws.String(username)}
	devices, err := updater.iamClient.ListMFADevices(deviceInput)
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
	token, err := updater.stsClient.GetSessionToken(tokenInput)
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
	section.NewKey("aws_access_key_id", *token.Credentials.AccessKeyId)
	section.NewKey("aws_secret_access_key", *token.Credentials.SecretAccessKey)
	section.NewKey("aws_session_token", *token.Credentials.SessionToken)

	err = awsFile.SaveTo(awsFilePath)
	if err != nil {
		log.Fatal(err)
	}

	expires := int(token.Credentials.Expiration.Sub(time.Now()).Hours())

	fmt.Printf(`
Sucessfully update access tokens for profile %s.
Access will be valid for %d hours. You can now your profile.

export AWS_PROFILE=%s

`,
		viper.GetString("destination"),
		expires,
		viper.GetString("destination"),
	)
}

func parseUsername(input *sts.GetCallerIdentityOutput) string {
	var arn string
	arn = *input.Arn
	arr := strings.Split(arn, "/")
	if len(arr) == 0 || arr[1] == "" {
		log.Fatalf("Could not detect user name from %v", input)
	}
	return arr[1]
}
