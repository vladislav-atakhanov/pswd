package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info -p <vault> -k <key>",
	Short: "Display vault contents",
	Args:  cobra.NoArgs,
	RunE: withVault(func(ctx *vaultContext, args []string) error {
		fmt.Println("Devices:")
		for _, d := range ctx.Vault.Devices {
			key := d.PublicKey()
			fmt.Printf("\t%s | %s\n", d.Name(), base64.URLEncoding.EncodeToString(key[:]))
		}
		fmt.Println("Passwords:")
		for _, entry := range ctx.Vault.List() {
			t := time.Unix(int64(entry.LastUpdate), 0).Format("2006-01-02 15:04:05")
			fmt.Printf("\t%s | %s | %s\n", entry.Label, hex.EncodeToString(entry.ID[:]), t)
		}
		return nil
	}),
}
