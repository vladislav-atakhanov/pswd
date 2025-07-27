package styles

import "github.com/charmbracelet/lipgloss"

var (
	grayStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	blueStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("4"))
	greenStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	redStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	yellowStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	cyanStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("6"))
	magentaStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	whiteStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("15"))
	blackStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("0"))

	boldStyle      = lipgloss.NewStyle().Bold(true)
	underlineStyle = lipgloss.NewStyle().Underline(true)
	resetStyle     = lipgloss.NewStyle()
)

var (
	Key    = greenStyle
	Error  = redStyle
	Dir    = blueStyle
	Data   = yellowStyle
	Pass   = cyanStyle
	Prefix = grayStyle
	Magic  = magentaStyle
	File   = resetStyle
)
