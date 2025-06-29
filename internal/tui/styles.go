package tui

import "github.com/charmbracelet/lipgloss"

// styles holds all Lip Gloss style definitions for the TUI.
type styles struct {
	App      lipgloss.Style
	Title    lipgloss.Style
	InfoBox  lipgloss.Style
	ErrorBox lipgloss.Style
	Success  lipgloss.Style
	Error    lipgloss.Style
	Help     lipgloss.Style
	Choice   lipgloss.Style
}

// defaultStyles returns an opinionated set of default styles used by the UI.
func defaultStyles() styles {
	return styles{
		App: lipgloss.NewStyle().
			Margin(1, 2),

		Title: lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFDF5")).
			Background(lipgloss.Color("#5A56E0")).
			Padding(0, 1).
			Bold(true),

		InfoBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1, 2),

		ErrorBox: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("9")).
			Padding(1, 2),

		Success: lipgloss.NewStyle().
			Foreground(lipgloss.Color("42")),

		Error: lipgloss.NewStyle().
			Foreground(lipgloss.Color("9")),

		Help: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")),

		Choice: lipgloss.NewStyle().
			Foreground(lipgloss.Color("252")),
	}
}
