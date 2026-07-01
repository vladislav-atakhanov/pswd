package main

import (
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

var removeDeviceCmd = &cobra.Command{
	Use:   "remove-device <hex-pubkey>",
	Short: "Remove a device from the vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		yes, _ := cmd.Flags().GetBool("yes")
		return withVault(func(ctx *vaultContext, args []string) error {
			raw, err := hex.DecodeString(args[0])
			if err != nil {
				return fmt.Errorf("decode public key: %w", err)
			}
			if len(raw) != 32 {
				return fmt.Errorf("invalid public key length: got %d bytes, expected 32", len(raw))
			}
			var pub [32]byte
			copy(pub[:], raw)

			name := hex.EncodeToString(raw)
			for _, d := range ctx.Vault.Devices {
				if d.PublicKey() == pub {
					name = d.Name()
					break
				}
			}

			if !yes && !confirm(fmt.Sprintf("Remove device %q? [y/N]: ", name)) {
				fmt.Println("Canceled")
				return nil
			}

			if err := ctx.Vault.RemoveDevice(pub, ctx.File, ctx.Priv); err != nil {
				if errors.Is(err, vault.ErrDeviceNotFound) {
					return fmt.Errorf("device %s not found in vault", hex.EncodeToString(raw))
				}
				return fmt.Errorf("remove device: %w", err)
			}
			if err := ctx.Vault.Save(ctx.File); err != nil {
				return fmt.Errorf("save vault: %w", err)
			}

			fmt.Println("Device removed from vault")
			return nil
		})(cmd, args)
	},
}

func init() {
	removeDeviceCmd.Flags().BoolP("yes", "y", false, "skip confirmation prompt")
}
