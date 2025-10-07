package cli

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
)

var insertCmd = &cobra.Command{
	Use:   "insert name [password]",
	Short: "Insert new password to storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		var password, name string

		p, err := pswd.NewPswd("", "")
		if err != nil {
			return err
		}
		switch len(args) {
		case 0:
			return PassArgumentsErr("name")
		case 1:
			{
				name = args[0]
				if err := checkInsert(p, name); err != nil {
					return err
				}
				password, err = readPasswordFromStdin()
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
		if err := checkInsert(p, name); err != nil {
			return err
		}
		if _, err := p.Insert(name, password); err != nil {
			return err
		}
		fmt.Println("New password saved as", s.Passname.Render(name))
		return nil
	},
}

func checkInsert(p *pswd.Pswd, name string) error {
	if _, err := p.Type(name); err == nil {
		return fmt.Errorf("%s already exists", s.Passname.Render(name))
	}
	return nil
}

func readPasswordFromStdin() (string, error) {
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
