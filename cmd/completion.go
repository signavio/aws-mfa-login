package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate completion script",
	Long: `To load completions:

Bash:

  $ source <(yourprogram completion bash)

  # To load completions for each session, execute once:
  # Linux:
  $ aws-mfa-login completion bash > /etc/bash_completion.d/aws-mfa-login
  # macOS:
  $ aws-mfa-login completion bash > /usr/local/etc/bash_completion.d/aws-mfa-login

Zsh:

  # If shell completion is not already enabled in your environment,
  # you will need to enable it.  You can execute the following once:

  $ echo "autoload -U compinit; compinit" >> ~/.zshrc

  # To load completions for each session, execute once:
  $ aws-mfa-login completion zsh > ~/.oh-my-zsh/completions/_aws-mfa-login
  
  # verify that ~/.oh-my-zsh/completions is in your fpath
  $ print -l $fpath

  # You will need to start a new shell for this setup to take effect.

fish:

  $ aws-mfa-login completion fish | source

  # To load completions for each session, execute once:
  $ aws-mfa-login completion fish > ~/.config/fish/completions/aws-mfa-login.fish

PowerShell:

  PS> aws-mfa-login completion powershell | Out-String | Invoke-Expression

  # To load completions for every new session, run:
  PS> aws-mfa-login completion powershell > aws-mfa-login.ps1
  # and source this file from your PowerShell profile.
`,
	DisableFlagsInUseLine: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			_ = cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			_ = cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			_ = cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			_ = cmd.Root().GenPowerShellCompletionWithDesc(os.Stdout)
		}
	},
}

func init() {
	rootCmd.AddCommand(completionCmd)
}
