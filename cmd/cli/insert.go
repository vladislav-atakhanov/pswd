package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd/pkg/pswd"
)

var insertCmd = &cobra.Command{
	Use:   "insert name [password]",
	Short: "Insert new password to storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		var password, name string
		var err error
		switch len(args) {
		case 0:
			return PassArgumentsErr("name")
		case 1:
			{
				name = args[0]
				password, err = readPassword()
				if err != nil {
					return err
				}
			}
		case 2:
			{
				name = args[0]
				password = args[1]
			}
		default:
			return TooManyArgumentsErr()
		}
		p, err := pswd.NewPswd("")
		if err != nil {
			return err
		}
		passfile, err := p.Insert(name, password)
		if err != nil {
			return err
		}
		fmt.Println("saved to", passfile)
		return nil
	},
}

func readPassword() (string, error) {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return "", err
	}
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", err
		}
		return strings.TrimSpace(string(data)), nil
	}
	return promptPassword("Enter password: ", "Repeat password: ")
}

func registerInsert(c *cobra.Command) {
	c.AddCommand(insertCmd)
}
