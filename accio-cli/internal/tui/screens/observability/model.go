package observability

import (
	"context"
	"fmt"

	"accio/internal/api"
	"accio/internal/models"
	"accio/internal/tui/styles"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Model struct {
	client api.Client
	table  table.Model
	err    error
}

func New(client api.Client) Model {
	columns := []table.Column{
		{Title: "Name", Width: 20},
		{Title: "Cloud", Width: 10},
		{Title: "Region", Width: 15},
		{Title: "Status", Width: 15},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	t.SetStyles(s)

	return Model{
		client: client,
		table:  t,
	}
}

func (m Model) Init() tea.Cmd {
	return m.fetchStacks
}

func (m Model) fetchStacks() tea.Msg {
	stacks, err := m.client.GetStacks(context.Background())
	if err != nil {
		return err
	}
	return stacks
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.table.SetWidth(msg.Width)

	case error:
		m.err = msg
		return m, nil

	case []models.Stack:
		rows := []table.Row{}
		for _, s := range msg {
			rows = append(rows, table.Row{
				s.Name,
				s.Cloud,
				s.Region,
				string(s.Status),
			})
		}
		m.table.SetRows(rows)
	}

	var cmd tea.Cmd
	m.table, cmd = m.table.Update(msg)
	return m, cmd
}

func (m Model) View() string {
	if m.err != nil {
		return styles.StatusInvalidStyle.Render(fmt.Sprintf("Error: %v", m.err))
	}
	return styles.BoxStyle.Render(
		lipgloss.JoinVertical(lipgloss.Left,
			styles.TitleStyle.Render("Cluster Observability"),
			m.table.View(),
		),
	)
}
