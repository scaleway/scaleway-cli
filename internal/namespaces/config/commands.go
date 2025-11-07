package config

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/config"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		configRoot(),
		configGetCommand(),
		configSetCommand(),
		configUnsetCommand(),
		configDumpCommand(),
		configProfileCommand(),
		configDeleteProfileCommand(),
		configActivateProfileCommand(),
		configResetCommand(),
		configDestroyCommand(),
		configInfoCommand(),
		configImportCommand(),
		configValidateCommand(),
		configEditCommand(),
	)
}

func configRoot() *core.Command {
	configPath := scw.GetConfigPath()
	envVarTable := bytes.Buffer{}
	w := tabwriter.NewWriter(&envVarTable, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)
	for _, envVar := range [][5]string{
		{"|", "Environment Variable", "|", "Description", "|"},
		{"|", "--", "|", "--", "|"},
		{"|", scw.ScwAccessKeyEnv, "|", "The access key of a token (create a token at https://console.scaleway.com/iam/api-keys)", "|"},
		{"|", scw.ScwSecretKeyEnv, "|", "The secret key of a token (create a token at https://console.scaleway.com/iam/api-keys)", "|"},
		{"|", scw.ScwDefaultOrganizationIDEnv, "|", "The default organization ID (get your organization ID at https://console.scaleway.com/iam/api-keys)", "|"},
		{"|", scw.ScwDefaultProjectIDEnv, "|", "The default project ID (get your project ID at https://console.scaleway.com/iam/api-keys)", "|"},
		{"|", scw.ScwDefaultRegionEnv, "|", "The default region", "|"},
		{"|", scw.ScwDefaultZoneEnv, "|", "The default availability zone", "|"},
		{"|", scw.ScwAPIURLEnv, "|", "URL of the API", "|"},
		{"|", scw.ScwInsecureEnv, "|", "Set this to true to enable the insecure mode", "|"},
		{"|", scw.ScwActiveProfileEnv, "|", "Set the config profile to use", "|"},
	} {
		fmt.Fprintf(
			w,
			"  %s%s%s%s%s\n",
			envVar[0],
			terminal.Style(envVar[1], color.Bold, color.FgBlue),
			envVar[2],
			envVar[3],
			envVar[4],
		)
	}
	w.Flush()

	return &core.Command{
		Groups: []string{"config"},
		Short:  `Config file management`,
		Long: interactive.RemoveIndent(`
			Config management engine is common across all Scaleway developer tools (CLI, terraform, SDK, ... ). It allows to handle Scaleway config through two ways: environment variables and/or config file.
			Default path for configuration file is based on the following priority order:

			- $SCW_CONFIG_PATH
			- $XDG_CONFIG_HOME/scw/config.yaml
			- $HOME/.config/scw/config.yaml
			- $USERPROFILE/.config/scw/config.yaml

			In this CLI, ` + terminal.Style(`environment variables have priority over the configuration file`, color.Bold) + `.

			The following environment variables are supported:

			` + envVarTable.String() + `
			Read more about the config management engine at https://github.com/scaleway/scaleway-sdk-go/tree/master/scw#scaleway-config
		`),
		Namespace: "config",
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Init your Scaleway config",
				Command: "scw init",
			},
			{
				Short:   "Set a config attribute",
				Command: "scw config set",
			},
			{
				Short:   "Set a config attribute",
				Command: "scw config get",
			},
			{
				Short:   "Dump the config",
				Command: "scw config dump",
			},
			{
				Short:   "Display the actual config file",
				Command: "cat " + configPath,
			},
		},
	}
}

// configGetCommand gets one or many values for the scaleway config
func configGetCommand() *core.Command {
	type configGetArgs struct {
		Key string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Get a value from the config file`,
		Namespace:            "config",
		Resource:             "get",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configGetArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key",
				Short:      "the key to get from the config",
				Required:   true,
				EnumValues: getProfileKeys(),
				Positional: true,
			},
		},
		Examples: []*core.Example{
			{
				Short: "Get the default organization ID",
				Raw:   "scw config get default_organization_id",
			},
			{
				Short: "Get the default region of the profile 'prod'",
				Raw:   "scw -p prod config get default_region",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Config management help",
				Command: "scw config",
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			config, err := scw.LoadConfigFromPath(core.ExtractConfigPath(ctx))
			if err != nil {
				return nil, err
			}
			key := argsI.(*configGetArgs).Key

			profileName := core.ExtractProfileName(ctx)
			profile, err := getProfile(config, profileName)
			if err != nil {
				return nil, err
			}

			return getProfileValue(profile, key)
		},
	}
}

// configSetCommand sets a value for the scaleway config
func configSetCommand() *core.Command {
	allRegions := []string(nil)
	for _, region := range scw.AllRegions {
		allRegions = append(allRegions, region.String())
	}
	allZones := []string(nil)
	for _, zone := range scw.AllZones {
		allZones = append(allZones, zone.String())
	}

	return &core.Command{
		Groups: []string{"config"},
		Short:  `Set a line from the config file`,
		Long: `This commands overwrites the configuration file parameters with user input.
The only allowed attributes are access_key, secret_key, default_organization_id, default_region, default_zone, api_url, insecure`,
		Namespace:            "config",
		Resource:             "set",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(scw.Profile{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:  "access-key",
				Short: "A Scaleway access key",
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if !reflect.ValueOf(value).IsNil() &&
						!validation.IsAccessKey(*value.(*string)) {
						return core.InvalidAccessKeyError(*value.(*string))
					}

					return nil
				},
			},
			{
				Name:  "secret-key",
				Short: "A Scaleway secret key",
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if !reflect.ValueOf(value).IsNil() &&
						!validation.IsSecretKey(*value.(*string)) {
						return core.InvalidSecretKeyError(*value.(*string))
					}

					return nil
				},
			},
			{
				Name:  "api-url",
				Short: "Scaleway API URL",
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if !reflect.ValueOf(value).IsNil() && !validation.IsURL(*value.(*string)) {
						return fmt.Errorf("%s is not a valid URL", *value.(*string))
					}

					return nil
				},
			},
			{
				Name:  "insecure",
				Short: "Set to true to allow insecure HTTPS connections",
			},
			{
				Name:  "default-organization-id",
				Short: "A default Scaleway organization id",
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if !reflect.ValueOf(value).IsNil() &&
						!validation.IsOrganizationID(*value.(*string)) {
						return core.InvalidOrganizationIDError(*value.(*string))
					}

					return nil
				},
			},
			{
				Name:  "default-project-id",
				Short: "A default Scaleway project id",
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if !reflect.ValueOf(value).IsNil() &&
						!validation.IsProjectID(*value.(*string)) {
						return core.InvalidProjectIDError(*value.(*string))
					}

					return nil
				},
			},
			{
				Name:       "default-region",
				Short:      "A default Scaleway region",
				EnumValues: allRegions,
			},
			{
				Name:       "default-zone",
				Short:      "A default Scaleway zone",
				EnumValues: allZones,
			},
			{
				Name:  "send-telemetry",
				Short: "Set to false to disable telemetry",
			},
		},
		Examples: []*core.Example{
			{
				Short: "Update the default organization ID",
				Raw:   "scw config set default_organization_id=12903058-d0e8-4366-89c3-6e666abe1f6f",
			},
			{
				Short: "Update the default region of the profile 'prod'",
				Raw:   "scw -p prod config set default_region=nl-ams",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Config management help",
				Command: "scw config",
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			// Validate arguments
			args := argsI.(*scw.Profile)

			// Execute
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}

			// send_telemetry is the only key that is not in a profile but in the config object directly
			profileName := core.ExtractProfileName(ctx)
			profile := &config.Profile
			if profileName != scw.DefaultProfileName {
				var exist bool
				profile, exist = config.Profiles[profileName]
				if !exist {
					if config.Profiles == nil {
						config.Profiles = map[string]*scw.Profile{}
					}
					config.Profiles[profileName] = &scw.Profile{}
					profile = config.Profiles[profileName]
				}
			}

			argValue := reflect.ValueOf(args).Elem()
			profileValue := reflect.ValueOf(profile).Elem()
			for i := range argValue.NumField() {
				field := argValue.Field(i)
				if !field.IsNil() {
					profileValue.Field(i).Set(field)
				}
			}

			// Save
			err = config.SaveTo(configPath)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully update config",
			}, nil
		},
	}
}

// configDumpCommand unsets a value for the scaleway config
func configUnsetCommand() *core.Command {
	type configUnsetArgs struct {
		Key string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Unset a line from the config file`,
		Namespace:            "config",
		Resource:             "unset",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configUnsetArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "key",
				Short:      "the config config key name to unset",
				Required:   true,
				EnumValues: getProfileKeys(),
				Positional: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}
			key := argsI.(*configUnsetArgs).Key

			profileName := core.ExtractProfileName(ctx)
			profile, err := getProfile(config, profileName)
			if err != nil {
				return nil, err
			}
			err = unsetProfileValue(profile, key)
			if err != nil {
				return nil, err
			}

			err = config.SaveTo(configPath)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully unset " + key,
			}, nil
		},
	}
}

// configDumpCommand dumps the scaleway config
func configDumpCommand() *core.Command {
	type configDumpArgs struct{}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Dump the config file`,
		Namespace:            "config",
		Resource:             "dump",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configDumpArgs{}),
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Config management help",
				Command: "scw config",
			},
		},
		Run: func(ctx context.Context, _ any) (i any, e error) {
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}

			return config, nil
		},
	}
}

func configProfileCommand() *core.Command {
	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Allows the activation and deletion of a profile from the config file`,
		Namespace:            "config",
		Resource:             "profile",
		AllowAnonymousClient: true,
	}
}

// configDeleteProfileCommand deletes a profile from the config
func configDeleteProfileCommand() *core.Command {
	type configDeleteProfileArgs struct {
		Name string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Delete a profile from the config file`,
		Namespace:            "config",
		Resource:             "profile",
		Verb:                 "delete",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configDeleteProfileArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   true,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			profileName := argsI.(*configDeleteProfileArgs).Name
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}
			if _, exists := config.Profiles[profileName]; exists {
				delete(config.Profiles, profileName)
			} else {
				return nil, unknownProfileError(profileName)
			}
			err = config.SaveTo(configPath)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully delete profile " + profileName,
			}, nil
		},
	}
}

// configActivateProfileCommand mark a profile as active
func configActivateProfileCommand() *core.Command {
	type configActiveProfileArgs struct {
		ProfileName string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Mark a profile as active in the config file`,
		Namespace:            "config",
		Resource:             "profile",
		Verb:                 "activate",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configActiveProfileArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:             "profile-name",
				Required:         true,
				Positional:       true,
				AutoCompleteFunc: core.AutocompleteProfileName(),
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			profileName := argsI.(*configActiveProfileArgs).ProfileName
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}

			if profileName == scw.DefaultProfileName {
				config.ActiveProfile = nil
			} else {
				if _, exists := config.Profiles[profileName]; !exists {
					return nil, unknownProfileError(profileName)
				}
				config.ActiveProfile = &profileName
			}

			err = config.SaveTo(configPath)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully activate profile " + profileName,
			}, nil
		},
	}
}

// configResetCommand resets the config
func configResetCommand() *core.Command {
	type configResetArgs struct{}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Reset the config`,
		Namespace:            "config",
		Resource:             "reset",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configResetArgs{}),
		Run: func(_ context.Context, _ any) (i any, e error) {
			_, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			config := &scw.Config{}
			err = config.Save()
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully reset config",
			}, nil
		},
	}
}

// configDestroyCommand destroys the config
func configDestroyCommand() *core.Command {
	type configDestroyArgs struct{}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Destroy the config file`,
		Namespace:            "config",
		Resource:             "destroy",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configDestroyArgs{}),
		Run: func(ctx context.Context, _ any) (i any, e error) {
			configPath := core.ExtractConfigPath(ctx)
			err := os.Remove(configPath)
			if err != nil {
				return nil, err
			}

			return &core.SuccessResult{
				Message: "successfully destroy config",
			}, nil
		},
	}
}

// configInfoCommand values from the scaleway config for the current profile
func configInfoCommand() *core.Command {
	type configInfoArgs struct{}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                `Get config values from the config file for the current profile`,
		Namespace:            "config",
		Resource:             "info",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configInfoArgs{}),
		ArgSpecs:             core.ArgSpecs{},
		Examples: []*core.Example{
			{
				Short: "Get the default config values",
				Raw:   "scw config info",
			},
			{
				Short: "Get the config values of the profile 'prod'",
				Raw:   "scw -p prod config info",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Config management help",
				Command: "scw config",
			},
		},
		Run: func(ctx context.Context, _ any) (i any, e error) {
			config, err := scw.LoadConfigFromPath(core.ExtractConfigPath(ctx))
			if err != nil {
				return nil, err
			}

			profileEnv := scw.LoadEnvProfile()

			// Search for env variable that will override profile
			// Will be used to display them
			overriddenVariables := []string(nil)
			for _, key := range getProfileKeys() {
				value, err := getProfileField(profileEnv, key)
				if err == nil && !value.IsZero() {
					overriddenVariables = append(overriddenVariables, key)
				}
			}

			profileName := core.ExtractProfileName(ctx)
			// use config.GetProfile instead of getProfile as we want the profile merged with the default
			profile, err := config.GetProfile(profileName)
			if err != nil {
				return nil, err
			}

			profile = scw.MergeProfiles(profile, profileEnv)

			values := map[string]any{}
			for _, key := range getProfileKeys() {
				value, err := getProfileValue(profile, key)
				if err == nil && value != nil {
					values[key] = value
				}
			}

			if len(overriddenVariables) > 0 {
				msg := "Some variables are overridden by the environment: " + strings.Join(
					overriddenVariables,
					", ",
				)
				fmt.Println(terminal.Style(msg, color.FgRed))
			}

			return struct {
				ConfigPath  string
				ProfileName string
				Profile     map[string]any
			}{
				ConfigPath:  core.ExtractConfigPath(ctx),
				ProfileName: core.ExtractProfileName(ctx),
				Profile:     values,
			}, nil
		},
	}
}

// configImportCommand imports an external config
func configImportCommand() *core.Command {
	type configImportArgs struct {
		File string
	}

	return &core.Command{
		Groups:               []string{"config"},
		Short:                "Import configurations from another file",
		Namespace:            "config",
		Resource:             "import",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configImportArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "file",
				Short:      "Path to the configuration file to import",
				Required:   true,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, argsI any) (i any, e error) {
			args := argsI.(*configImportArgs)
			configPath := core.ExtractConfigPath(ctx)

			currentConfig, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}
			currentProfileName := core.ExtractProfileName(ctx)
			currentProfile, err := currentConfig.GetProfile(currentProfileName)
			if err != nil {
				return nil, err
			}

			// Read the content of the file to import
			importedConfig, err := scw.LoadConfigFromPath(args.File)
			if err != nil {
				return nil, err
			}
			importedProfile := importedConfig.Profile

			// Merge the imported configurations into the existing configuration
			currentConfig.Profile = *scw.MergeProfiles(currentProfile, &importedProfile)

			for profileName, profile := range importedConfig.Profiles {
				existingProfile, exists := currentConfig.Profiles[profileName]
				if exists {
					currentConfig.Profiles[profileName] = scw.MergeProfiles(
						existingProfile,
						profile,
					)
				} else {
					currentConfig.Profiles[profileName] = profile
				}
			}

			err = currentConfig.SaveTo(configPath)
			if err != nil {
				return nil, fmt.Errorf("failed to save updated configuration: %v", err)
			}

			return &core.SuccessResult{
				Message: "successfully import config",
			}, nil
		},
	}
}

// configValidateCommand validates the config
func configValidateCommand() *core.Command {
	type configValidateArgs struct{}

	return &core.Command{
		Short: `Validate the config`,
		Long: `This command validates the configuration of your Scaleway CLI tool.

It performs the following checks:

	- YAML syntax correctness: It checks whether your config file is a valid YAML file.
	- Field validity: It checks whether the fields present in the config file are valid and expected fields. This includes fields like AccessKey, SecretKey, DefaultOrganizationID, DefaultProjectID, DefaultRegion, DefaultZone, and APIURL.
	- Field values: For each of the fields mentioned above, it checks whether the value assigned to it is valid. For example, it checks if the AccessKey and SecretKey are non-empty and meet the format expectations.

The command goes through each profile present in the config file and validates it.`,
		Namespace:            "config",
		Resource:             "validate",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(configValidateArgs{}),
		Run: func(ctx context.Context, _ any) (i any, e error) {
			configPath := core.ExtractConfigPath(ctx)
			config, err := scw.LoadConfigFromPath(configPath)
			if err != nil {
				return nil, err
			}

			// validate default profile
			err = validateProfile(&config.Profile)
			if err != nil {
				return nil, err
			}
			// validate the remaining profiles
			for _, profile := range config.Profiles {
				err = validateProfile(profile)
				if err != nil {
					return nil, err
				}
			}

			return &core.SuccessResult{
				Message: "successfully validate config",
			}, nil
		},
	}
}

func configEditCommand() *core.Command {
	type configEditArgs struct{}

	return &core.Command{
		Namespace:            "config",
		Resource:             "edit",
		Short:                "Edit the configuration file",
		Long:                 "Edit the configuration file with the default editor",
		ArgsType:             reflect.TypeOf(configEditArgs{}),
		AllowAnonymousClient: true,
		Run: func(ctx context.Context, _ any) (i any, e error) {
			configPath := core.ExtractConfigPath(ctx)

			defaultEditor := config.GetDefaultEditor()
			args := []string{configPath}

			cmd := exec.Command(defaultEditor, args...)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout

			err := cmd.Run()
			if err != nil {
				return nil, fmt.Errorf("failed to edit file %q: %w", configPath, err)
			}

			return &core.SuccessResult{
				Message: "successfully wrote config",
			}, nil
		},
	}
}

// Helper functions
func getProfileValue(profile *scw.Profile, fieldName string) (any, error) {
	field, err := getProfileField(profile, fieldName)
	if err != nil {
		return nil, err
	}

	return field.Interface(), nil
}

func unsetProfileValue(profile *scw.Profile, key string) error {
	field, err := getProfileField(profile, key)
	if err != nil {
		return err
	}
	field.Set(reflect.Zero(field.Type()))

	return nil
}

func getProfileField(profile *scw.Profile, key string) (reflect.Value, error) {
	field := reflect.ValueOf(profile).Elem().FieldByName(strcase.ToPublicGoName(key))
	if !field.IsValid() {
		return reflect.ValueOf(nil), invalidProfileKeyError(key)
	}

	return field, nil
}

func getProfileKeys() []string {
	t := reflect.TypeOf(scw.Profile{})
	keys := []string{}
	for i := range t.NumField() {
		field := t.Field(i)
		switch field.Name {
		case "APIURL":
			keys = append(keys, "api-url")
		default:
			keys = append(keys, strcase.ToBashArg(t.Field(i).Name))
		}
	}

	return keys
}

// getProfile return a config profile by its name.
// Warning: This return the profile pointer directly so it can be modified by commands.
// For this reason we cannot rely on config.GetProfileByName method as it create a copy.
func getProfile(config *scw.Config, profileName string) (*scw.Profile, error) {
	if profileName == scw.DefaultProfileName {
		return &config.Profile, nil
	}
	profile, exist := config.Profiles[profileName]
	if !exist {
		return nil, unknownProfileError(profileName)
	}

	return profile, nil
}

func validateProfile(profile *scw.Profile) error {
	if err := validateAccessKey(profile); err != nil {
		return err
	}
	if err := validateSecretKey(profile); err != nil {
		return err
	}
	if err := validateDefaultOrganizationID(profile); err != nil {
		return err
	}
	if err := validateDefaultProjectID(profile); err != nil {
		return err
	}
	if err := validateDefaultRegion(profile); err != nil {
		return err
	}
	if err := validateDefaultZone(profile); err != nil {
		return err
	}

	return validateAPIURL(profile)
}

func validateAccessKey(profile *scw.Profile) error {
	if profile.AccessKey != nil {
		if *profile.AccessKey == "" {
			return &core.CliError{
				Err: errors.New("access key cannot be empty"),
			}
		}

		if !validation.IsAccessKey(*profile.AccessKey) {
			return core.InvalidAccessKeyError(*profile.AccessKey)
		}
	}

	return nil
}

func validateSecretKey(profile *scw.Profile) error {
	if profile.SecretKey != nil {
		if *profile.SecretKey == "" {
			return &core.CliError{
				Err: errors.New("secret key cannot be empty"),
			}
		}

		if !validation.IsSecretKey(*profile.SecretKey) {
			return core.InvalidSecretKeyError(*profile.SecretKey)
		}
	}

	return nil
}

func validateDefaultOrganizationID(profile *scw.Profile) error {
	if profile.DefaultOrganizationID != nil {
		if *profile.DefaultOrganizationID == "" {
			return &core.CliError{
				Err: errors.New("default organization ID cannot be empty"),
			}
		}

		if !validation.IsOrganizationID(*profile.DefaultOrganizationID) {
			return core.InvalidOrganizationIDError(*profile.DefaultOrganizationID)
		}
	}

	return nil
}

func validateDefaultProjectID(profile *scw.Profile) error {
	if profile.DefaultProjectID != nil {
		if *profile.DefaultProjectID == "" {
			return &core.CliError{
				Err: errors.New("default project ID cannot be empty"),
			}
		}

		if !validation.IsProjectID(*profile.DefaultProjectID) {
			return core.InvalidProjectIDError(*profile.DefaultProjectID)
		}
	}

	return nil
}

func validateDefaultRegion(profile *scw.Profile) error {
	if profile.DefaultRegion != nil {
		if *profile.DefaultRegion == "" {
			return &core.CliError{
				Err: errors.New("default region cannot be empty"),
			}
		}

		if !validation.IsRegion(*profile.DefaultRegion) {
			return core.InvalidRegionError(*profile.DefaultRegion)
		}
	}

	return nil
}

func validateDefaultZone(profile *scw.Profile) error {
	if profile.DefaultZone != nil {
		if *profile.DefaultZone == "" {
			return &core.CliError{
				Err: errors.New("default zone cannot be empty"),
			}
		}

		if !validation.IsZone(*profile.DefaultZone) {
			return core.InvalidZoneError(*profile.DefaultZone)
		}
	}

	return nil
}

func validateAPIURL(profile *scw.Profile) error {
	if profile.APIURL != nil {
		if *profile.APIURL != "" && !validation.IsURL(*profile.APIURL) {
			return core.InvalidAPIURLError(*profile.APIURL)
		}
	}

	return nil
}
