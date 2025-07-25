package cli

import (
	"fmt"
	"pswd/pkg/pswd"

	"github.com/spf13/cobra"
)

var showCmd = &cobra.Command{
	Use:   "show passname",
	Short: "show password by name",
	Run: withError(func(cmd *cobra.Command, args []string) error {
		p, err := pswd.NewPswd("")
		if err != nil {
			return err
		}
		var name string
		switch len(args) {
		case 0:
			return fmt.Errorf("recieve name")
		case 1:
			name = args[0]
		default:
			return fmt.Errorf("too many arguments")
		}
		data, err := p.Show(name, func() (string, error) {
			return promptPassword(false, "")
		})
		if err != nil {
			return err
		}
		fmt.Println(data)
		return nil
	}),
}

func registerShow(c *cobra.Command) {
	c.AddCommand(showCmd)
}

func withError(f func(cmd *cobra.Command, args []string) error) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		if err := f(cmd, args); err != nil {
			fmt.Println(err)
		}
	}
}
