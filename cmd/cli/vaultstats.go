package main

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Display vault statistics and fragmentation",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, err := openVault(cmd)
		if err != nil {
			return err
		}
		defer closeVault(ctx)

		stat, err := ctx.File.Stat()
		if err != nil {
			return err
		}
		fileSize := int(stat.Size())

		s := ctx.Vault.Stats()

		indexSize := fileSize - s.DataEnd

		fmt.Println("Devices:")
		for _, d := range ctx.Vault.Devices() {
			key := d.PublicKey()
			fmt.Printf("\t%s | %s\n", d.Name(), base64.URLEncoding.EncodeToString(key[:]))
		}
		fmt.Println()

		fmt.Println("File layout:")
		fmt.Printf("  Header\t%s\n", formatSize(s.HeaderSize))
		fmt.Printf("  Body  \t%s\n", formatSize(s.BodySize))
		fmt.Printf("  Index \t%s\n", formatSize(indexSize))
		fmt.Printf("  Total \t%s\n", formatSize(fileSize))

		if s.OrphanedCount > 0 {
			pct := float64(s.OrphanedSize) / float64(s.BodySize) * 100
			fmt.Printf("\nFragmentation:\n")
			fmt.Printf("  Orphaned spans\t%d\n", s.OrphanedCount)
			fmt.Printf("  Orphaned space\t%s\n", formatSize(s.OrphanedSize))
			fmt.Printf("  Fragmentation \t%.1f%%\n", pct)
		}

		return nil
	},
}

func formatSize(b int) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	if b < unit*unit {
		return fmt.Sprintf("%.1f KiB", float64(b)/unit)
	}
	return fmt.Sprintf("%.1f MiB", float64(b)/(unit*unit))
}


