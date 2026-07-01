package main

import (
	"encoding/hex"
	"fmt"

	"github.com/spf13/cobra"
)

var removeDeviceCmd = &cobra.Command{
	Use:   "remove-device <hex-pubkey>",
	Short: "Remove a device from the vault",
	Args:  cobra.ExactArgs(1),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		raw, err := hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("decode public key: %w", err)
		}
		if len(raw) != 32 {
			return fmt.Errorf("invalid public key length: got %d bytes, expected 32", len(raw))
		}
		var pub [32]byte
		copy(pub[:], raw)

		if err := ctx.Vault.RemoveDevice(pub, ctx.File, ctx.Priv); err != nil {
			return fmt.Errorf("remove device: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Device removed from vault")
		return nil
	}),
}
