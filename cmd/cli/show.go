package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/uuid"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

var showCmd = &cobra.Command{
	Use:   "show <hex-id>",
	Short: "Show a password entry",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := openVault(cmd)
		if err != nil {
			return err
		}
		defer closeVault(ctx)

		raw, err := hex.DecodeString(args[0])
		if err != nil {
			return fmt.Errorf("decode hex id: %w", err)
		}
		if len(raw) != 16 {
			return fmt.Errorf("invalid hex id length: got %d bytes, expected 16", len(raw))
		}
		var id uuid.V4
		copy(id[:], raw)

		r, err := ctx.Vault.Read(ctx.File, id, ctx.Priv)
		if err != nil {
			if errors.Is(err, vault.ErrNotFound) {
				return fmt.Errorf("password %q not found", args[0])
			}
			return fmt.Errorf("read password: %w", err)
		}

		content, err := io.ReadAll(r)
		if err != nil {
			return err
		}

		clip, _ := cmd.Flags().GetBool("clip")
		if clip {
			firstLine, _, _ := strings.Cut(string(content), "\n")
			if err := clipboardWrite(firstLine); err != nil {
				return fmt.Errorf("clipboard: %w", err)
			}
			return nil
		}

		fmt.Print(string(content))
		return nil
	},
}

func init() {
	showCmd.Flags().BoolP("clip", "c", false, "copy password to clipboard")
}
