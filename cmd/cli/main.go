package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"golang.org/x/term"

	"github.com/vladislav-atakhanov/pswd/internal/uuid"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

func readPassword(prompt string) ([]byte, error) {
	fmt.Print(prompt)
	password, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	return password, err
}

var rootCmd = &cobra.Command{
	Use:   "pswd",
	Short: "Password manager",
	Long:  "A CLI password manager with X25519 encryption.",
	Args:  cobra.ArbitraryArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			return cmd.Help()
		}
		hexID := args[0]

		ctx, err := openVault(cmd)
		if err != nil {
			return err
		}
		defer closeVault(ctx)

		raw, err := hex.DecodeString(hexID)
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
				return fmt.Errorf("password %q not found", hexID)
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
	rootCmd.PersistentFlags().StringP("vault", "p", "", "path to vault file")
	rootCmd.PersistentFlags().StringP("key", "k", "", "path to private key file")
	rootCmd.PersistentFlags().BoolP("clip", "c", false, "copy password to clipboard")
}

func main() {
	rootCmd.AddCommand(genkeyCmd, exportCmd, masterCmd, initCmd, addDeviceCmd, addCmd, searchCmd, removeCmd, renameCmd, compactCmd, removeDeviceCmd, generateCmd, editCmd, listCmd, statsCmd)

	ctx := trapSignals()
	exitCode := 0
	defer func() { os.Exit(exitCode) }()

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		fmt.Fprintln(os.Stderr, err)
		exitCode = 1
		return
	}
}
