module github.com/signavio/aws-mfa-login

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.14.0
	github.com/aws/aws-sdk-go-v2/config v1.14.0
	github.com/aws/aws-sdk-go-v2/credentials v1.9.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.19.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.17.0
	github.com/aws/aws-sdk-go-v2/service/sts v1.15.0
	github.com/aws/smithy-go v1.11.0
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.66.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.3.0
	github.com/spf13/viper v1.10.1
	github.com/stretchr/testify v1.7.0
	k8s.io/client-go v0.23.4
)
