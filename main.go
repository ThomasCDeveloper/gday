package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	appWidth = 80
)

var (
	m model
)

func (m model) quit() (tea.Model, tea.Cmd) {
	return m, tea.Quit
}

type tickMsg time.Time

func tickCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*200, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), textinput.Blink)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m.quit()
		case "enter":
			if m.textInput.Value() != "" {
				m.events = append(m.events, Event{Time: m.time, Message: m.textInput.Value()})
				m.textInput.Reset()
				if err := m.Save(); err != nil {
					panic(err)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.viewportWidth = msg.Width
		m.viewportHeight = msg.Height

	case tickMsg:
		m.time = time.Now().Format("Mon 02/01 15:04:05")
		return m, tickCmd()
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)

	return m, cmd
}

func (m model) View() string {
	if m.viewportWidth < appWidth || m.viewportHeight < 6 {
		return lipgloss.Place(
			m.viewportWidth,
			m.viewportHeight,
			lipgloss.Center,
			0,
			"Terminal too small")
	}

	page := m.Print()

	leftMargin := (m.viewportWidth - appWidth) / 2
	return lipgloss.NewStyle().MarginLeft(leftMargin).Render(page)
}

func main() {
	fileName := time.Now().Format("2006_01_02")
	if len(os.Args) > 1 {
		fileName = strings.Join(os.Args[1:], "_")
	}

	InitStyles()
	InitModel(fileName)

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
