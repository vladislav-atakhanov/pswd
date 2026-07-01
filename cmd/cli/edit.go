package main

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/uuid"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

var editCmd = &cobra.Command{
	Use:   "edit <hex-id>",
	Short: "Edit a password entry",
	Long: `Edit a password entry content using an external editor.

By default, decrypts the entry, opens the system editor (EDITOR/VISUAL env),
and re-encrypts on save if the content changed.

Use --manual to decrypt to a temp file and wait for Enter to re-encrypt.`,
	Args: cobra.ExactArgs(1),
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

		original, err := io.ReadAll(r)
		if err != nil {
			return err
		}
		defer mem.ZeroBytes(original)

		tmpFile, err := os.CreateTemp("", "pswd-edit-*")
		if err != nil {
			return fmt.Errorf("create temp file: %w", err)
		}
		tmpPath := tmpFile.Name()
		defer os.Remove(tmpPath)

		if _, err := tmpFile.Write(original); err != nil {
			tmpFile.Close()
			return fmt.Errorf("write temp file: %w", err)
		}
		if err := tmpFile.Close(); err != nil {
			return fmt.Errorf("close temp file: %w", err)
		}

		manual, _ := cmd.Flags().GetBool("manual")

		if manual {
			fmt.Printf("Decrypted to: %s\nPress Enter to update or Ctrl+C to cancel", tmpPath)
			fmt.Scanln()
		} else {
			editor := findEditor()
			editorCmd := exec.Command(editor, tmpPath)
			editorCmd.Stdin = os.Stdin
			editorCmd.Stdout = os.Stdout
			editorCmd.Stderr = os.Stderr
			if err := editorCmd.Run(); err != nil {
				return fmt.Errorf("editor %q: %w", editor, err)
			}
		}

		newContent, err := os.ReadFile(tmpPath)
		if err != nil {
			return fmt.Errorf("read edited file: %w", err)
		}
		defer mem.ZeroBytes(newContent)

		if bytes.Equal(original, newContent) {
			fmt.Println("No changes made")
			return nil
		}

		if err := ctx.Vault.Update(id, bytes.NewReader(newContent)); err != nil {
			return fmt.Errorf("update password: %w", err)
		}
		if err := ctx.Vault.Save(ctx.File); err != nil {
			return fmt.Errorf("save vault: %w", err)
		}

		fmt.Println("Password updated")
		return nil
	},
}

func findEditor() string {
	if e := os.Getenv("VISUAL"); e != "" {
		return e
	}
	if e := os.Getenv("EDITOR"); e != "" {
		return e
	}
	if runtime.GOOS == "windows" {
		return "notepad.exe"
	}
	return "vim"
}

func init() {
	editCmd.Flags().Bool("manual", false, "decrypt to file and wait for Enter to update")
}
