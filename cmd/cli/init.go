package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

var initCmd = &cobra.Command{
	Use:   "init <name> <public-key>",
	Short: "Create a new vault",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		name, pub, err := parsePublicKey(args[1])
		if err != nil {
			return err
		}
		file, err := os.Create(args[0])
		if err != nil {
			return fmt.Errorf("open vault: %w", err)
		}
		defer file.Close()
		v := vault.New(pub, name)
		return v.Save(file)
	},
}
