package main

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:   "search <label>",
	Short: "Search passwords by label",
	Long:  "Fuzzy search password labels in the vault and print matching entries with their IDs.",
	Args:  cobra.ExactArgs(1),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		query := strings.ToLower(args[0])
		for _, entry := range ctx.Vault.List() {
			if strings.Contains(strings.ToLower(entry.Label), query) {
				fmt.Printf("%s - %s\n", entry.Label, hex.EncodeToString(entry.ID[:]))
			}
		}
		return nil
	}),
}
