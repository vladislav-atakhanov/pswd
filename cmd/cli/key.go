package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
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

		password, err := readNewPassword("Enter master password: ", "Confirm master password: ")
		if err != nil {
			return err
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
