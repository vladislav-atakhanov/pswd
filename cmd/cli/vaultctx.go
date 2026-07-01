package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"

	"github.com/vladislav-atakhanov/pswd/internal/crypto"
	"github.com/vladislav-atakhanov/pswd/internal/mem"
	"github.com/vladislav-atakhanov/pswd/internal/vault"
)

type vaultContext struct {
	Vault *vault.Vault
	File  *os.File
	Priv  [32]byte
}

func openVault(cmd *cobra.Command) (*vaultContext, error) {
	vaultPath, _ := cmd.Flags().GetString("vault")
	keyPath, _ := cmd.Flags().GetString("key")

	if vaultPath == "" {
		vaultPath = os.Getenv("PSWD_VAULT")
	}
	if keyPath == "" {
		keyPath = os.Getenv("PSWD_KEY")
	}

	if vaultPath == "" {
		return nil, fmt.Errorf("flag --vault (-p) or PSWD_VAULT environment variable is required")
	}
	if keyPath == "" {
		return nil, fmt.Errorf("flag --key (-k) or PSWD_KEY environment variable is required")
	}

	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	password, err := readPassword("Enter master password: ")
	if err != nil {
		return nil, fmt.Errorf("read password: %w", err)
	}
	mem.Lock(password)
	defer mem.Unlock(password)
	defer mem.ZeroBytes(password)

	priv, err := crypto.DecryptPrivateKey(data, password)
	if err != nil {
		if errors.Is(err, crypto.ErrWrongPassword) {
			return nil, fmt.Errorf("wrong password: check your master password or key file")
		}
		return nil, fmt.Errorf("decrypt private key: %w", err)
	}

	file, err := os.OpenFile(vaultPath, os.O_RDWR, 0666)
	if err != nil {
		mem.ZeroArray32(&priv)
		return nil, fmt.Errorf("open vault: %w", err)
	}

	stat, err := file.Stat()
	if err != nil {
		file.Close()
		mem.ZeroArray32(&priv)
		return nil, err
	}

	v, err := vault.Open(file, int(stat.Size()), priv)
	if err != nil {
		file.Close()
		mem.ZeroArray32(&priv)
		if errors.Is(err, vault.ErrAccessDenied) {
			return nil, fmt.Errorf("access denied: your key is not authorized for this vault")
		}
		return nil, fmt.Errorf("open vault: %w", err)
	}

	return &vaultContext{Vault: v, File: file, Priv: priv}, nil
}

func closeVault(ctx *vaultContext) {
	ctx.File.Close()
	mem.ZeroArray32(&ctx.Priv)
}

func readNewPassword(prompt1, prompt2 string) ([]byte, error) {
	p1, err := readPassword(prompt1)
	if err != nil {
		return nil, fmt.Errorf("read password: %w", err)
	}
	mem.Lock(p1)
	defer mem.Unlock(p1)
	defer mem.ZeroBytes(p1)

	if len(p1) == 0 {
		return nil, errors.New("password cannot be empty")
	}

	p2, err := readPassword(prompt2)
	if err != nil {
		return nil, fmt.Errorf("read confirmation: %w", err)
	}
	mem.Lock(p2)
	defer mem.Unlock(p2)
	defer mem.ZeroBytes(p2)

	if !bytes.Equal(p1, p2) {
		return nil, errors.New("passwords do not match")
	}

	result := make([]byte, len(p1))
	copy(result, p1)
	return result, nil
}

func confirm(prompt string) bool {
	fmt.Print(prompt)
	reader := bufio.NewReader(os.Stdin)
	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(strings.ToLower(line))
	return line == "y" || line == "yes"
}

func withVault(run func(ctx *vaultContext, args []string) error) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		ctx, err := openVault(cmd)
		if err != nil {
			return err
		}
		defer closeVault(ctx)
		return run(ctx, args)
	}
}
