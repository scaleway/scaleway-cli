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

Edit the selected resource from a list view:
    e       open it in your $EDITOR as YAML (a comment documents each
            field); save and quit to apply the update

Delete the selected resource from a list view:
    d       delete it (asks for confirmation; set
            tui.delete_confirm: false in cli.yaml to skip it)

Keybindings and the : commands can be customized in the tui section of the CLI config file (cli.yaml):

    tui:
        keybindings:
            refresh: R
        commands:
            mydb: rdb instance
        delete_confirm: false

List every keybinding with its effective key: scw tui keybinding list`,
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
			// Keybindings, `:` commands and the delete-confirmation
			// toggle, overridable in the tui section of cli.yaml.
			keys := tui.DefaultKeymap()
			commands := map[string]string{}
			confirmDelete := true
			if cfg := core.ExtractCliConfig(ctx); cfg != nil && cfg.Tui != nil {
				keys = tui.NewKeymap(cfg.Tui.KeyBindings)
				commands = cfg.Tui.Commands
				if cfg.Tui.DeleteConfirm != nil {
					confirmDelete = *cfg.Tui.DeleteConfirm
				}
			}

			return nil, tui.Run(ctx, client, meta, resources, keys, commands, confirmDelete, loc)
		},
	}
}

// GetCommands returns the tui commands.
func GetCommands() *core.Commands {
	return core.NewCommands(NewTUICommand(), keybindingListCommand())
}
