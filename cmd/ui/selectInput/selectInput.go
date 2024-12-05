// Package multiInput provides functions that
// help define and draw a multi-input step
package selectInput

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/mr-smith-org/mr/cmd/ui/textInput"
	"github.com/mr-smith-org/mr/cmd/ui/utils/program"
	"github.com/mr-smith-org/mr/cmd/ui/utils/steps"
	"github.com/mr-smith-org/mr/pkg/style"
)

// A Selection represents a choice made in a multiInput step
type Selection struct {
	Choice string
}

// Update changes the value of a Selection's Choice
func (s *Selection) Update(value string) {
	s.Choice = value
}

// A multiInput.model contains the data for the multiInput step.
//
// It has the required methods that make it a bubbletea.Model
type model struct {
	cursor   int
	choices  []steps.Item
	selected map[int]struct{}
	choice   *Selection
	header   string
	other    bool
	program  *program.Program
}

func (m model) Init() tea.Cmd {
	return nil
}

// InitialSelectInputModel initializes a multiInput step with
// the given data
func InitialSelectInputModel(choices []steps.Item, selection *Selection, header string, other bool, program *program.Program) model {
	m := model{
		choices:  choices,
		selected: make(map[int]struct{}),
		choice:   selection,
		header:   style.TitleStyle.Render(header),
		program:  program,
		other:    other,
	}
	return m
}

// Update is called when "things happen", it checks for
// important keystrokes to signal when to quit, change selection,
// and confirm the selection.
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.program.Exit = true
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			m.selected[m.cursor] = struct{}{}
			for selectedKey := range m.selected {
				m.choice.Update(m.choices[selectedKey].Value)
				m.cursor = selectedKey
			}
			return m, tea.Quit
		case "o":
			if m.other {
				textValue := &textInput.Output{}
				p := tea.NewProgram(textInput.InitialTextInputModel(textValue, "", "Enter another option", m.program))
				_, err := p.Run()
				if err != nil {
					style.ErrorPrint("error running program: " + err.Error())
					m.program.Exit = true
					return m, tea.Quit
				}
				m.choice.Update(textValue.Output)
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// View is called to draw the multiInput step
func (m model) View() string {
	s := m.header + "\n\n"
	for i, choice := range m.choices {
		cursor := " "
		if m.cursor == i {
			cursor = style.SelectedItemStyle.Render(">")
			choice.Label = style.SelectedItemStyle.Render(choice.Label)
			choice.Description = style.DescriptionStyle.Render(choice.Description)
			choice.Tags = style.TagsStyle.Render(choice.Tags)
		}

		checked := " "
		if _, ok := m.selected[i]; ok {
			checked = style.SelectedItemStyle.Render("*")
		}

		label := style.FocusedStyle.Render(choice.Label)
		description := style.DescriptionStyle.Render(choice.Description)
		tags := style.TagsStyle.Render(choice.Tags)

		s += fmt.Sprintf("%s [%s] %s%s%s\n\n", cursor, checked, label, description, tags)
	}
	if m.other {
		s += fmt.Sprintf("Press %s to text another option.\n", style.FocusedStyle.Render("o"))
	}
	s += fmt.Sprintf("Press %s to quit.\n", style.FocusedStyle.Render("q"))
	s += "\n"
	return s
}
