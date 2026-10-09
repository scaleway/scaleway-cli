//go:build !wasm

package tui

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/tui"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type args struct {
	Zone scw.Zone
}

// NewTUICommand returns the `scw tui` command.
func NewTUICommand() *core.Command {
	return &core.Command{
		Short: "Launch the interactive TUI",
		Long: `Launch the interactive TUI to browse and monitor Scaleway resources.

Change the session zone, region, profile or project at any time:
    ctrl+z  pick the zone
    ctrl+r  pick the region
    ctrl+p  pick the profile
    ctrl+e  pick the project
The choice applies to the current session only (not persisted).

Keybindings and the : commands can be customized in the tui section of the CLI config file (cli.yaml):

    tui:
        keybindings:
            refresh: R
        commands:
            mydb: rdb instance`,
		Namespace: "tui",
		ArgsType:  reflect.TypeFor[args](),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
			),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*args)
			client := core.ExtractClient(ctx)
			meta := core.ExtractMeta(ctx)
			zone := args.Zone
			if zone == "" {
				zone, _ = client.GetDefaultZone()
			}
			if zone == "" {
				zone = scw.ZoneFrPar1
			}
			// The session locality (zone/region), changeable at runtime
			// with the zone/region pickers and shared with the fetches.
			loc := tui.NewLocality(zone)
			// Every `scw <ns> <res> list` command becomes a browsable
			// resource. The bootstrap meta carries the authenticated
			// client and the full command registry.
			resources := tui.ResourcesFromCommands(
				meta.Commands.GetAll(),
				meta,
				loc,
			)
			// Keybindings and `:` commands, overridable in the tui
			// section of cli.yaml.
			keys := tui.DefaultKeymap()
			var commands map[string]string
			if cfg := core.ExtractCliConfig(ctx); cfg != nil && cfg.Tui != nil {
				keys = tui.NewKeymap(cfg.Tui.KeyBindings)
				commands = cfg.Tui.Commands
			}

			return nil, tui.Run(ctx, client, meta, resources, keys, commands, loc)
		},
	}
}

// GetCommands returns the tui commands.
func GetCommands() *core.Commands {
	return core.NewCommands(NewTUICommand())
}
