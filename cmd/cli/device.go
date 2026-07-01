package main

import (
	"github.com/spf13/cobra"
)

var deviceCmd = &cobra.Command{
	Use:   "device <command>",
	Short: "Manage devices",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func init() {
	deviceCmd.AddCommand(deviceAddCmd, deviceRemoveCmd, deviceExportCmd)
}
