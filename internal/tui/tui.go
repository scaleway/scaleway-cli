package tui

import (
	"context"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Run starts the TUI and blocks until the user quits (q or ctrl+c).
// meta owns the API client (profile switches reload it through it),
// keys is the effective keymap (nil: defaults), commands the user
// `:` commands from cli.yaml, confirmDelete toggles the delete
// confirmation modal (tui.delete_confirm), loc the session zone/region
// (nil: fr-par-1).
func Run(
	ctx context.Context,
	client *scw.Client,
	meta *core.Meta,
	resources map[string]*Resource,
	keys *Keymap,
	commands map[string]string,
	confirmDelete bool,
	loc *Locality,
) error {
	if loc == nil {
		loc = NewLocality(scw.ZoneFrPar1)
	}
	app := NewApp(ctx, client, meta, resources, keys, commands, loc)
	app.confirmDelete = confirmDelete
	p := tea.NewProgram(app, tea.WithContext(ctx), tea.WithAltScreen())
	_, err := p.Run()

	return err
}
