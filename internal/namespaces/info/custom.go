package info

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		infosRoot(),
	)
}

type infoResult struct {
	BuildInfo *core.BuildInfo
	Settings  []*setting
}

func (i infoResult) MarshalHuman() (string, error) {
	type tmp infoResult
	return human.Marshal(tmp(i), &human.MarshalOpt{
		Sections: []*human.MarshalSection{
			{
				FieldName: "build-info",
				Title:     "BuildInfo",
			},
			{
				FieldName: "settings",
				Title:     "Settings",
			},
		},
	})
}

type setting struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Origin string `json:"origin"`
}

func infosRoot() *core.Command {
	type infoArgs struct {
		ShowSecret bool
	}

	return &core.Command{
		Short:                `Get settings about current settings`,
		Namespace:            "info",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(infoArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "show-secret",
				Short:    `Reveal secret`,
				Required: false,
				Default:  core.DefaultValueSetter("false"),
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			req := argsI.(*infoArgs)
			config, _ := scw.LoadConfigFromPath(core.ExtractConfigPath(ctx))
			profileName := core.ExtractProfileName(ctx)
			return &infoResult{
				BuildInfo: core.ExtractBuildInfo(ctx),
				Settings: []*setting{
					configPath(ctx),
					defaultRegion(ctx, config, profileName),
					defaultZone(ctx, config, profileName),
					defaultOrganizationId(ctx, config, profileName),
					accessKey(ctx, config, profileName),
					secretKey(ctx, config, profileName, req.ShowSecret),
					profile(ctx),
				},
			}, nil
		},
	}
}

func configPath(ctx context.Context) *setting {
	setting := &setting{
		Key:   "config_path",
		Value: core.ExtractConfigPath(ctx),
	}
	switch {
	case core.ExtractConfigPathFlag(ctx) != "":
		setting.Origin = "flag --config/-c"
	case os.Getenv(scw.ScwConfigPathEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwConfigPathEnv)
	default:
		setting.Origin = "default"
	}
	return setting
}

func defaultRegion(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_region"}
	client := core.ExtractClient(ctx)
	defaultRegion, exists := client.GetDefaultRegion()
	if exists {
		setting.Value = defaultRegion.String()
	}
	switch {
	// Environment variable check
	case os.Getenv(scw.ScwDefaultRegionEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultRegionEnv)
	// There is no config file
	case config == nil:
		setting.Origin = "default"
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultRegion != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.Profile.DefaultRegion != nil:
		setting.Origin = "default profile"
	default:
		setting.Origin = "default"
	}
	return setting
}

func defaultZone(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_zone"}
	client := core.ExtractClient(ctx)
	defaultZone, exists := client.GetDefaultZone()
	if exists {
		setting.Value = defaultZone.String()
	}
	switch {
	// Environment variable check
	case os.Getenv(scw.ScwDefaultZoneEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultZoneEnv)
	// There is no config file
	case config == nil:
		setting.Origin = "default"
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultZone != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.Profile.DefaultZone != nil:
		setting.Origin = "default profile"
	default:
		setting.Origin = "default"
	}
	return setting
}

func defaultOrganizationId(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_organization_id"}
	client := core.ExtractClient(ctx)
	defaultOrganizationId, exists := client.GetDefaultOrganizationID()
	if exists {
		setting.Value = defaultOrganizationId
	}
	switch {
	// Environment variable check
	case os.Getenv(scw.ScwDefaultOrganizationIDEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultOrganizationIDEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultOrganizationID != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.Profile.DefaultOrganizationID != nil:
		setting.Origin = "default profile"
	default:
		setting.Origin = "unknown"
	}
	return setting
}

func accessKey(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "access_key"}
	client := core.ExtractClient(ctx)
	aK, exists := client.GetAccessKey()
	if exists {
		setting.Value = aK
	}
	switch {
	// Environment variable check
	case os.Getenv(scw.ScwAccessKeyEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwAccessKeyEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].AccessKey != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.Profile.AccessKey != nil:
		setting.Origin = "default profile"
	default:
		setting.Origin = "unknown"
	}
	return setting
}

func hideSecretKey(k string) string {
	switch {
	case len(k) == 0:
		return ""
	case len(k) > 8:
		return k[0:8] + "-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	default:
		return "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx"
	}
}

func secretKey(ctx context.Context, config *scw.Config, profileName string, showSecret bool) *setting {
	setting := &setting{Key: "secret_key"}
	client := core.ExtractClient(ctx)
	sK, exists := client.GetSecretKey()
	if exists {
		if showSecret {
			setting.Value = sK
		} else {
			setting.Value = hideSecretKey(sK)
		}
	}
	switch {
	// Environment variable check
	case os.Getenv(scw.ScwSecretKeyEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwSecretKeyEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].SecretKey != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.Profile.SecretKey != nil:
		setting.Origin = "default profile"
	default:
		setting.Origin = "unknown"
	}

	return setting
}

func profile(ctx context.Context) *setting {
	setting := &setting{
		Key:   "profile",
		Value: core.ExtractProfileName(ctx),
	}
	switch {
	case core.ExtractProfileFlag(ctx) != "":
		setting.Origin = "flag --profile/-p"
	case os.Getenv(scw.ScwActiveProfileEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwActiveProfileEnv)
	default:
		setting.Origin = ""
	}
	return setting
}
