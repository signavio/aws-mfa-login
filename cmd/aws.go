package cmd

import (
	"fmt"
	"github.com/signavio/aws-mfa-login/action"
	"github.com/spf13/cobra"
	"log"
)

var awsCommand = &cobra.Command{
	Use:   "aws [setup|view]",
	Short: "setup or view your aws config",
	Long:  "setup or view your aws config",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("use aws setup or aws view")
	},
}

// configureAWSCmd represents the configureAWS command
var configureAWSCmd = &cobra.Command{
	Use:   "setup",
	Short: "setup profiles for aws",
	Long:  "setup profiles for aws by setting profile name, assumed role and source profile",
	Run: func(cmd *cobra.Command, args []string) {
		clusters := &action.Clusters{}
		clusters.InitConfig()
		err := clusters.WriteAll("")
		if err != nil {
			log.Fatal(err)
		}
	},
}

var viewAwsConfig = &cobra.Command{
	Use:   "view",
	Short: "view your current aws config",
	Long:  "view your current aws config",
	Run: func(cmd *cobra.Command, args []string) {
		action.PrintAwsConfig("")
	},
}

func init() {
	rootCmd.AddCommand(awsCommand)
	awsCommand.AddCommand(configureAWSCmd, viewAwsConfig)
}
