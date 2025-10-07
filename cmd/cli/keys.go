package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
	"github.com/vladislav-atakhanov/pswd/pkg/keys"
	"golang.org/x/term"
)

var genereteKeysCmd = &cobra.Command{
	Use:   "gen-key id [password]",
	Short: "generate new key-pair and save to ~/.keys",
	RunE: func(cmd *cobra.Command, args []string) error {
		var id, password string
		var err error

		ks, err := keys.NewKeyStore("")
		if err != nil {
			return err
		}

		switch len(args) {
		case 0:
			return PassArgumentsErr("id")
		case 1:
			id = args[0]
			if err := check(ks, id); err != nil {
				return err
			}
			keyLabel := s.KeyID.Render(id)
			password, err = promptPassword(
				fmt.Sprintf("Enter password for %s key: ", keyLabel),
				fmt.Sprintf("Repeat password for %s key: ", keyLabel),
			)
			if err != nil {
				return err
			}
		case 2:
			id = args[0]
			if err := check(ks, id); err != nil {
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
		d, err := ks.Save(id, priv, pub)
		if err != nil {
			return err
		}
		fmt.Println("New keys generated and saved to", s.Dir.Render(d))
		return nil
	},
}

func check(ks *keys.KeyStore, id string) error {
	if ks.Exists(id) {
		return fmt.Errorf("Key pair %s already exists", s.KeyID.Render(id))
	}
	return nil
}

func promptPassword(label, confirmLabel string) (string, error) {
	fmt.Fprint(os.Stderr, label)
	passwordBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", err
	}

	if confirmLabel != "" {
		fmt.Fprint(os.Stderr, confirmLabel)
		confirmBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
		fmt.Fprintln(os.Stderr)
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
