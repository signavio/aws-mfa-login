module github.com/signavio/aws-mfa-login

go 1.19

require (
	github.com/aws/aws-sdk-go-v2 v1.16.7
	github.com/aws/aws-sdk-go-v2/config v1.15.14
	github.com/aws/aws-sdk-go-v2/credentials v1.12.9
	github.com/aws/aws-sdk-go-v2/service/eks v1.21.4
	github.com/aws/aws-sdk-go-v2/service/iam v1.18.9
	github.com/aws/aws-sdk-go-v2/service/sts v1.16.9
	github.com/aws/smithy-go v1.12.0
	github.com/fatih/color v1.13.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.66.6
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.5.0
	github.com/spf13/viper v1.12.0
	github.com/stretchr/testify v1.8.0
	k8s.io/client-go v0.24.3
)
