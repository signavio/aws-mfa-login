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

func UpdateSessionCredentials(sourceProfile string) {
	PrintConfig()
	config := &aws.Config{
		Credentials: credentials.NewSharedCredentials("", sourceProfile),
	}
	session, err := session.NewSession(config)
	if err != nil {
		log.Fatal(err)
	}
	iamClient := iam.New(session)
	stsClient := sts.New(session)

	callerOutput, err := stsClient.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatal(err)
	}
	username := parseUsername(callerOutput)

	var mfaId string
	deviceInput := &iam.ListMFADevicesInput{UserName: aws.String(username)}
	devices, err := iamClient.ListMFADevices(deviceInput)

	if err != nil {
		log.Fatal(err)
	}
	if len(devices.MFADevices) == 0 {
		log.Fatal("user has no MFA Device activated", )
	} else {
		mfaId = *devices.MFADevices[0].SerialNumber
		fmt.Printf("detected MFA device with serial number %s\n", mfaId)
	}

	// get mfa code from user
	fmt.Print("enter 6-digit MFA code: ")
	userInput := bufio.NewReader(os.Stdin)
	in, _, err := userInput.ReadLine()
	code := string(in)

	// get session token
	tokenInput := &sts.GetSessionTokenInput{
		DurationSeconds: nil,
		SerialNumber:    &mfaId,
		TokenCode:       &code,
	}
	token, err := stsClient.GetSessionToken(tokenInput)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(token)

	// update aws config
	home, err := os.UserHomeDir()
	awsFilePath := fmt.Sprintf("%s/.aws/credentials", home)
	if err != nil {
		log.Fatal(err)
	}
	awsFile, err := ini.Load(awsFilePath)
	if err != nil {
		log.Fatal(err)
	}
	section, err := awsFile.NewSection(viper.GetString("profile"))
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
		viper.GetString("profile"),
		expires,
		viper.GetString("profile"),
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

func PrintConfig() {
	fmt.Println("Current Config\n#####")
	for _, key := range viper.AllKeys() {
		fmt.Printf("%v: %v\n", key, viper.Get(key))
	}
	fmt.Print("\n")
}
