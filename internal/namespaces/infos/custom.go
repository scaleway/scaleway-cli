package infos

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		infosRoot(),
	)
}

type info struct {
	Key    string
	Value  string
	Origin string
}

func infosRoot() *core.Command {
	return &core.Command{
		Short:     `Send feedback to the Scaleway CLI Team!`,
		Namespace: "infos",
		ArgsType:  reflect.TypeOf(struct{}{}),
		ArgSpecs:  core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			return []info{
				{
					Key:   "config_path",
					Value: "~/.config/scw/config.yaml",
				},
				defaultRegion(ctx),
				defaultZone(ctx),
				defaultOrganizationId(ctx),
				accessKey(ctx),
				secretKey(ctx),
				profile(ctx),
			}, nil
		},
	}
}

func defaultRegion(ctx context.Context) info {
	info := info{Key: "default_region"}
	client := core.ExtractClient(ctx)
	defaultRegion, exists := client.GetDefaultRegion()
	if exists {
		info.Value = defaultRegion.String()
		info.Origin = "SCW_SECRET_KEY env"
	}
	return info
}

func defaultZone(ctx context.Context) info {
	info := info{Key: "default_zone"}
	client := core.ExtractClient(ctx)
	defaultZone, exists := client.GetDefaultZone()
	if exists {
		info.Value = defaultZone.String()
		info.Origin = "config_file"
	}
	return info
}

func defaultOrganizationId(ctx context.Context) info {
	info := info{Key: "default_organization_id"}
	client := core.ExtractClient(ctx)
	defaultOrganizationId, exists := client.GetDefaultOrganizationID()
	if exists {
		info.Value = defaultOrganizationId
		info.Origin = "config_file"
	}
	return info
}

func accessKey(ctx context.Context) info {
	info := info{Key: "access_key"}
	client := core.ExtractClient(ctx)
	aK, exists := client.GetAccessKey()
	if exists {
		info.Value = aK
		info.Origin = "SCW_ACCESS_KEY env"
	}
	return info
}

func secretKey(ctx context.Context) info {
	info := info{Key: "secret_key"}
	client := core.ExtractClient(ctx)
	sK, exists := client.GetSecretKey()
	if exists {
		info.Value = sK
		info.Origin = "config_file"
	}
	return info
}

func profile(ctx context.Context) info {
	return info{
		Key:    "profile",
		Value:  "default",
		Origin: "flag | SCW_PROFILE | config_file",
	}
}
