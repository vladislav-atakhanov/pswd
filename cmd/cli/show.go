package cli

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
)

var showCmd = &cobra.Command{
	Use:   "show name",
	Short: "show password by name",
	Run: withError(func(cmd *cobra.Command, args []string) error {
		var name string
		switch len(args) {
		case 0:
			return showTree("")
		case 1:
			name = args[0]
		default:
			return TooManyArgumentsErr()
		}
		clip, _ := cmd.Flags().GetBool("clip")
		p, err := pswd.NewPswd("", "")
		if err != nil {
			return err
		}
		t, err := p.Type(name)
		if err != nil {
			return err
		}
		switch t {
		case pswd.PassDir:
			return showTree(name)
		case pswd.PassUnknown:
			return showTree(name)
		}
		data, err := p.ShowLazy(name, enterMasterPassword)
		if err != nil {
			return err
		}
		if clip {
			password := strings.TrimSpace(firstLine(data))
			if err := clipboard.WriteAll(password); err != nil {
				return err
			}
			fmt.Printf("Copied %s to clipboard\n", s.Passname.Render(name))
		} else {
			fmt.Println(s.Data.Render(data))
		}
		return nil
	}),
}

func enterMasterPassword(key string) (string, error) {
	return promptPassword(fmt.Sprintf("Enter password for %s key: ", s.KeyID.Render(key)), "")
}

func showTree(name string) error {
	p, err := pswd.NewPswd("", "")
	if err != nil {
		return err
	}
	tree, err := p.Tree(name)
	if err != nil {
		return err
	}
	fmt.Println(s.Dir.Render(tree.Name))
	for i, child := range tree.Children {
		printTree(child, "", i == len(tree.Children)-1)
	}
	return nil
}

func printTree(node *pswd.TreeNode, prefix string, isLast bool) {
	branch := "├── "
	if isLast {
		branch = "└── "
	}
	color := s.File
	if node.IsDir {
		color = s.Dir
	}

	fmt.Print(s.Secondary.Render(prefix))
	fmt.Print(s.Secondary.Render(branch))
	fmt.Println(color.Render(node.Name))

	newPrefix := prefix
	if isLast {
		newPrefix += "    "
	} else {
		newPrefix += "│   "
	}

	for i, child := range node.Children {
		printTree(child, newPrefix, i == len(node.Children)-1)
	}
}

func firstLine(s string) string {
	lines := strings.SplitN(s, "\n", 2)
	return lines[0]
}

func registerShow(c *cobra.Command) {
	c.AddCommand(showCmd)
	showCmd.Flags().BoolP("clip", "c", false, "copy first line to clipboard")
}

func withError(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			fmt.Println("Error:", err)
		}
	}
}
