package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Actions are the TUI keybinding names. They are the keys of the
// `tui.keybindings` section of cli.yaml (action -> key) and select what
// a pressed key does.
const (
	ActionQuit     = "quit"      // q: quit the TUI
	ActionMenu     = "menu"      // m, g: open the goto menu
	ActionCommand  = "command"   // :: open the command prompt
	ActionFilter   = "filter"    // /: filter (browser) / find (details) prompt
	ActionRefresh  = "refresh"   // r: re-fetch the current screen
	ActionBack     = "back"      // esc: back (pop the navigation stack)
	ActionOpen     = "open"      // enter: open the selected row / menu item
	ActionCopy     = "copy"      // y: copy the selected id / attribute
	ActionZone     = "zone"      // ctrl+z: pick the session zone
	ActionRegion   = "region"    // ctrl+r: pick the session region
	ActionProfile  = "profile"   // ctrl+p: pick the scw profile
	ActionProject  = "project"   // ctrl+e: pick the session project
	ActionNext     = "next"      // n: next search match (details)
	ActionPrev     = "prev"      // p: previous search match (details)
	ActionUp       = "up"        // k, up
	ActionDown     = "down"      // j, down
	ActionPageUp   = "page-up"   // pgup
	ActionPageDown = "page-down" // pgdown
	ActionHome     = "home"      // home
	ActionEnd      = "end"       // end
)

// defaultKeys lists the default key(s) of each action, in hint order
// (the first key is the one shown in the header hints).
var defaultKeys = map[string][]string{
	ActionQuit:     {"q"},
	ActionMenu:     {"m", "g"},
	ActionCommand:  {":"},
	ActionFilter:   {"/"},
	ActionRefresh:  {"r"},
	ActionBack:     {"esc"},
	ActionOpen:     {"enter"},
	ActionCopy:     {"y"},
	ActionNext:     {"n"},
	ActionPrev:     {"p"},
	ActionUp:       {"k", "up"},
	ActionDown:     {"j", "down"},
	ActionPageUp:   {"pgup"},
	ActionPageDown: {"pgdown"},
	ActionHome:     {"home"},
	ActionEnd:      {"end"},
	ActionZone:     {"ctrl+z"},
	ActionRegion:   {"ctrl+r"},
	ActionProfile:  {"ctrl+p"},
	ActionProject:  {"ctrl+e"},
}

// Keymap maps a pressed key to its action. It is built from the default
// bindings, each action overridable from cli.yaml (`tui.keybindings`:
// action -> key). An override replaces the default key(s) of the action;
// the other actions keep theirs. Unknown actions or empty keys are
// ignored.
type Keymap struct {
	keys    map[string]string // key -> action
	actions map[string][]string
}

var defaultKeymap = buildKeymap(nil)

// DefaultKeymap returns the keymap without any override.
func DefaultKeymap() *Keymap {
	return defaultKeymap
}

// NewKeymap builds the effective keymap from overrides (action -> key).
func NewKeymap(overrides map[string]string) *Keymap {
	return buildKeymap(overrides)
}

func buildKeymap(overrides map[string]string) *Keymap {
	keys := map[string]string{}
	actions := map[string][]string{}
	for action, defKeys := range defaultKeys {
		ks := defKeys
		if k, ok := overrides[action]; ok && k != "" {
			ks = []string{k}
		}
		actions[action] = ks
		for _, k := range ks {
			keys[k] = action
		}
	}

	return &Keymap{keys: keys, actions: actions}
}

// Action returns the action bound to key, "" if none. A nil keymap
// falls back to the defaults.
func (km *Keymap) Action(key string) string {
	if km == nil {
		return defaultKeymap.keys[key]
	}

	return km.keys[key]
}

// Key returns the first key bound to an action, for the header hints.
func (km *Keymap) Key(action string) string {
	if km == nil {
		km = defaultKeymap
	}
	if keys := km.actions[action]; len(keys) > 0 {
		return keys[0]
	}

	return "?"
}

// keyOf normalizes a tea.KeyMsg to the key spelling of the keymap
// (tea's own: a rune as-is, or "esc", "enter", "pgup", "ctrl+r", ...).
func keyOf(k tea.KeyMsg) string {
	return k.String()
}

// hint formats one header key hint ("[ enter ] open").
func hint(key, label string) string {
	return "[ " + key + " ] " + label
}
