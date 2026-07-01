package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

var genkeyCmd = &cobra.Command{
	Use:   "genkey <file>",
	Short: "Generate a new encryption key pair",
	Long:  "Generate a new X25519 key pair, encrypt the private key with a master password, and save it to a file.",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filepath := args[0]

		priv, _, err := crypto.GenerateKeys()
		if err != nil {
			return fmt.Errorf("generate keys: %w", err)
		}

		password, err := readPassword("Enter master password: ")
		if err != nil {
			return fmt.Errorf("read password: %w", err)
		}
		mem.Lock(password)
		defer mem.Unlock(password)
		defer mem.ZeroBytes(password)

		confirm, err := readPassword("Confirm master password: ")
		if err != nil {
			return fmt.Errorf("read confirmation: %w", err)
		}
		mem.Lock(confirm)
		defer mem.Unlock(confirm)
		defer mem.ZeroBytes(confirm)

		if len(password) == 0 {
			return fmt.Errorf("password cannot be empty")
		}

		if !bytes.Equal(password, confirm) {
			return fmt.Errorf("passwords do not match")
		}

		encrypted, err := crypto.EncryptPrivateKey(priv, password)
		if err != nil {
			return fmt.Errorf("encrypt private key: %w", err)
		}

		if err := os.WriteFile(filepath, encrypted, 0600); err != nil {
			return fmt.Errorf("write key file: %w", err)
		}

		fmt.Println("Private key save to", filepath)

		return nil
	},
}
