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

type infoRequest struct {
	ShowSecret *bool
}

func infosRoot() *core.Command {
	return &core.Command{
		Short: `Get current config status`,
		// TODO status, infos ?
		Namespace: "status",
		ArgsType:  reflect.TypeOf(infoRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "show-secret",
				Short:    `Reveal secret`,
				Required: false,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			req := argsI.(*infoRequest)
			return []info{
				configPath(ctx),
				defaultRegion(ctx),
				defaultZone(ctx),
				defaultOrganizationId(ctx),
				accessKey(ctx),
				secretKey(ctx, req.ShowSecret),
				profile(ctx),
			}, nil
		},
	}
}

func configPath(ctx context.Context) info {
	return info{
		Key:   "config_path",
		Value: "~/.config/scw/config.yaml",
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

func HideSecretKey(k string) string {
	switch {
	case len(k) == 0:
		return ""
	case len(k) > 8:
		return k[0:8] + "-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	default:
		return "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
}

func secretKey(ctx context.Context, secrets *bool) info {
	info := info{Key: "secret_key"}
	client := core.ExtractClient(ctx)
	sK, exists := client.GetSecretKey()
	if exists {
		showSecrets := false
		if secrets != nil {
			showSecrets = *secrets
		}
		if showSecrets {
			info.Value = sK
		} else {
			info.Value = HideSecretKey(sK)
		}
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
