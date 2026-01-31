package styles

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	PrimaryColor   = lipgloss.Color("#7D56F4")
	SecondaryColor = lipgloss.Color("#FFD93D")
	SubtleColor    = lipgloss.Color("#626262")
	ErrorColor     = lipgloss.Color("#FF5555")
	SuccessColor   = lipgloss.Color("#50FA7B")

	// Text Styles
	TitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	StatusValidStyle = lipgloss.NewStyle().
				Foreground(SuccessColor)

	StatusInvalidStyle = lipgloss.NewStyle().
				Foreground(ErrorColor)

	SubtleStyle = lipgloss.NewStyle().
			Foreground(SubtleColor)

	// Box Styles
	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(PrimaryColor).
			Padding(1, 2)

	FocusedBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(SecondaryColor).
				Padding(1, 2)
)
