package main

import "github.com/charmbracelet/lipgloss"

var styles map[string]lipgloss.Style = map[string]lipgloss.Style{}

func InitStyles() {
	base := lipgloss.NewStyle().Width(appWidth)
	styles["base"] = base

	styles["header"] = base.MarginTop(1).Bold(true)
	styles["msg"] = base
	styles["eventTime"] = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	styles["eventMessage"] = lipgloss.NewStyle()
}
