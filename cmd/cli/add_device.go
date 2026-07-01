package main

import (
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

func parsePublicKey(token string) (string, [32]byte, error) {

	tokenRaw, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", [32]byte{}, fmt.Errorf("decode token: %w", err)
	}
	if len(tokenRaw) < 34 {
		return "", [32]byte{}, fmt.Errorf("invalid token length: got %d bytes, expected at least 34+", len(tokenRaw))
	}

	var pub [32]byte
	copy(pub[:], tokenRaw)

	var name string
	if len(tokenRaw) > 32 {
		nameLen := binary.BigEndian.Uint16(tokenRaw[32:34])
		if 34+int(nameLen) > len(tokenRaw) {
			return "", [32]byte{}, fmt.Errorf("invalid token: name length %d exceeds token size %d", nameLen, len(tokenRaw))
		}
		name = string(tokenRaw[34 : 34+nameLen])
	}
	return name, pub, nil
}

var addDeviceCmd = &cobra.Command{
	Use:   "add-device <token>",
	Short: "Add a device to the vault",
	Long:  "Add a new device to the vault. Token can be a raw base64-URL-encoded public key (32 bytes) or the full output of pswd export.",
	Args:  cobra.ExactArgs(1),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		name, pub, err := parsePublicKey(args[0])
		if err != nil {
			return err
		}
		if err := ctx.Vault.AddDevice(pub, name, ctx.File, ctx.Priv); err != nil {
			if errors.Is(err, vault.ErrDeviceExists) {
				return fmt.Errorf("device %q is already in the vault", name)
			}
			return fmt.Errorf("add device: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Device", name, "added to vault")
		return nil
	}),
}
