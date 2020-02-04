# aws-mfa-login
Small CLI tool to do aws with mfa and update credentials in local aws config.
It will create or update a destination profile with temporary credentials for `aws_access_key_id`, `aws_secret_access_key` and `aws_session_token`.
Those credentials will be valid for 12 hours by default.

## Getting started

### Install using go...

For this, go must be installed on your system. 

Install executable with golang
```bash
go get github.com/signavio/aws-mfa-login
```
.. or update by setting the -u flag: 
```bash
go get -u github.com/signavio/aws-mfa-login
``` 

Make sure your go path is part of your PATH environment variable: 
```
export GOPATH="~/go"
export PATH="${PATH}:${GOPATH}/bin/"
```

### .. or download from releases

```console
curl -L https://github.com/signavio/aws-mfa-login/releases/latest/download/aws-mfa-login_$(uname)_amd64.gz -o aws-mfa-login.gz
gunzip aws-mfa-login.gz && chmod +x aws-mfa-login && sudo mv aws-mfa-login /usr/local/bin/aws-mfa-login
```

### Post-install

Check your installation - this should work now: 

```console
$ aws-mfa-login -h
CLI tool to update your temporary AWS credentials

Usage:
  aws-mfa-login [flags]
  aws-mfa-login [command]

Available Commands:
  aws         setup or view your aws config
  cluster     view or setup your kubeconfig
  completion  Generates bash completion scripts
  help        Help about any command

Flags:
      --config string        config file (default is $HOME/.aws-mfa.yaml)
  -d, --destination string   destination profile for temporary aws credentials
  -h, --help                 help for aws-mfa-login
  -s, --source string        source profile where mfa is activated
      --version              version for aws-mfa-login
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

## Setup kubernetes access
you can provide information for static clusters in the yaml see example:
```yaml
source: suite
destination: suite-mfa
clusters:
    - name: eks-staging
      alias: suite-staging
      accountId: "1234"
      role: DeveloperAccessRole
      region: eu-central-1
    - name: eks-prod
      alias: suite-academic
      accountId: "4321"
      role: DeveloperAccessRole
      region: eu-central-1
```
Then you can setup the assumed roles in your aws config and also update the kubeconfig to access the cluster.
```bash
aws-mfa-login aws setup
> Updated aws credentials in ~/.aws/credentials
> 2 sections updated and 0 sections created

aws-mfa-login cluster setup
> Updated context suite-staging in C:\Users\Karl\.kube\config

> you can switch to cluster e.g. with:
> kubectl config use-context suite-staging
```

## Autocompletion

```console
aws-mfa-login completion -h
> To enable autocompletion one-time run
> source <(aws-mfa-login completion)
> To enable autocompletion for all terminal sessions add this your bashrc
> # ~/.bashrc or ~/.profile
> source <(aws-mfa-login completion)
```

# Development

## Versioning
In order to increase version your commit message (or squash merge) should start with `major:`, `minor:` or `patch:`.
See https://github.com/stevenmatthewt/semantics#how-it-works
The CI will publish artifacts to releases page and increment version.
Also increase `VERSION` in [root.go](cmd/root.go) matching to your increment string.

## Install from sources

```console
export GO111MODULE=on
go build .
go install .
aws-mfa-login
```