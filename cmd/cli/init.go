package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd/pkg/pswd"
	"golang.org/x/term"
)

func registerInit(c *cobra.Command) {
	c.AddCommand(initCmd)
	initCmd.Flags().StringP("path", "p", "", "subfolder")
}

var initCmd = &cobra.Command{
	Use:   "init [password]",
	Short: "Initialize new password storage",
	RunE:  initStorage,
}

func initStorage(cmd *cobra.Command, args []string) error {
	subfolder, _ := cmd.Flags().GetString("path")
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
	}
	p, err := pswd.NewPswd("")
	if err != nil {
		return err
	}
	inited := p.IsInit(subfolder)
	if inited {
		fmt.Println("reinit", subfolder)
	}
	d, names, err := p.Init(subfolder, func() (string, error) {
		return getPassword(args, "new master password")
	}, func() (string, error) {
		return promptPassword(false, "old master password")
	})
	if err != nil {
		return err
	}
	for _, n := range names {
		fmt.Println(n, "reencrypt")
	}
	if inited {
		fmt.Printf("Password store at %s reinitialized\n", d)
	} else {
		fmt.Println("New password store initialized at", d)
	}
	return nil
}

func getPassword(args []string, label string) (string, error) {
	if len(args) == 1 {
		return args[0], nil
	}
	return promptPassword(true, label)
}

func promptPassword(confirm bool, label string) (string, error) {
	if label == "" {
		label = "password"
	}
	fmt.Printf("Enter %s: ", label)
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if confirm {
		fmt.Printf("Repeat %s: ", label)
		confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Println()
		if err != nil {
			return "", err
		}
		if string(passwordBytes) != string(confirmBytes) {
			return "", fmt.Errorf("passwords do not match")
		}
	}

	return string(passwordBytes), nil
}
