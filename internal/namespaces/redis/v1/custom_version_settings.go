package redis

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func redisVersionSettingsCommand() *core.Command {
	type versionSettingsArgs struct {
		Version string
		Zone    scw.Zone
	}

	return &core.Command{
		Short:     "List available settings from a Redis™ version",
		Long:      "List available settings from a Redis™ version.",
		Namespace: "redis",
		Resource:  "version",
		Verb:      "settings",
		ArgsType:  reflect.TypeOf(versionSettingsArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "version",
				Short:    "Redis™ engine version",
				Required: true,
			},
			core.ZoneArgSpec(
				scw.ZoneFrPar1,
				scw.ZoneFrPar2,
				scw.ZoneNlAms1,
				scw.ZoneNlAms2,
				scw.ZonePlWaw1,
				scw.ZonePlWaw2,
			),
		},
		Examples: []*core.Example{
			{
				Short:    "List settings for Redis™ 7.2.11",
				ArgsJSON: `{"version": "7.2.11"}`,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*versionSettingsArgs)
			api := redis.NewAPI(core.ExtractClient(ctx))

			resp, err := api.ListClusterVersions(&redis.ListClusterVersionsRequest{
				Zone:    args.Zone,
				Version: &args.Version,
			})
			if err != nil {
				return nil, err
			}

			for _, version := range resp.Versions {
				if version.Version == args.Version {
					return version.AvailableSettings, nil
				}
			}

			return []*redis.AvailableClusterSetting{}, nil
		},
	}
}
