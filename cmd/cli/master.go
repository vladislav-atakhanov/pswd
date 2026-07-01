package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

var masterCmd = &cobra.Command{
	Use:   "master <keyfile>",
	Short: "Change the master password",
	Long:  "Decrypt the private key with the old password, then re-encrypt it with a new password and save to the same file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyfile := args[0]

		data, err := os.ReadFile(keyfile)
		if err != nil {
			return fmt.Errorf("read key file: %w", err)
		}

		oldPassword, err := readPassword("Enter current master password: ")
		if err != nil {
			return fmt.Errorf("read current password: %w", err)
		}
		defer mem.ZeroBytes(oldPassword)

		priv, err := crypto.DecryptPrivateKey(data, oldPassword)
		if err != nil {
			return fmt.Errorf("decrypt private key: %w", err)
		}
		defer mem.ZeroArray32(&priv)

		newPassword, err := readPassword("Enter new master password: ")
		if err != nil {
			return fmt.Errorf("read new password: %w", err)
		}
		defer mem.ZeroBytes(newPassword)

		confirm, err := readPassword("Confirm new master password: ")
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}
		defer mem.ZeroBytes(confirm)

		if len(newPassword) == 0 {
			return fmt.Errorf("password cannot be empty")
		}

		if string(newPassword) != string(confirm) {
			return fmt.Errorf("passwords do not match")
		}

		encrypted, err := crypto.EncryptPrivateKey(priv, newPassword)
		if err != nil {
			return fmt.Errorf("encrypt private key: %w", err)
		}

		if err := os.WriteFile(keyfile, encrypted, 0600); err != nil {
			return fmt.Errorf("write key file: %w", err)
		}

		fmt.Println("Master password updated")
		return nil
	},
}
