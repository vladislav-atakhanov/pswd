package main

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/uuid"
)

var removeCmd = &cobra.Command{
	Use:   "remove <hex-id>",
	Short: "Remove a password from the vault",
	Args:  cobra.ExactArgs(1),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		raw, err := hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("decode hex id: %w", err)
		}
		if len(raw) != 16 {
			return fmt.Errorf("invalid hex id length: got %d bytes, expected 16", len(raw))
		}
		var id uuid.V4
		copy(id[:], raw)

		label := id.String()
		for _, e := range ctx.Vault.List() {
			if e.ID == id {
				label = e.Label
				break
			}
		}

		if !confirm(fmt.Sprintf("Remove password %q? [y/N]: ", label)) {
			fmt.Println("Canceled")
			return nil
		}

		if err := ctx.Vault.Remove(id); err != nil {
			return fmt.Errorf("remove password: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Password removed")
		return nil
	}),
}
