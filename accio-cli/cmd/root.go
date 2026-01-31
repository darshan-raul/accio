package cmd

import (
	"fmt"
	"os"

	"accio/internal/ui"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

type menuItem struct {
	title, desc string
	action      func()
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.desc }
func (i menuItem) FilterValue() string { return i.title }

type menuModel struct {
	list   list.Model
	choice string
	quit   bool
}

func (m menuModel) Init() tea.Cmd {
	return nil
}

func (m menuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.quit = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(menuItem)
			if ok {
				m.choice = i.title
				m.quit = true
				return m, tea.Quit
			}
		}
	case tea.WindowSizeMsg:
		m.list.SetWidth(msg.Width)
		return m, nil
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m menuModel) View() string {
	if m.quit {
		return ""
	}

	header := ui.RenderHeader()
	welcome := ui.RenderWelcome()

	doc := lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		welcome,
		"\n",
		m.list.View(),
	)

	return lipgloss.Place(
		80, 30,
		lipgloss.Center, lipgloss.Top,
		doc,
	)
}

var rootCmd = &cobra.Command{
	Use:   "accio",
	Short: "Accio is an AI-powered cloud infrastructure assistant",
	Long: `Accio enables you to manage cloud infrastructure across AWS, Azure, and GCP
using natural language and GitOps workflows.`,
	Run: func(cmd *cobra.Command, args []string) {
		items := []list.Item{
			menuItem{title: "🔐 Login", desc: "Authenticate with Accio Platform"},
			menuItem{title: "🏗️  Check Infrastructure", desc: "View current stacks and resources"},
			menuItem{title: "➕ Create Resource", desc: "Launch the infrastructure creation wizard"},
			menuItem{title: "💬 Chat", desc: "Interactive AI assistant"},
			menuItem{title: "⚙️  Settings", desc: "Configure API endpoint and context"},
			menuItem{title: "🚪 Exit", desc: "Exit Accio CLI"},
		}

		l := list.New(items, list.NewDefaultDelegate(), 0, 14)
		l.Title = "What would you like to do?"
		l.SetShowStatusBar(false)
		l.SetFilteringEnabled(false)
		l.SetShowHelp(true)

		m := menuModel{list: l}
		p := tea.NewProgram(m)

		finalModel, err := p.Run()
		if err != nil {
			fmt.Println("Error running menu:", err)
			os.Exit(1)
		}

		if m, ok := finalModel.(menuModel); ok && m.choice != "" {
			handleMenuChoice(m.choice)
		}
	},
}

func handleMenuChoice(choice string) {
	switch choice {
	case "🔐 Login":
		loginCmd.Run(loginCmd, []string{})
	case "🏗️  Check Infrastructure":
		infraCmd.Run(infraCmd, []string{})
	case "➕ Create Resource":
		createCmd.Run(createCmd, []string{})
	case "💬 Chat":
		chatCmd.Run(chatCmd, []string{})
	case "⚙️  Settings":
		fmt.Println("\n⚙️  Settings configuration coming soon!")
	case "🚪 Exit":
		fmt.Println("\n👋 Goodbye!")
		os.Exit(0)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
