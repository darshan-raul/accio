package ui

import "github.com/charmbracelet/lipgloss"

var (
	logoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7D56F4"))

	subHeaderStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#626262")).
			Align(lipgloss.Center)

	versionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD93D")).
			Italic(true)
)

func RenderHeader() string {
	logo := `
   █████╗  ██████╗ ██████╗██╗ ██████╗ 
  ██╔══██╗██╔════╝██╔════╝██║██╔═══██╗
  ███████║██║     ██║     ██║██║   ██║
  ██╔══██║██║     ██║     ██║██║   ██║
  ██║  ██║╚██████╗╚██████╗██║╚██████╔╝
  ╚═╝  ╚═╝ ╚═════╝ ╚═════╝╚═╝ ╚═════╝ 
`

	tagline := "AI-Powered Cloud Infrastructure Platform"
	version := "v1.0.0"

	return lipgloss.JoinVertical(
		lipgloss.Center,
		logoStyle.Render(logo),
		"",
		subHeaderStyle.Render(tagline),
		versionStyle.Render(version),
		"",
	)
}

func RenderWelcome() string {
	return subHeaderStyle.Render("✨ Manage AWS, Azure, and GCP with natural language ✨\n")
}
