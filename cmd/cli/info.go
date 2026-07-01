package main

import (
	"os"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info -p <vault> -k <key>",
	Short: "Display vault contents",
	Args:  cobra.NoArgs,
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		return ctx.Vault.Print(os.Stdout)
	}),
}
