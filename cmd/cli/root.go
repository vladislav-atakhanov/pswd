package cli

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "pswd is CLI interface for PSWD",
}

func Execute() error {
	registerInit(rootCmd)
	registerInsert(rootCmd)
	registerShow(rootCmd)
	return rootCmd.Execute()
}
