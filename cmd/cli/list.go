package main

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all passwords",
	Args:  cobra.NoArgs,
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		for _, entry := range ctx.Vault.List() {
			t := time.Unix(int64(entry.LastUpdate), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("%s | %s | %s\n", entry.Label, hex.EncodeToString(entry.ID[:]), t)
		}
		return nil
	}),
}
