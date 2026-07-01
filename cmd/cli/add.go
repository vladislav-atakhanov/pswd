package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add <label> [filename]",
	Short: "Add a password to the vault",
	Long:  "Add a password entry to the vault. If filename is given, read from file; otherwise prompt for the password.",
	Args:  cobra.RangeArgs(1, 2),
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		label := args[0]

		var plain io.Reader
		if len(args) == 2 {
			f, err := os.Open(args[1])
			if err != nil {
				return fmt.Errorf("open file: %w", err)
			}
			defer f.Close()
			plain = f
		} else {
			p1, err := readPassword("Enter password: ")
			if err != nil {
				return fmt.Errorf("read password: %w", err)
			}
			if len(p1) == 0 {
				return fmt.Errorf("password cannot be empty")
			}
			p2, err := readPassword("Confirm password: ")
			if err != nil {
				return fmt.Errorf("read confirmation: %w", err)
			}
			if !bytes.Equal(p1, p2) {
				return fmt.Errorf("passwords do not match")
			}
			plain = strings.NewReader(string(p1))
		}

		if err := ctx.Vault.Add(plain, label); err != nil {
			return fmt.Errorf("add password: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Password added with label", label)
		return nil
	}),
}
