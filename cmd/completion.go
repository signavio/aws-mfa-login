package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates bash completion scripts",
	Long: `To enable autocompletion one-time run

source <(aws-mfa-login completion)

To enable autocompletion for all terminal sessions add this your bashrc

# ~/.bashrc or ~/.profile
source <(aws-mfa-login completion)
`,
	Run: func(cmd *cobra.Command, args []string) {
		rootCmd.GenBashCompletion(os.Stdout)
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
