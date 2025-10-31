//go:build !wasm

package interactive

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type ListPrompt struct {
	// Prompt that will be printed when showing the list
	Prompt  string
	Choices []string
	// DefaultIndex is the element that will be selected when starting prompt
	DefaultIndex int

	cursor    int
	cancelled bool
}

func (m *ListPrompt) Init() tea.Cmd {
	return nil
}

func (m *ListPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		// Key is pressed
		switch msg.String() {
		case "ctrl+c", "q":
			m.cancelled = true

			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.Choices)-1 {
				m.cursor++
			}
		case "enter", " ":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m *ListPrompt) View() string {
	var builder strings.Builder
	builder.WriteString(m.Prompt + "\n\n")

	for i, choice := range m.Choices {
		if m.cursor == i {
			fmt.Fprintf(&builder, "> %s\n", choice)
		} else {
			builder.WriteString(choice + "\n")
		}
	}

	builder.WriteString("\nPress enter or space for select.\n")

	return builder.String()
}

// Execute start the prompt and return the selected index
func (m *ListPrompt) Execute(ctx context.Context) (int, error) {
	m.cursor = m.DefaultIndex

	opts := []tea.ProgramOption{
		tea.WithContext(ctx),
	}

	if hasMockedResponse(ctx) {
		opts = append(opts, tea.WithInput(&mockResponseReader{
			ctx:           ctx,
			defaultReader: os.Stdin,
		}))
		opts = append(opts, tea.WithOutput(bytes.NewBuffer([]byte{})))
	}

	p := tea.NewProgram(m, opts...)
	_, err := p.Run()
	if err != nil {
		return -1, fmt.Errorf("error running prompt: %w", err)
	}

	if m.cancelled {
		return -1, errors.New("prompt cancelled")
	}

	return m.cursor, nil
}
