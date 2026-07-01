package main

import (
	"fmt"
	"os"

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
		return nil, fmt.Errorf("flag --vault (-p) is required")
	}
	if keyPath == "" {
		return nil, fmt.Errorf("flag --key (-k) is required")
	}

	data, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read key file: %w", err)
	}

	password, err := readPassword("Enter master password: ")
	if err != nil {
		return nil, fmt.Errorf("read password: %w", err)
	}
	defer mem.ZeroBytes(password)

	priv, err := crypto.DecryptPrivateKey(data, password)
	if err != nil {
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
		return nil, fmt.Errorf("open vault: %w", err)
	}

	return &vaultContext{Vault: v, File: file, Priv: priv}, nil
}

func closeVault(ctx *vaultContext) {
	ctx.File.Close()
	mem.ZeroArray32(&ctx.Priv)
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
