package main

import (
	"encoding/json"
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
	fileName = ""

	baseStyle = lipgloss.NewStyle().Width(appWidth)

	headerStyle = baseStyle.
			MarginTop(1).
			Bold(true)

	msgStyle = baseStyle

	eventTimeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))

	eventMessageStyle = lipgloss.NewStyle()
)

type model struct {
	viewportWidth  int
	viewportHeight int

	time string

	textInput textinput.Model
	events    []Event
}

type Event struct {
	Time    string
	Message string
}

func (e Event) Lipglossed() string {
	return eventTimeStyle.Render(e.Time) + " " + eventMessageStyle.Render(e.Message)
}

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

func (m model) Save() error {
	file, err := os.Create(fileName + ".json")
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")

	return encoder.Encode(m.events)
}

func Load() ([]Event, error) {
	empty := []Event{}
	file, err := os.Open(fileName + ".json")
	if err != nil {
		if os.IsNotExist(err) {
			return empty, nil
		}
		return empty, err
	}
	defer file.Close()

	var loaded []Event
	decoder := json.NewDecoder(file)
	if err := decoder.Decode(&loaded); err != nil {
		return empty, err
	}

	return loaded, nil
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

	header := headerStyle.Render("Gday, sir! " + fileName)

	visible := min(25, m.viewportHeight-6)
	if visible > len(m.events) {
		visible = len(m.events)
	}
	recent := m.events[len(m.events)-visible:]

	events := ""
	for _, e := range recent {
		events = e.Lipglossed() + "\n" + events
	}

	underline := strings.Repeat("─", appWidth)
	content := baseStyle.Render(m.time+" "+m.textInput.View()) + "\n" + baseStyle.Render(underline) + "\n" + msgStyle.Render(events)

	page := header + "\n\n" + content

	leftMargin := (m.viewportWidth - appWidth) / 2
	return lipgloss.NewStyle().MarginLeft(leftMargin).Render(page)
}

func main() {
	fileName = time.Now().Format("2006_01_02")
	if len(os.Args) > 1 {
		fileName = strings.Join(os.Args[1:], "_")
	}

	ti := textinput.New()
	ti.Placeholder = "Type your item"
	ti.Focus()
	ti.CharLimit = appWidth - 19
	ti.Width = appWidth - 19
	ti.Prompt = ""

	m := model{
		time:      time.Now().Format("Mon 02/01 15:04:05"),
		events:    []Event{},
		textInput: ti,
	}

	events, err := Load()
	if err != nil {
		panic(err)
	}
	m.events = events

	if _, err := tea.NewProgram(m, tea.WithAltScreen()).Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
