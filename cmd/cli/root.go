package cli

import (
	"github.com/spf13/cobra"
	"os"
)

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "pswd is a password manager inspired pass",
}

func Execute() error {
	registerInit(rootCmd)
	registerInsert(rootCmd)
	registerShow(rootCmd)
	registerGenerate(rootCmd)
	registerGenerateKeys(rootCmd)

	args := os.Args[1:]

	cmd, _, err := rootCmd.Find(args)
	if err != nil || !cmd.IsAvailableCommand() || cmd == rootCmd {
		os.Args = append(os.Args[:1], append([]string{"show"}, args...)...)
	}

	return rootCmd.Execute()
}
