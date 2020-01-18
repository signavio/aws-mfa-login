package cmd

import (
	"github.com/signavio/aws-mfa-login/action"
	"github.com/spf13/cobra"
)

var clusterCommand = &cobra.Command{
	Use:   "cluster [setup|view]",
	Short: "view or setup your kubeconfig",
	Long:  "view or setup your kubeconfig",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var clusterSetupCommand = &cobra.Command{
	Use:   "setup",
	Short: "setup your kubeconfig",
	Long:  "setup your kubeconfig",
	Run: func(cmd *cobra.Command, args []string) {
		action.SetupClusters()
	},
}

var clusterViewCommand = &cobra.Command{
	Use:   "view",
	Short: "view all cluster names and metadata",
	Long:  "view all cluster names and metadata",
	Run: func(cmd *cobra.Command, args []string) {
		action.ListClusters()
	},
}

func init() {
	rootCmd.AddCommand(clusterCommand)
	clusterCommand.AddCommand(clusterSetupCommand, clusterViewCommand)
}
