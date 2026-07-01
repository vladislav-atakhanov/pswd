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
		for id, item := range ctx.Vault.Content {
			if strings.Contains(strings.ToLower(item.Label), query) {
				fmt.Printf("%s - %s\n", item.Label, hex.EncodeToString(id[:]))
			}
		}
		return nil
	}),
}
