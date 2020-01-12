package cmd

import (
	"github.com/signavio/aws-mfa-login/action"
	"github.com/spf13/cobra"
	"log"
)

// configureAWSCmd represents the configureAWS command
var configureAWSCmd = &cobra.Command{
	Use:   "configure-aws",
	Short: "configure profiles for aws",
	Long:  "configure profiles for aws by setting profile name, assumed role and source profile",
	Run: func(cmd *cobra.Command, args []string) {
		_, err := action.WriteAll("")
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configureAWSCmd)
}
