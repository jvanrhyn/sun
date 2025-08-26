package ui

import "github.com/charmbracelet/lipgloss"

// Theme contains reusable styles for the TUI components.
// Use DefaultTheme to obtain a sane preset.
type Theme struct {
	Base   lipgloss.Style
	Header lipgloss.Style
	Footer lipgloss.Style
	Error  lipgloss.Style
}

// DefaultTheme returns a preset Theme with borders, header/footer emphasis,
// and an error style in red for visibility.
func DefaultTheme() Theme {
	base := lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("240"))
	header := lipgloss.NewStyle().Bold(true)
	footer := lipgloss.NewStyle().Faint(true)
	err := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	return Theme{Base: base, Header: header, Footer: footer, Error: err}
}
