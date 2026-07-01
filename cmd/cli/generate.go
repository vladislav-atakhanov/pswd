package main

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
)

const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+"

var generateCmd = &cobra.Command{
	Use:   "generate <label> [length]",
	Short: "Generate a random password and add it to the vault",
	Args:  cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := openVault(cmd)
		if err != nil {
			return err
		}
		defer closeVault(ctx)

		label := args[0]

		length := 24
		if len(args) == 2 {
			n, err := fmt.Sscanf(args[1], "%d", &length)
			if err != nil || n != 1 || length < 8 {
				return fmt.Errorf("invalid length: must be at least 8")
			}
		}

		password := make([]byte, length)
		for i := range password {
			idx, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
			if err != nil {
				return fmt.Errorf("generate password: %w", err)
			}
			password[i] = charset[idx.Int64()]
		}
		defer mem.ZeroBytes(password)

		if err := ctx.Vault.Add(bytes.NewReader(password), label); err != nil {
			return fmt.Errorf("add password: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		clip, _ := cmd.Flags().GetBool("clip")
		if clip {
			if err := clipboardWrite(string(password)); err != nil {
				return fmt.Errorf("clipboard: %w", err)
			}
			return nil
		}

		fmt.Printf("Password generated with label %s\n%s\n", label, string(password))
		return nil
	},
}
