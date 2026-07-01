package main

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/spf13/cobra"
)

var addDeviceCmd = &cobra.Command{
	Use:   "add-device <token>",
	Short: "Add a device to the vault",
	Long:  "Add a new device to the vault. Token can be a raw base64-URL-encoded public key (32 bytes) or the full output of pswd export.",
	Args:  cobra.ExactArgs(1),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		token := args[0]

		tokenRaw, err := base64.URLEncoding.DecodeString(token)
		if err != nil {
			return fmt.Errorf("decode token: %w", err)
		}
		if len(tokenRaw) < 32 {
			return fmt.Errorf("invalid token length: got %d bytes, expected at least 32", len(tokenRaw))
		}

		var pub [32]byte
		copy(pub[:], tokenRaw)

		var deviceName string
		if len(tokenRaw) > 32 {
			nameLen := binary.BigEndian.Uint16(tokenRaw[32:34])
			deviceName = string(tokenRaw[34 : 34+nameLen])
		} else {
			nameBytes, err := readPassword("Enter device name: ")
			if err != nil {
				return fmt.Errorf("read device name: %w", err)
			}
			deviceName = string(nameBytes)
		}

		if err := ctx.Vault.AddDevice(pub, deviceName, ctx.File, ctx.Priv); err != nil {
			return fmt.Errorf("add device: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Device", deviceName, "added to vault")
		return nil
	}),
}
