package cli

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
)

func registerInit(c *cobra.Command) {
	c.AddCommand(initCmd)
	initCmd.Flags().StringP("path", "p", "", "subfolder")
}

var initCmd = &cobra.Command{
	Use:   "init key-id",
	Short: "Initialize new password storage with key",
	RunE: func(cmd *cobra.Command, args []string) error {
		subfolder, _ := cmd.Flags().GetString("path")
		var keyId string
		switch len(args) {
		case 0:
			return PassArgumentsErr("key-id")
		case 1:
			keyId = args[0]
		default:
			return TooManyArgumentsErr()
		}
		p, err := pswd.NewPswd("", "")
		if err != nil {
			return err
		}
		names := make(chan string)
		go func() {
			for n := range names {
				fmt.Println("Password", s.Passname.Render(n), "reencrypt")
			}
		}()
		d, reinit, err := p.Init(subfolder, keyId, enterMasterPassword, names)
		if err != nil {
			return err
		}
		if reinit {
			fmt.Printf("Password store at %s reinitialized\n", s.Dir.Render(d))
		} else {
			fmt.Println("New password store initialized at", s.Dir.Render(d))
		}
		return nil
	},
}
