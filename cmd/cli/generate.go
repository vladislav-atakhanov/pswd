package cli

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/sethvargo/go-password/password"
	"github.com/spf13/cobra"
	"github.com/vladislav-atakhanov/pswd"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
)

const defaultLength = 25

var generateCmd = &cobra.Command{
	Use:   "generate passname [length]",
	Short: "generate password and save by name",
	RunE: func(cmd *cobra.Command, args []string) error {
		var name string
		length := defaultLength
		switch len(args) {
		case 0:
			return PassArgumentsErr("name")
		case 1:
			name = args[0]
		case 2:
			name = args[0]
			l, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("Length must be an integer. Passed %s", s.Error.Render(args[1]))
			}
			length = l
		default:
			return TooManyArgumentsErr()
		}
		clip, _ := cmd.Flags().GetBool("clip")
		noSymbols, _ := cmd.Flags().GetBool("no-symbols")
		p, err := pswd.NewPswd("")
		if err != nil {
			return err
		}
		password, err := generatePassword(length, noSymbols)
		if err != nil {
			return err
		}
		if _, err := p.Insert(name, password); err != nil {
			return err
		}
		if clip {
			if err := clipboard.WriteAll(password); err != nil {
				return err
			}
			fmt.Printf("Password saved as %s and copied to clipboard\n", s.Pass.Render(name))
		} else {
			fmt.Printf("Password saved as %s and here:\n%s\n", s.Pass.Render(name), s.Data.Render(password))
		}
		return nil
	},
}

func randInt(min, max int) (int, error) {
	if min >= max {
		return 0, fmt.Errorf("invalid range: min >= max")
	}

	nBig, err := rand.Int(rand.Reader, big.NewInt(int64(max-min)))
	if err != nil {
		return 0, err
	}

	return int(nBig.Int64()) + min, nil
}

func generatePassword(length int, noSymbols bool) (string, error) {
	numDigits, err := randInt(2, length/3)
	if err != nil {
		return "", err
	}
	numSymbols := 0
	if !noSymbols {
		numSymbols, err = randInt(2, length/3)
		if err != nil {
			return "", err
		}
	}
	return password.Generate(length, numDigits, numSymbols, false, true)
}

func registerGenerate(c *cobra.Command) {
	c.AddCommand(generateCmd)
	generateCmd.Flags().BoolP("clip", "c", false, "copy generated password to clipboard")
	generateCmd.Flags().BoolP("no-symbols", "n", false, "copy generated password to clipboard")
}
