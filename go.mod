module github.com/signavio/aws-mfa-login

go 1.14

require (
	github.com/aws/aws-sdk-go-v2 v1.6.0
	github.com/aws/aws-sdk-go-v2/config v1.1.7
	github.com/aws/aws-sdk-go-v2/credentials v1.1.7
	github.com/aws/aws-sdk-go-v2/internal/ini v1.0.0 // indirect
	github.com/aws/aws-sdk-go-v2/service/eks v1.3.0
	github.com/aws/aws-sdk-go-v2/service/iam v1.3.2
	github.com/aws/aws-sdk-go-v2/service/sts v1.3.1
	github.com/aws/smithy-go v1.4.0
	github.com/fatih/color v1.12.0
	github.com/ghodss/yaml v1.0.0
	github.com/go-ini/ini v1.62.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/spf13/cobra v1.1.3
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.7.0
	k8s.io/client-go/v12 v12.0.0
)
