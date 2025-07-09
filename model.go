package main

import (
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
)

type Event struct {
	Time    string
	Message string
}

func (e Event) Lipglossed() string {
	return styles["eventTime"].Render(e.Time) + " " + styles["eventMessage"].Render(e.Message)
}

type model struct {
	fileName string

	viewportWidth  int
	viewportHeight int

	time      string
	textInput textinput.Model
	events    []Event
}

func InitModel(f string) {
	ti := textinput.New()
	ti.Placeholder = "Type your item"
	ti.Focus()
	ti.CharLimit = appWidth - 19
	ti.Width = appWidth - 19
	ti.Prompt = ""

	m = model{fileName: f}

	events, err := Load()
	if err != nil {
		panic(err)
	}

	m.time = time.Now().Format("Mon 02/01 15:04:05")
	m.textInput = ti
	m.events = events
}

func (m model) PrintHeader() string {
	return styles["header"].Render("Gday, sir! " + m.fileName)
}

func (m model) PrintContent() string {
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

	return styles["base"].Render(m.time+" "+m.textInput.View()) + "\n" + styles["base"].Render(underline) + "\n" + styles["msg"].Render(events)
}

func (m model) Print() string {
	return m.PrintHeader() + "\n\n" + m.PrintContent()
}
