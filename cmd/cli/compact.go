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
	RunE: func(cmd *cobra.Command, args []string) error {
		yes, _ := cmd.Flags().GetBool("yes")
		return withVault(func(ctx *vaultContext, _ []string) error {
			if !yes && !confirm("Are you sure you want to rewrite the entire vault? [y/N]: ") {
				fmt.Println("Canceled")
				return nil
			}

			if err := ctx.Vault.Compact(ctx.File, ctx.Priv); err != nil {
				return fmt.Errorf("compact vault: %w", err)
			}
			if err := ctx.Vault.Save(ctx.File); err != nil {
				return fmt.Errorf("save vault: %w", err)
			}

			fmt.Println("Vault compacted")
			return nil
		})(cmd, args)
	},
}

func init() {
	compactCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")
}
