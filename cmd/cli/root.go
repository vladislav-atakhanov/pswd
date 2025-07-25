package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "pswd is a password manager inspired pass",
}

func Execute() error {
	registerInit(rootCmd)
	registerInsert(rootCmd)
	registerShow(rootCmd)
	registerGenerate(rootCmd)
	return rootCmd.Execute()
}
