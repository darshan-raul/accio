package tui

import (
	"fmt"
	"os"

	"accio/internal/api"
	"accio/internal/config"
	"accio/internal/tui/common"
	"accio/internal/tui/screens/dashboard"
	"accio/internal/tui/screens/observability"
	"accio/internal/tui/screens/wizard"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	cfg    *config.Config
	client api.Client

	currentScreen common.Screen
	dashboard     tea.Model
	wizard        tea.Model
	observability tea.Model
	// Add other screens here
}

func NewModel(cfg *config.Config, client api.Client) Model {
	return Model{
		cfg:           cfg,
		client:        client,
		currentScreen: common.ScreenDashboard,
		dashboard:     dashboard.New(),
		wizard:        wizard.New(),
		observability: observability.New(client),
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.dashboard.Init(),
		tea.EnterAltScreen,
	)
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			// Simple navigation for now
			m.currentScreen = (m.currentScreen + 1) % 3
		}
	}

	// Handle global messages or screen switching
	switch msg := msg.(type) {
	case common.SwitchScreenMsg:
		m.currentScreen = common.Screen(msg)
		// Re-init the target screen if needed?
		// For now just switch.
		return m, nil
	}

	// Forward messages to the current screen
	switch m.currentScreen {
	case common.ScreenDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
		cmds = append(cmds, cmd)
	case common.ScreenWizard:
		m.wizard, cmd = m.wizard.Update(msg)
		cmds = append(cmds, cmd)
	case common.ScreenObservability:
		m.observability, cmd = m.observability.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	switch m.currentScreen {
	case common.ScreenDashboard:
		return m.dashboard.View()
	case common.ScreenWizard:
		return m.wizard.View()
	case common.ScreenObservability:
		return m.observability.View()
	default:
		return "Unknown Screen"
	}
}

func Start(cfg *config.Config, client api.Client) {
	p := tea.NewProgram(NewModel(cfg, client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}

// StartObservability launches only the observability screen
func StartObservability(cfg *config.Config, client api.Client) {
	p := tea.NewProgram(observability.New(client))
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}

// StartWizard launches only the wizard screen
func StartWizard(cfg *config.Config, client api.Client) {
	p := tea.NewProgram(wizard.New())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
