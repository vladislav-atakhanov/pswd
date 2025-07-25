package cli

import (
	"fmt"
	"io"
	"pswd/pkg/pswd"

	"github.com/spf13/cobra"
)

var insertCmd = &cobra.Command{
	Use:   "insert name [password]",
	Short: "Insert new password to storage",
	RunE: func(cmd *cobra.Command, args []string) error {
		p, err := pswd.NewPswd("")
		if err != nil {
			return err
		}
		var password, name string
		switch len(args) {
		case 0:
			return fmt.Errorf("recieve name")
		case 1:
			{
				name = args[0]
				var inputReader io.Reader = cmd.InOrStdin()
				c, err := io.ReadAll(inputReader)
				if err != nil {
					return err
				}
				if len(c) > 0 {
					password = string(c)
				} else {
					password, err = promptPassword(true, "")
					if err != nil {
						return err
					}
				}
			}
		case 2:
			{
				name = args[0]
				password = args[1]
			}
		default:
			return fmt.Errorf("too many arguments")
		}
		passfile, err := p.Insert(name, password)
		if err != nil {
			return err
		}
		fmt.Println("saved to", passfile)
		return nil
	},
}

func registerInsert(c *cobra.Command) {
	c.AddCommand(insertCmd)
}
