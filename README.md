# aws-mfa-login
Small CLI tool to do aws with mfa and update credentials in local aws config.
It will create or update a destination profile with temporary credentials for `aws_access_key_id`, `aws_secret_access_key` and `aws_session_token`.
Those credentials will be valid for 12 hours by default.

## Getting started
Install executable with golang
```bash
go get -u github.com/signavio/aws-mfa-login
```
or download from releases
```yaml
curl -L https://github.com/signavio/aws-mfa-login/releases/latest/download/aws-mfa-login_linux_amd64.gz -o aws-mfa-login.gz
gunzip aws-mfa-login.gz && chmod +x aws-mfa-login && sudo mv aws-mfa-login /usr/local/bin/aws-mfa-login
```

```console
$ aws-mfa-login -h
CLI tool to update your temporary AWS credentials

Usage:
  aws-mfa-login [flags]

Flags:
      --config string        config file (default is $HOME/.aws-mfa.yaml)
  -d, --destination string   destination profile for temporary aws credentials
  -h, --help                 help for aws-mfa-login
  -s, --source string        source profile where mfa is activated
```
Create application configuration to `~/.aws-mfa.yaml`.
```yaml
source: some-source-profile
destination: some-destination-profile
```
`Source` is source profile where MFA is already activated and the key and secret id is configured.
The tool will create a new profile entry if `destination` profile does not exist yet or update accordingly.
Run tool to update session token in your local aws credentials.

```console
$ aws-mfa-login 
located config file on 
Current Config located in ~/.aws-mfa.yaml
#####
source: suite
destination: suite-mfa

detected MFA device with serial number arn:aws:iam::123456:mfa/username
enter 6-digit MFA code: 123456

Sucessfully update access tokens for profile suite-mfa.
Access will be valid for 11 hours. You can now your profile.

export AWS_PROFILE=suite-mfa
```

## Install from sources

```console
export GO111MODULE=on
go get .
go install .
aws-mfa-login
```