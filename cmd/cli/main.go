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

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "Password manager",
	Long:  "A CLI password manager with X25519 encryption.",
}

func main() {
	rootCmd.AddCommand(genkeyCmd, exportCmd, masterCmd)
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
