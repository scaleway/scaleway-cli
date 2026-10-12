//go:build !wasm

package tui

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/tui"
)

// keybindingListCommand returns the `scw tui keybinding list` command:
// every TUI keybinding with the tui.keybindings overrides from the CLI
// config file (cli.yaml) applied.
func keybindingListCommand() *core.Command {
	return &core.Command{
		Short:                "List the TUI keybindings",
		Namespace:            "tui",
		Resource:             "keybinding",
		Verb:                 "list",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(struct{}{}),
		Run: func(ctx context.Context, _ any) (any, error) {
			var overrides map[string]string
			if cfg := core.ExtractCliConfig(ctx); cfg != nil && cfg.Tui != nil {
				overrides = cfg.Tui.KeyBindings
			}

			return tui.NewKeymap(overrides).Bindings(), nil
		},
	}
}
