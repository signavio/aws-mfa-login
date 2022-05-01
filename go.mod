module github.com/signavio/aws-mfa-login

go 1.16

require (
	github.com/aws/aws-sdk-go-v2 v1.16.3
	github.com/aws/aws-sdk-go-v2/config v1.15.4
	github.com/aws/aws-sdk-go-v2/credentials v1.12.0
	github.com/aws/aws-sdk-go-v2/service/eks v1.20.6
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.4
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.4
	github.com/aws/smithy-go v1.11.2
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.66.4
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.4.0
	github.com/spf13/viper v1.11.0
	github.com/stretchr/testify v1.7.1
	k8s.io/client-go v0.23.6
)
