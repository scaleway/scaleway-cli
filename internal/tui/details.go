package tui

import (
	"fmt"
	"net"
	"reflect"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

var (
	timeType  = reflect.TypeFor[time.Time]()
	netIPType = reflect.TypeFor[net.IP]()
)

// Details is the key/value "describe" view of a single object. It is the
// k9s describe live-view.
type Details struct {
	title  string
	lines  []string
	width  int
	height int
	sel    int // selected line
	off    int // scroll offset (top line)
	// matches holds the line indices matching the last `/` search,
	// match the current one (less-style n/p navigation).
	matches []int
	match   int
	keys    *Keymap
}

// NewDetails builds a Details view from a (SDK) object via reflection.
func NewDetails(title string, obj any, keys *Keymap) *Details {
	d := &Details{title: title, keys: keys}
	var lines []string
	d.appendLines(reflect.ValueOf(obj), "", &lines)
	d.lines = lines

	return d
}

// Name implements Screen.
func (d *Details) Name() string {
	return d.title
}

// Hints implements Screen. The keys come from the keymap so the hints
// track cli.yaml overrides.
func (d *Details) Hints() []string {
	k := d.keys

	return []string{
		hint(k.Key(ActionFilter), "find"),
		hint(k.Key(ActionNext)+"/"+k.Key(ActionPrev), "next/prev"),
		hint(k.Key(ActionCopy), "copy"),
		hint(k.Key(ActionDown)+"/"+k.Key(ActionUp), "move"),
		hint(k.Key(ActionProfile), "profile"),
		hint(k.Key(ActionProject), "project"),
		hint(k.Key(ActionBack), "back"),
		hint(k.Key(ActionQuit), "quit"),
	}
}

// Size implements Screen.
func (d *Details) Size(width, height int) {
	d.width, d.height = width, height
}

// Init implements Screen.
func (d *Details) Init() tea.Cmd {
	return nil
}

// Update implements Screen.
func (d *Details) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch d.keys.Action(keyOf(km)) {
		case ActionUp:
			d.move(-1)
		case ActionDown:
			d.move(1)
		case ActionPageUp:
			d.move(-d.height)
		case ActionPageDown:
			d.move(d.height)
		case ActionHome:
			d.sel, d.off = 0, 0
		case ActionEnd:
			d.sel = max(0, len(d.lines)-1)
			d.off = max(0, len(d.lines)-d.height)
		case ActionNext:
			if cmd := d.nextMatch(1); cmd != nil {
				return d, cmd
			}
		case ActionPrev:
			if cmd := d.nextMatch(-1); cmd != nil {
				return d, cmd
			}
		case ActionCopy:
			if key, value := d.selected(); key != "" {
				return d, func() tea.Msg { return copyMsg{value: value, label: key} }
			}
		}
	}

	return d, nil
}

// SearchAndCopy finds all attributes whose key contains the query
// (case-insensitive), selects the first, and returns a cmd copying its
// value to the clipboard. Subsequent n/p keys walk the remaining
// matches (less-style, wrapping around).
func (d *Details) SearchAndCopy(query string) tea.Cmd {
	query = strings.ToLower(strings.TrimSpace(query))
	if query == "" {
		return nil
	}
	d.matches = nil
	for i, l := range d.lines {
		key, _, _ := strings.Cut(l, "\t")
		if strings.Contains(strings.ToLower(key), query) {
			d.matches = append(d.matches, i)
		}
	}
	if len(d.matches) == 0 {
		return func() tea.Msg { return flashMsg{msg: "no attribute: " + query} }
	}
	d.match = 0

	return d.copyMatch()
}

// View implements Screen.
func (d *Details) View() string {
	if d.width <= 0 {
		return ""
	}

	if len(d.lines) == 0 {
		lines := []string{dimStyle.Render("no data")}
		for len(lines) < d.height {
			lines = append(lines, "")
		}

		return strings.Join(lines, "\n")
	}
	start := d.off
	if start > len(d.lines)-d.height {
		start = max(0, len(d.lines)-d.height)
	}
	visible := d.lines[start:min(start+d.height, len(d.lines))]
	keyW := 20
	for _, l := range visible {
		if k, _, _ := strings.Cut(l, "\t"); len(k) > keyW {
			keyW = len(k)
		}
	}
	keyW = min(keyW, 32)
	valueW := d.width - cursorGutter - keyW - 2
	out := make([]string, 0, d.height)
	for i, l := range visible {
		k, v, _ := strings.Cut(l, "\t")
		if start+i == d.sel {
			out = append(out, cursorStyle.Render("> "+pad(k, keyW)+"  "+truncate(v, valueW)))
		} else {
			out = append(
				out,
				"  "+detailKeyStyle.Render(
					pad(k, keyW),
				)+"  "+detailValueStyle.Render(
					truncate(v, valueW),
				),
			)
		}
	}
	for len(out) < d.height {
		out = append(out, "")
	}

	return strings.Join(out, "\n")
}

// copyMatch selects the current match and returns a cmd copying its
// value to the clipboard.
func (d *Details) copyMatch() tea.Cmd {
	line := d.matches[d.match]
	d.sel = line
	d.off = max(0, min(line, max(0, len(d.lines)-d.height)))
	key, value, _ := strings.Cut(d.lines[line], "\t")

	return func() tea.Msg { return copyMsg{value: value, label: key} }
}

// nextMatch moves to the next/previous search match, wrapping around
// (less n/p), and returns a cmd copying the value at the new position.
// Before any search it returns a hint flash instead.
func (d *Details) nextMatch(delta int) tea.Cmd {
	if len(d.matches) == 0 {
		return func() tea.Msg { return flashMsg{msg: "no search yet — press /"} }
	}
	d.match = (d.match + delta + len(d.matches)) % len(d.matches)

	return d.copyMatch()
}

func (d *Details) move(delta int) {
	d.sel = min(max(0, d.sel+delta), max(0, len(d.lines)-1))
	if d.sel < d.off {
		d.off = d.sel
	}
	if d.sel >= d.off+d.height {
		d.off = d.sel - d.height + 1
	}
	d.off = max(0, d.off)
}

func (d *Details) selected() (string, string) {
	if d.sel < 0 || d.sel >= len(d.lines) {
		return "", ""
	}
	k, v, _ := strings.Cut(d.lines[d.sel], "\t")

	return k, v
}

// appendLines walks a struct and appends "key\tvalue" lines. Nested
// structs are flattened with dotted keys.
func (d *Details) appendLines(v reflect.Value, prefix string, lines *[]string) {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return
		}
		v = v.Elem()
	}
	if v.Type() == timeType {
		if prefix != "" {
			*lines = append(*lines, prefix+"\t"+v.Interface().(time.Time).Format(time.RFC3339))
		}

		return
	}
	if v.Kind() != reflect.Struct {
		if prefix != "" {
			*lines = append(*lines, prefix+"\t"+formatValue(v))
		}

		return
	}
	t := v.Type()
	for i := range t.NumField() {
		f := t.Field(i)
		if !f.IsExported() {
			continue
		}
		name := goFieldName(f.Name)
		if prefix != "" {
			name = prefix + "." + name
		}
		d.appendLines(v.Field(i), name, lines)
	}
}

// goFieldName renders a Go field name as a snake_case key, keeping
// acronyms like IP/ID intact (strcase.ToSnake would emit "public_i_ps").
func goFieldName(name string) string {
	s := strings.ReplaceAll(strcase.ToSnake(name), "_i_", "_i")
	if rest, ok := strings.CutPrefix(s, "i_"); ok {
		s = "i" + rest
	}

	return s
}

// formatValue renders a leaf value.
// ponytail: collections are summarized (first 5 items), no per-element
// drill-down; add a nested list view if needed.
func formatValue(v reflect.Value) string {
	for v.Kind() == reflect.Pointer || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return "-"
		}
		v = v.Elem()
	}
	switch {
	case v.Type() == timeType:
		return v.Interface().(time.Time).Format(time.RFC3339)
	case v.Type() == netIPType:
		return v.Interface().(net.IP).String()
	case v.Kind() == reflect.String:
		s := v.String()
		if s == "" {
			return "-"
		}

		return s
	case v.Kind() == reflect.Slice, v.Kind() == reflect.Array:
		if v.Len() == 0 {
			return "-"
		}
		parts := make([]string, 0, min(v.Len(), 5))
		for i := 0; i < v.Len() && i < 5; i++ {
			parts = append(parts, formatValue(v.Index(i)))
		}
		s := strings.Join(parts, ", ")
		if v.Len() > 5 {
			s += fmt.Sprintf(" (%d total)", v.Len())
		}

		return s
	case v.Kind() == reflect.Map:
		if v.Len() == 0 {
			return "-"
		}
		parts := make([]string, 0, min(v.Len(), 5))
		it := v.MapRange()
		n := 0
		for it.Next() && n < 5 {
			parts = append(parts, it.Key().String()+"="+formatValue(it.Value()))
			n++
		}
		s := strings.Join(parts, ", ")
		if v.Len() > 5 {
			s += fmt.Sprintf(" (%d total)", v.Len())
		}

		return s
	default:
		return fmt.Sprintf("%v", v.Interface())
	}
}
