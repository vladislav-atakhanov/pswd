package main

import (
	"testing"

	"github.com/spf13/cobra"
)

type cmdCheck struct {
	cmd  *cobra.Command
	args cobra.PositionalArgs
	name string
}

func collectCommands(cmd *cobra.Command) []cmdCheck {
	var cmds []cmdCheck
	for _, sub := range cmd.Commands() {
		cmds = append(cmds, cmdCheck{cmd: sub, args: sub.Args, name: sub.Name()})
	}
	return cmds
}

func TestRootCmdHasAllSubcommands(t *testing.T) {
	rootCmd.ResetCommands()
	rootCmd.AddCommand(
		addCmd, showCmd, editCmd, removeCmd, renameCmd, listCmd,
		searchCmd, generateCmd,
		deviceCmd,
		initCmd, compactCmd, statsCmd, passwdCmd,
		genkeyCmd, completionCmd,
	)

	expected := map[string]bool{
		"add": true, "show": true, "edit": true, "remove": true,
		"rename": true, "list": true, "search": true, "generate": true,
		"device": true,
		"init": true, "compact": true, "stats": true, "passwd": true,
		"genkey": true, "completion": true,
	}
	for _, sub := range rootCmd.Commands() {
		delete(expected, sub.Name())
	}
	if len(expected) > 0 {
		for name := range expected {
			t.Errorf("missing subcommand: %s", name)
		}
	}
}

func TestCommandArgs(t *testing.T) {
	rootCmd.ResetCommands()
	rootCmd.AddCommand(
		addCmd, showCmd, editCmd, removeCmd, renameCmd, listCmd,
		searchCmd, generateCmd,
		deviceCmd,
		initCmd, compactCmd, statsCmd, passwdCmd,
		genkeyCmd, completionCmd,
	)

	expected := map[string]cobra.PositionalArgs{
		"compact":    cobra.NoArgs,
		"stats":      cobra.NoArgs,
		"list":       cobra.NoArgs,
		"add":        cobra.RangeArgs(1, 2),
		"generate":   cobra.RangeArgs(1, 2),
		"show":       cobra.ExactArgs(1),
		"edit":       cobra.ExactArgs(1),
		"remove":     cobra.ExactArgs(1),
		"rename":     cobra.ExactArgs(2),
		"search":     cobra.ExactArgs(1),
		"init":       cobra.ExactArgs(2),
		"genkey":     cobra.ExactArgs(1),
		"passwd":     cobra.ExactArgs(1),
		"device":     cobra.NoArgs,
		"completion": cobra.ExactArgs(1),
	}
	for _, cc := range collectCommands(rootCmd) {
		want, ok := expected[cc.name]
		if !ok {
			t.Errorf("unexpected command: %s", cc.name)
			continue
		}
		if cc.args == nil && want != nil {
			t.Errorf("%s: args is nil, want non-nil", cc.name)
		}
	}
}

func TestCommandYesFlag(t *testing.T) {
	rootCmd.ResetCommands()
	rootCmd.AddCommand(
		addCmd, showCmd, editCmd, removeCmd, renameCmd, listCmd,
		searchCmd, generateCmd,
		deviceCmd,
		initCmd, compactCmd, statsCmd, passwdCmd,
		genkeyCmd, completionCmd,
	)

	for _, cc := range collectCommands(rootCmd) {
		flag := cc.cmd.Flags().Lookup("yes")
		switch cc.name {
		case "compact", "remove":
			if flag == nil {
				t.Errorf("%s: expected --yes flag, got nil", cc.name)
			} else if flag.Value.Type() != "bool" {
				t.Errorf("%s: --yes should be bool, got %s", cc.name, flag.Value.Type())
			}
		default:
			if flag != nil {
				t.Errorf("%s: unexpected --yes flag", cc.name)
			}
		}
	}
}
