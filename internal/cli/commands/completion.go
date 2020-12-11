package commands

import (
	"github.com/spf13/cobra"
	"os"
)

var CompletionCmd = &cobra.Command{
	Use:                   "completion [bash|zsh|fish|powershell]",
	Short:                 "Generate completion script",
	Long: `To load completions:

Bash:

$ source <(bgc completion bash)

# To load completions for each session, execute once:
Linux:
  $ bgc completion bash > /etc/bash_completion.d/bgc
MacOS:
  $ bgc completion bash > /usr/local/etc/bash_completion.d/bgc

Zsh:

# If shell completion is not already enabled in your environment you will need
# to enable it.  You can execute the following once:

$ echo "autoload -U compinit; compinit" >> ~/.zshrc

# To load completions for each session, execute once:
$ bgc completion zsh > "${fpath[1]}/_bgc"

# You will need to start a new shell for this setup to take effect.

Fish:

$ bgc completion fish | source

# To load completions for each session, execute once:
$ bgc completion fish > ~/.config/fish/completions/bgc.fish
`,
	DisableFlagsInUseLine: true,
	Hidden: true,
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	Args:                  cobra.ExactValidArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "bash":
			cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
	},
}
