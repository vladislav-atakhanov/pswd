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

	resetStyle = lipgloss.NewStyle()
)

var (
	KeyID     = yellowStyle
	Passname  = cyanStyle.Bold(true)
	Error     = redStyle
	Dir       = blueStyle.Bold(true)
	Data      = greenStyle
	Secondary = grayStyle
	File      = resetStyle
	Program   = cyanStyle
)
