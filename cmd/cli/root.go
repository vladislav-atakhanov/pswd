package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
	s "github.com/vladislav-atakhanov/pswd/cmd/cli/styles"
)

const banner = `
                              '||
... ...   .... ... ... ...  .. ||
 ||'  || ||. '  ||  ||  | .'  '||
 ||    | . '|..  ||| |||  |.   ||
 ||...'  |'..|'   |   |   '|..'||.
 ||
''''
`

var program = s.Program.Render("pswd")
var rootCmd = &cobra.Command{
	Use:  s.Program.Render("pswd"),
	Long: fmt.Sprintf("%s is a password manager inspired by %s", program, s.Data.Render("pass")),
}

func printCenter(elem ...string) {
	var w int
	for _, e := range elem {
		for line := range strings.Lines(e) {
			l := lipgloss.Width(strings.TrimSpace(line))
			if l > w {
				w = l
			}
		}
	}
	for _, e := range elem {
		padding := (w - lipgloss.Width(e)) / 2
		fmt.Println(lipgloss.NewStyle().PaddingLeft(padding).Render(e))
	}
}

func printBanner() {
	printCenter(
		s.Program.Render(strings.Trim(strings.ReplaceAll(banner, "$", "`"), "\n")),
		strings.Join([]string{"secure", "cross-platform", "fast"}, s.Secondary.Render(" - ")),
		"",
		fmt.Sprintf("Use \"%s %s\" for show password", program, s.Passname.Render("pass-name")),
	)
}

func Execute() error {
	registerInit(rootCmd)
	registerInsert(rootCmd)
	registerShow(rootCmd)
	registerGenerate(rootCmd)
	registerGenerateKeys(rootCmd)

	args := os.Args[1:]

	if len(args) == 0 {
		printBanner()
		return nil
	}
	switch args[0] {
	case "help", "--help", "-h":
		return rootCmd.Help()
	}
	cmd, _, err := rootCmd.Find(args)
	if err != nil || !cmd.IsAvailableCommand() || cmd == rootCmd {
		os.Args = append(os.Args[:1], append([]string{"show"}, args...)...)
	}

	return rootCmd.Execute()
}
