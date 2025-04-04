package info

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	defaultOrigin        = "default"
	defaultProfileOrigin = "default profile"
	unknownOrigin        = "unknown"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		infosRoot(),
	)
}

type infoResult struct {
	BuildInfo *core.BuildInfo `json:"build_info"`
	Settings  []*setting      `json:"settings"`
}

func (i infoResult) MarshalHuman() (string, error) {
	type tmp infoResult

	return human.Marshal(tmp(i), &human.MarshalOpt{
		Sections: []*human.MarshalSection{
			{
				FieldName: "BuildInfo",
				Title:     "Build Info",
			},
			{
				FieldName: "Settings",
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
		Groups:               []string{"config"},
		Short:                `Get info about current settings`,
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
					profile(ctx, config),
					defaultRegion(ctx, config, profileName),
					defaultZone(ctx, config, profileName),
					defaultOrganizationID(ctx, config, profileName),
					defaultProjectID(ctx, config, profileName),
					accessKey(ctx, config, profileName),
					secretKey(ctx, config, profileName, req.ShowSecret),
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
		setting.Value = core.ExtractConfigPathFlag(ctx)
	case core.ExtractEnv(ctx, scw.ScwConfigPathEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwConfigPathEnv)
		setting.Value = core.ExtractEnv(ctx, scw.ScwConfigPathEnv)
	default:
		setting.Origin = defaultOrigin
	}

	return setting
}

func profile(ctx context.Context, config *scw.Config) *setting {
	setting := &setting{
		Key:   "profile",
		Value: core.ExtractProfileName(ctx),
	}

	switch {
	case core.ExtractProfileFlag(ctx) != "":
		setting.Origin = "flag --profile/-p"
		setting.Value = core.ExtractProfileFlag(ctx)
	case core.ExtractEnv(ctx, scw.ScwActiveProfileEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwActiveProfileEnv)
		setting.Value = core.ExtractEnv(ctx, scw.ScwActiveProfileEnv)
	case config != nil && config.ActiveProfile != nil:
		setting.Origin = "active_profile in config file"
		setting.Value = *config.ActiveProfile
	default:
		setting.Origin = ""
	}

	return setting
}

func defaultRegion(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_region"}
	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwDefaultRegionEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultRegionEnv)
		setting.Value = core.ExtractEnv(ctx, scw.ScwDefaultRegionEnv)
	// There is no config file
	case config == nil:
		setting.Origin = defaultOrigin
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultRegion != nil:
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
		setting.Value = *config.Profiles[profileName].DefaultRegion
	// Default config
	case config.DefaultRegion != nil:
		setting.Value = *config.DefaultRegion
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = defaultOrigin
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
	case core.ExtractEnv(ctx, scw.ScwDefaultZoneEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultZoneEnv)
		setting.Value = core.ExtractEnv(ctx, scw.ScwDefaultZoneEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultZone != nil:
		setting.Value = *config.Profiles[profileName].DefaultZone
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.DefaultZone != nil:
		setting.Value = *config.DefaultZone
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = defaultOrigin
	}

	return setting
}

func defaultOrganizationID(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_organization_id"}
	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwDefaultOrganizationIDEnv) != "":
		setting.Value = core.ExtractEnv(ctx, scw.ScwDefaultOrganizationIDEnv)
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultOrganizationIDEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultOrganizationID != nil:
		setting.Value = *config.Profiles[profileName].DefaultOrganizationID
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.DefaultOrganizationID != nil:
		setting.Value = *config.DefaultOrganizationID
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = unknownOrigin
	}

	return setting
}

func defaultProjectID(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "default_project_id"}
	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwDefaultProjectIDEnv) != "":
		setting.Value = core.ExtractEnv(ctx, scw.ScwDefaultProjectIDEnv)
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwDefaultProjectIDEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].DefaultProjectID != nil:
		setting.Value = *config.Profiles[profileName].DefaultProjectID
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.DefaultProjectID != nil:
		setting.Value = *config.DefaultProjectID
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = unknownOrigin
	}

	return setting
}

func accessKey(ctx context.Context, config *scw.Config, profileName string) *setting {
	setting := &setting{Key: "access_key"}
	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwAccessKeyEnv) != "":
		setting.Value = core.ExtractEnv(ctx, scw.ScwAccessKeyEnv)
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwAccessKeyEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].AccessKey != nil:
		setting.Value = *config.Profiles[profileName].AccessKey
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.AccessKey != nil:
		setting.Value = *config.AccessKey
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = unknownOrigin
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

func secretKey(
	ctx context.Context,
	config *scw.Config,
	profileName string,
	showSecret bool,
) *setting {
	setting := &setting{Key: "secret_key"}
	switch {
	// Environment variable check
	case core.ExtractEnv(ctx, scw.ScwSecretKeyEnv) != "":
		setting.Origin = fmt.Sprintf("env (%s)", scw.ScwSecretKeyEnv)
		setting.Value = core.ExtractEnv(ctx, scw.ScwSecretKeyEnv)
	// There is no config file
	case config == nil:
		setting.Origin = ""
	// Config file with profile name
	case config.Profiles[profileName] != nil && config.Profiles[profileName].SecretKey != nil:
		setting.Value = *config.Profiles[profileName].SecretKey
		setting.Origin = fmt.Sprintf("profile (%s)", profileName)
	// Default config
	case config.SecretKey != nil:
		setting.Value = *config.SecretKey
		setting.Origin = defaultProfileOrigin
	default:
		setting.Origin = unknownOrigin
	}
	if !showSecret {
		setting.Value = hideSecretKey(setting.Value)
	}

	return setting
}
