package cli

import (
	"fmt"
	"os"
	"pswd/pkg/pswd"

	"github.com/spf13/cobra"
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
	p, err := pswd.NewPswd("")
	if err != nil {
		return err
	}
	}
	var password string
	if len(args) == 1 {
		password = args[0]
	} else {
		password, err = promptPassword(true)
		if err != nil {
			return err
		}
	}
	d, err := p.Init(subfolder, password)
	if err != nil {
		return err
	}
	fmt.Println("New password store initialized at", d)
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
