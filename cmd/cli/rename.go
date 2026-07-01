package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/uuid"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

var renameCmd = &cobra.Command{
	Use:   "rename <hex-id> <new-label>",
	Short: "Rename a password label",
	Args:  cobra.ExactArgs(2),
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

		if err := ctx.Vault.Rename(id, args[1]); err != nil {
			if errors.Is(err, vault.ErrNotFound) {
				return fmt.Errorf("password %q not found", hex.EncodeToString(raw))
			}
			return fmt.Errorf("rename password: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Password renamed to", args[1])
		return nil
	}),
}
