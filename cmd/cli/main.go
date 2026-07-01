package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"golang.org/x/term"
)

func readPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return password, err
}

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "Password manager",
	Long:  "A CLI password manager with X25519 encryption.",
	Run: func(cmd *cobra.Command, _ []string) {
		cmd.Help()
	},
}

func init() {
	rootCmd.PersistentFlags().StringP("vault", "p", "", "path to vault file")
	rootCmd.PersistentFlags().StringP("key", "k", "", "path to private key file")

	rootCmd.AddGroup(
		&cobra.Group{ID: "password", Title: "Password Commands:"},
		&cobra.Group{ID: "device", Title: "Device Commands:"},
		&cobra.Group{ID: "vault", Title: "Vault Commands:"},
		&cobra.Group{ID: "utility", Title: "Utility:"},
	)
}

func main() {
	rootCmd.AddCommand(
		addCmd, showCmd, editCmd, removeCmd, renameCmd, listCmd,
		searchCmd, generateCmd,
		deviceCmd,
		initCmd, compactCmd, statsCmd, passwdCmd,
		genkeyCmd, completionCmd,
	)

	addCmd.GroupID = "password"
	showCmd.GroupID = "password"
	editCmd.GroupID = "password"
	removeCmd.GroupID = "password"
	renameCmd.GroupID = "password"
	listCmd.GroupID = "password"
	searchCmd.GroupID = "password"
	generateCmd.GroupID = "password"

	deviceCmd.GroupID = "device"

	initCmd.GroupID = "vault"
	compactCmd.GroupID = "vault"
	statsCmd.GroupID = "vault"
	passwdCmd.GroupID = "vault"

	genkeyCmd.GroupID = "utility"
	completionCmd.GroupID = "utility"

	ctx := trapSignals()
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
