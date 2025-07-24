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
	p, err := pswd.NewPswd("")
	if err != nil {
		return err
	}
	if len(args) > 1 {
		return fmt.Errorf("too many arguments")
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

func promptPassword(confirm bool) (string, error) {
	fmt.Print("Enter password: ")
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if confirm {
		fmt.Print("Repeat password: ")
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
