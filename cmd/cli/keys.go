package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd/pkg/keys"
	"golang.org/x/term"
)

var genereteKeysCmd = &cobra.Command{
	Use:   "gen-key id [password]",
	Short: "generate new key-pair and save to ~/.keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		var id, password string
		var err error
		switch len(args) {
		case 0:
			return PassArgumentsErr("id")
		case 1:
			id = args[0]
			if err := check(id); err != nil {
				return err
			}
			keyLabel := keyColor(id)
			password, err = promptPassword(
				fmt.Sprintf("Enter password for %s key: ", keyLabel),
				fmt.Sprintf("Repeat password for %s key: ", keyLabel),
			)
			if err != nil {
				return err
			}
		case 2:
			id = args[0]
			if err := check(id); err != nil {
				return err
			}
			password = args[1]
		default:
			return TooManyArgumentsErr()
		}
		priv, pub, err := keys.Generate(password)
		if err != nil {
			return err
		}
		d, err := keys.Save(id, priv, pub)
		if err != nil {
			return err
		}
		fmt.Println("New keys generated and saved to", dirColor(d))
		return nil
	},
}

func check(id string) error {
	if keys.Has(id) {
		return fmt.Errorf("Key pair %s already exists", keyColor(id))
	}
	return nil
}

func promptPassword(label, confirmLabel string) (string, error) {
	fmt.Print(label)
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println()
	if err != nil {
		return "", err
	}

	if confirmLabel != "" {
		fmt.Print(confirmLabel)
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

func registerGenerateKeys(c *cobra.Command) {
	c.AddCommand(genereteKeysCmd)
}
