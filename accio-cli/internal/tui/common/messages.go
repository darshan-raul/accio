package common

type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenWizard
	ScreenObservability
)

type SwitchScreenMsg Screen
