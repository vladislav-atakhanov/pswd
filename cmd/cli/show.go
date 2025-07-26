package cli

import (
	"fmt"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd"
)

var showCmd = &cobra.Command{
	Use:   "show name",
	Short: "show password by name",
	Run: withError(func(cmd *cobra.Command, args []string) error {
		var name string
		switch len(args) {
		case 0:
			return PassArgumentsErr("name")
		case 1:
			name = args[0]
		default:
			return TooManyArgumentsErr()
		}
		clip, _ := cmd.Flags().GetBool("clip")
		p, err := pswd.NewPswd("")
		if err != nil {
			return err
		}
		data, err := p.ShowLazy(name, func(key string) (string, error) {
			return promptPassword(fmt.Sprintf("Enter password for %s key: ", key), "")
		})
		if err != nil {
			return err
		}
		if clip {
			password := strings.TrimSpace(firstLine(data))
			if err := clipboard.WriteAll(password); err != nil {
				return err
			}
			fmt.Printf("Copied %s to clipboard\n", name)
		} else {
			fmt.Println(data)
		}
		return nil
	}),
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
