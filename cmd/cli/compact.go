package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var compactCmd = &cobra.Command{
	Use:   "compact",
	Short: "Compact the vault",
	Long:  "Re-encrypt all content in-memory, clear orphaned spans, and rewrite the vault file.",
	Args:  cobra.NoArgs,
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		if err := ctx.Vault.Compact(ctx.File, ctx.Priv); err != nil {
			return fmt.Errorf("compact vault: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Vault compacted")
		return nil
	}),
}
