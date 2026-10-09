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
// `:` commands from cli.yaml, loc the session zone/region (nil:
// fr-par-1).
func Run(
	ctx context.Context,
	client *scw.Client,
	meta *core.Meta,
	resources map[string]*Resource,
	keys *Keymap,
	commands map[string]string,
	loc *Locality,
) error {
	if loc == nil {
		loc = NewLocality(scw.ZoneFrPar1)
	}
	p := tea.NewProgram(
		NewApp(ctx, client, meta, resources, keys, commands, loc),
		tea.WithContext(ctx),
		tea.WithAltScreen(),
	)
	_, err := p.Run()

	return err
}
