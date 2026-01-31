package wizard

import (
	"fmt"

	"accio/internal/tui/common"
	"accio/internal/tui/styles"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Step int

const (
	StepCloud Step = iota
	StepCapability
	StepDetails
	StepReview
)

type Model struct {
	step Step

	// Selections
	selectedCloud      string
	selectedCapability string
	resourceName       string

	// Components
	cloudList      list.Model
	capabilityList list.Model
	nameInput      textinput.Model

	// Dimensions
	width, height int
}

type item struct {
	title, desc string
	val         string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

func New() Model {
	// Cloud List
	clouds := []list.Item{
		item{title: "AWS", desc: "Amazon Web Services", val: "aws"},
		item{title: "GCP", desc: "Google Cloud Platform", val: "gcp"},
		item{title: "Azure", desc: "Microsoft Azure", val: "azure"},
	}
	cl := list.New(clouds, list.NewDefaultDelegate(), 0, 0)
	cl.Title = "Select Cloud Provider"
	cl.SetShowHelp(false)

	// Capability List
	caps := []list.Item{
		item{title: "Kubernetes Cluster", desc: "Managed K8s (EKS/GKE/AKS)", val: "kubernetes"},
		item{title: "Database", desc: "Managed SQL/NoSQL", val: "database"},
		item{title: "Compute Instance", desc: "Virtual Machine", val: "compute"},
	}
	cpl := list.New(caps, list.NewDefaultDelegate(), 0, 0)
	cpl.Title = "Select Capability"
	cpl.SetShowHelp(false)

	// Name Input
	ni := textinput.New()
	ni.Placeholder = "resource-name"
	ni.CharLimit = 32
	ni.Width = 30

	return Model{
		step:           StepCloud,
		cloudList:      cl,
		capabilityList: cpl,
		nameInput:      ni,
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		m.cloudList.SetWidth(msg.Width)
		m.capabilityList.SetWidth(msg.Width)

	case tea.KeyMsg:
		if msg.String() == "esc" {
			// Go back or cancel
			if m.step > StepCloud {
				m.step--
				return m, nil
			}
			return m, func() tea.Msg { return common.SwitchScreenMsg(common.ScreenDashboard) }
		}

		// Handle Enter based on step
		if msg.String() == "enter" {
			switch m.step {
			case StepCloud:
				i := m.cloudList.SelectedItem().(item)
				m.selectedCloud = i.val
				m.step++
				return m, nil
			case StepCapability:
				i := m.capabilityList.SelectedItem().(item)
				m.selectedCapability = i.val
				m.step++
				m.nameInput.Focus()
				return m, textinput.Blink
			case StepDetails:
				m.resourceName = m.nameInput.Value()
				m.step++
				return m, nil
			case StepReview:
				// Submit! (TODO: Call API)
				// For now, return to dashboard
				return m, func() tea.Msg { return common.SwitchScreenMsg(common.ScreenDashboard) }
			}
		}
	}

	// Forward Update to active component
	switch m.step {
	case StepCloud:
		m.cloudList, cmd = m.cloudList.Update(msg)
	case StepCapability:
		m.capabilityList, cmd = m.capabilityList.Update(msg)
	case StepDetails:
		m.nameInput, cmd = m.nameInput.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	var content string
	switch m.step {
	case StepCloud:
		content = m.cloudList.View()
	case StepCapability:
		content = m.capabilityList.View()
	case StepDetails:
		content = fmt.Sprintf(
			"Enter Resource Name:\n\n%s\n\n(Enter to continue, Esc to back)",
			m.nameInput.View(),
		)
	case StepReview:
		content = fmt.Sprintf(
			"Review Infra Intent:\n\nCloud: %s\nCapability: %s\nName: %s\n\n(Enter to Submit, Esc to back)",
			m.selectedCloud, m.selectedCapability, m.resourceName,
		)
	}

	return styles.BoxStyle.Render(content)
}
