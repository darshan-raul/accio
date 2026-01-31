package dashboard

import (
	"accio/internal/tui/common"
	"accio/internal/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type item struct {
	title, desc string
	screen      common.Screen
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type Model struct {
	list          list.Model
	width, height int
}

func New() Model {
	items := []list.Item{
		item{title: "Create Infrastructure", desc: "Provision new cloud resources via wizard", screen: common.ScreenWizard},
		item{title: "Observability", desc: "View stacks, status, and health", screen: common.ScreenObservability},
		item{title: "Settings", desc: "Configure API endpoint and context", screen: common.ScreenDashboard}, // Placeholder
	}

	l := list.New(items, list.NewDefaultDelegate(), 0, 0)
	l.Title = "Accio Platform"
	l.SetShowHelp(false)

	return Model{list: l}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		// Keep the list slightly smaller than the full width to look nice in the box
		m.list.SetWidth(msg.Width - 4)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "enter":
			i, ok := m.list.SelectedItem().(item)
			if ok {
				return m, func() tea.Msg {
					return common.SwitchScreenMsg(i.screen)
				}
			}
		}
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	doc := styles.BoxStyle.Render(m.list.View())
	if m.width > 0 && m.height > 0 {
		return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, doc)
	}
	return doc
}
