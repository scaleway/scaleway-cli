package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestKeymapDefaults(t *testing.T) {
	km := DefaultKeymap()
	for key, action := range map[string]string{
		"q": ActionQuit, "m": ActionMenu, "g": ActionMenu,
		":": ActionCommand, "/": ActionFilter, "r": ActionRefresh,
		"esc": ActionBack, "enter": ActionOpen, "y": ActionCopy,
		"n": ActionNext, "p": ActionPrev,
		"k": ActionUp, "up": ActionUp,
		"j": ActionDown, "down": ActionDown,
		"pgup": ActionPageUp, "pgdown": ActionPageDown,
		"home": ActionHome, "end": ActionEnd,
	} {
		if got := km.Action(key); got != action {
			t.Errorf("Action(%q) = %q, want %q", key, got, action)
		}
	}
	if got := km.Key(ActionMenu); got != "m" {
		t.Errorf("Key(menu) = %q, want m", got)
	}
	if got := km.Key(ActionDown); got != "j" {
		t.Errorf("Key(down) = %q, want j", got)
	}
}

func TestKeymapOverrides(t *testing.T) {
	km := NewKeymap(map[string]string{
		"menu":    "t",
		"refresh": "R",
		"copy":    "",  // empty key: ignored, default kept
		"bogus":   "z", // unknown action: ignored
	})
	if got := km.Action("t"); got != ActionMenu {
		t.Errorf("Action(t) = %q, want menu", got)
	}
	if got := km.Action("R"); got != ActionRefresh {
		t.Errorf("Action(R) = %q, want refresh", got)
	}
	// Overridden actions lose their default keys.
	for _, key := range []string{"m", "g", "r"} {
		if got := km.Action(key); got != "" {
			t.Errorf("Action(%q) should be unbound, got %q", key, got)
		}
	}
	// Ignored entries keep the defaults.
	if got := km.Action("y"); got != ActionCopy {
		t.Errorf("Action(y) = %q, want copy (default kept)", got)
	}
	// Unset actions keep their defaults.
	if got := km.Action("q"); got != ActionQuit {
		t.Errorf("Action(q) = %q, want quit", got)
	}
	if got := km.Key(ActionMenu); got != "t" {
		t.Errorf("Key(menu) = %q, want t", got)
	}
}

func TestKeymapNilFallsBackToDefaults(t *testing.T) {
	var km *Keymap
	if got := km.Action("q"); got != ActionQuit {
		t.Fatalf("nil keymap Action(q) = %q, want quit", got)
	}
	if got := km.Key(ActionMenu); got != "m" {
		t.Fatalf("nil keymap Key(menu) = %q, want m", got)
	}
}

func TestKeyOf(t *testing.T) {
	cases := []struct {
		msg  tea.KeyMsg
		want string
	}{
		{tea.KeyMsg{Type: tea.KeyEscape}, "esc"},
		{tea.KeyMsg{Type: tea.KeyEnter}, "enter"},
		{tea.KeyMsg{Type: tea.KeyPgUp}, "pgup"},
		{tea.KeyMsg{Type: tea.KeyCtrlR}, "ctrl+r"},
		{tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("j")}, "j"},
	}
	for _, c := range cases {
		if got := keyOf(c.msg); got != c.want {
			t.Errorf("keyOf(%v) = %q, want %q", c.msg.Type, got, c.want)
		}
	}
}

func TestExecPromptUserCommands(t *testing.T) {
	res := &Resource{Name: "rdb instance", Title: "rdb instance"}
	a := &App{
		resources: map[string]*Resource{"rdb instance": res},
		commands:  map[string]string{"mydb": "rdb instance", "bye": "quit"},
		prompt:    &Prompt{Active: true, Mode: PromptCmd, Buffer: "mydb"},
	}
	a.execPrompt()
	b, ok := a.current().(*Browser)
	if !ok || b.res.Name != "rdb instance" {
		t.Fatalf("user command must goto the target resource, got %T", a.current())
	}
	// A user command can alias a built-in command.
	a.prompt.Buffer = "bye"
	if cmd := a.execPrompt(); cmd == nil {
		t.Fatal("alias to 'quit' must run the built-in quit")
	}
	// Unknown word: a flash, no navigation.
	a.prompt.Buffer = "nope"
	a.execPrompt()
	if a.flash != "unknown command: nope" {
		t.Fatalf("flash = %q, want 'unknown command: nope'", a.flash)
	}
}
