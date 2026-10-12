package tui

import (
	"context"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestLineEdit(t *testing.T) {
	if got, ok := lineEdit(tea.KeyCtrlU, "abc"); !ok || got != "" {
		t.Fatalf("ctrl+u = %q ok=%v, want \"\" ok=true", got, ok)
	}
	if got, ok := lineEdit(tea.KeyCtrlU, ""); !ok || got != "" {
		t.Fatalf("ctrl+u on empty = %q ok=%v", got, ok)
	}
	if got, ok := lineEdit(tea.KeyBackspace, "ab"); !ok || got != "a" {
		t.Fatalf("backspace = %q ok=%v, want \"a\"", got, ok)
	}
	if got, ok := lineEdit(tea.KeyCtrlH, "ab"); !ok || got != "a" {
		t.Fatalf("ctrl+h = %q ok=%v, want \"a\"", got, ok)
	}
	if _, ok := lineEdit(tea.KeyBackspace, ""); !ok {
		t.Fatal("backspace on empty must be consumed (no-op)")
	}
	if _, ok := lineEdit(tea.KeyRunes, "ab"); ok {
		t.Fatal("printable chars are not editing keys")
	}
}

// TestMenuCtrlU verifies ctrl+u clears the goto menu query.
func TestMenuCtrlU(t *testing.T) {
	m := NewMenu(map[string]*Resource{"a b": {Name: "a b"}}, nil)
	for _, r := range "abc" {
		m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	if m.query != "abc" {
		t.Fatalf("query = %q, want abc", m.query)
	}
	m.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	if m.query != "" {
		t.Fatalf("query = %q, want empty after ctrl+u", m.query)
	}
}

// TestPromptCtrlU verifies ctrl+u clears the `:` prompt buffer and, in
// filter mode, resets the table filter with it.
func TestPromptCtrlU(t *testing.T) {
	res := &Resource{
		Name:  "instance server",
		Title: "instance server",
		Fetch: func(ctx context.Context) (listResult, error) {
			rows := []Row{{ID: "x", Cells: []string{"x"}}}

			return listResult{columns: []string{"ID"}, rows: rows}, nil
		},
	}
	a := &App{
		loc:       NewLocality("fr-par-1"),
		resources: map[string]*Resource{"instance server": res},
		stack:     []Screen{NewBrowser(res, nil)},
	}
	a.width, a.height = 80, 24
	a.sizeCurrent()

	// Filter mode: typing filters, ctrl+u clears buffer and filter.
	a.prompt = newFilterPrompt()
	for _, r := range "web" {
		a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	b := a.current().(*Browser)
	if got := b.Table().Filter; got != "web" {
		t.Fatalf("filter = %q, want web", got)
	}
	a.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	if a.prompt.Buffer != "" {
		t.Fatalf("buffer = %q, want empty after ctrl+u", a.prompt.Buffer)
	}
	if got := b.Table().Filter; got != "" {
		t.Fatalf("filter = %q, want reset after ctrl+u", got)
	}

	// Command mode: plain buffer clear.
	a.prompt = newCmdPrompt()
	for _, r := range "ip" {
		a.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}})
	}
	a.Update(tea.KeyMsg{Type: tea.KeyCtrlU})
	if a.prompt.Buffer != "" {
		t.Fatalf("buffer = %q, want empty after ctrl+u", a.prompt.Buffer)
	}
}
