package main

import (
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	gloss "github.com/charmbracelet/lipgloss"
	"github.com/matthewjwhite/ansibank/db"
)

// Define style "constants" - technically not constants, const not allowed for these.
var (
	cursorStyle = gloss.NewStyle().
			Bold(true).
			Foreground(gloss.Color("#A832A4"))
	cursorPoint = cursorStyle.Render(">")
	pathStyle   = gloss.NewStyle().
			Bold(true).
			Foreground(gloss.Color("#00a2ff"))
)

// Based on https://github.com/charmbracelet/bubbletea/blob/master/tutorials/basics.
// This tracks the state of the TUI, and maintains pointers to key data.
type listModel struct {
	choices []db.PathTime
	cursor  int
	db      db.DB
}

// Init is the initialization method expected for the bubbletea model. Since we have
// no initialization to do, it simply returns nil.
func (m listModel) Init() tea.Cmd {
	return nil
}

// Update modifies the model depending on user interaction.
func (m listModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up":
			if m.cursor > 0 {
				m.cursor--
			}

		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}

		case "enter":
			// Print output.
			output, err := m.db.GetOutput(m.choices[m.cursor])
			if err != nil {
				// Will not cause Start to fail, will quit normally. TUI still
				// did its job, despite failure internally. Not clear how to
				// propagate failure anyway through tea.
				fmt.Fprintf(os.Stderr, "\nFailed to retrieve output: %s\n", err)
			} else {
				// Output will be set if error is not nil.
				fmt.Print(output)
			}

			return m, tea.Quit
		}
	}

	return m, nil
}

// View renders the TUI.
func (m listModel) View() string {
	s := "Which playbook result would you like to view?\n\n"

	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = cursorPoint
		}

		choiceRender := pathStyle.Render(choice.Path) + " " +
			choice.StartTime.Local().Format(time.Stamp)
		s += fmt.Sprintf("%s %s\n", cursor, choiceRender)
	}

	s += "\nPress Ctrl+C to quit.\n"

	return s
}
