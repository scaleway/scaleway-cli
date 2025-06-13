package rdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func engineSettingsCommand() *core.Command {
	type engineSettingsArgs struct {
		Name    string
		Version string
		Region  scw.Region
	}

	return &core.Command{
		Short:     `List available settings from an engine.`,
		Long:      `List available settings from an engine.`,
		Namespace: "rdb",
		Resource:  "engine",
		Verb:      "settings",
		ArgsType:  reflect.TypeOf(engineSettingsArgs{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*engineSettingsArgs)
			api := rdb.NewAPI(core.ExtractClient(ctx))

			databaseEnginesRequest := &rdb.ListDatabaseEnginesRequest{
				Region:  args.Region,
				Name:    &args.Name,
				Version: &args.Version,
			}

			engines, err := api.ListDatabaseEngines(databaseEnginesRequest)
			if err != nil {
				return nil, err
			}

			var responseEngineSettings []*rdb.EngineSetting
			for _, e := range engines.Engines {
				for _, ev := range e.Versions {
					if ev.Version == args.Version {
						responseEngineSettings = ev.AvailableSettings
					}
				}
			}

			return responseEngineSettings, nil
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "name",
				Short:    `The name of your engine where you want list the available settings.`,
				Required: true,
			},
			{
				Name:     "version",
				Required: true,
				Short:    "The version of the engine.",
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Examples: []*core.Example{
			{
				Short:    "List Engine Settings",
				ArgsJSON: `{"name": "MySQL", "version": "8"}`,
			},
		},
	}
}
