package ui

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6")).
			Padding(0, 1)

	DoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			Strikethrough(true)

	ActiveStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F8F8F2"))

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB86C")).
			Bold(true)

	HelpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4"))

	FilterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8BE9FD"))

	InputPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#BD93F9"))

	BorderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#8BE9FD")).
			Padding(1, 2)
)
