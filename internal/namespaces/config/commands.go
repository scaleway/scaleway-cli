package config

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// TODO: add proper tests

func GetCommands() *core.Commands {
	return core.NewCommands(
		configRoot(),
		configGetCommand(),
		configSetCommand(),
		configUnsetCommand(),
		configDumpCommand(),
		configDeleteCommand(),
		configDeleteProfileCommand(),
		configResetCommand(),
	)
}

func configRoot() *core.Command {
	return &core.Command{
		Short:     `Config file management`,
		Long:      `Manage your Scaleway CLI config file.`,
		Namespace: "config",
	}
}

// configGetCommand gets one or many values for the scaleway config
func configGetCommand() *core.Command {
	return &core.Command{
		Short:     `Get a line from the config file`,
		Namespace: "config",
		Resource:  "get",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {

			// profileKeyValue is a custom type used for displaying configGetCommand result
			type profileKeyValue struct {
				Profile string `json:"profile"`
				Key     string `json:"key"`
				Value   string `json:"value"`
			}

			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			rawArgs := *(argsI.(*args.RawArgs))
			if len(rawArgs) == 0 {
				return nil, notEnoughArgsForConfigGetError()
			}
			profileKeyValues := []*profileKeyValue(nil)
			for _, arg := range rawArgs {
				profileName, key, err := splitProfileKey(arg)
				if err != nil {
					return nil, err
				}
				profile, err := getProfile(config, profileName)
				if err != nil {
					return nil, err
				}
				value, err := getProfileValue(profile, key)
				if err != nil {
					return nil, err
				}
				profileKeyValues = append(profileKeyValues, &profileKeyValue{
					Profile: profileName,
					Key:     key,
					Value:   value,
				})
			}
			return profileKeyValues, nil
		},
	}
}

// configSetCommand sets a value for the scaleway config
func configSetCommand() *core.Command {
	return &core.Command{
		Short: `Set a line from the config file`,
		Long: `This commands overwrites the configuration file parameters with user input.
The only allowed attributes are access_key, secret_key, default_organization_id, default_region, default_zone, api_url, insecure`,
		Namespace: "config",
		Resource:  "set",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			// Validate arguments
			rawArgs := *(argsI.(*args.RawArgs))
			profileName, key, value, err := validateRawArgsForConfigSet(rawArgs)
			if err != nil {
				return nil, err
			}

			// Execute
			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			profile, err := getProfile(config, profileName) // There can not be an error if profileName is empty
			if err != nil {
				// We create the profile if it doesn't exist
				if config.Profiles == nil {
					config.Profiles = map[string]*scw.Profile{}
				}
				config.Profiles[profileName] = &scw.Profile{}
				profile = config.Profiles[profileName]
			}
			err = setProfileValue(profile, key, value)
			if err != nil {
				return nil, err
			}

			// Save
			err = config.Save()
			if err != nil {
				return nil, err
			}

			// Inform success
			return configSetSuccess(rawArgs[0], value), nil
		},
	}
}

// configDumpCommand unsets a value for the scaleway config
func configUnsetCommand() *core.Command {
	return &core.Command{
		Short:     `Unset a line from the config file`,
		Namespace: "config",
		Resource:  "unset",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			rawArgs := *(argsI.(*args.RawArgs))
			if len(rawArgs) == 0 {
				return nil, notEnoughArgsForConfigUnsetError()
			}
			if len(rawArgs) > 1 {
				return nil, tooManyArgsForConfigUnsetError()
			}
			profileAndKey := rawArgs[0]
			profileName, key, err := splitProfileKey(profileAndKey)
			if err != nil {
				return nil, err
			}
			profile, err := getProfile(config, profileName)
			if err != nil {
				return nil, err
			}
			logger.Debugf("conf before: %v", config)
			err = unsetProfileValue(profile, key)
			if err != nil {
				return nil, err
			}
			logger.Debugf("conf after: %v", config)
			err = config.Save()
			if err != nil {
				return nil, err
			}

			return configUnsetSuccess(profileAndKey), nil
		},
	}
}

// configDumpCommand dumps the scaleway config
func configDumpCommand() *core.Command {
	return &core.Command{
		Short:     `Dump the config file`,
		Namespace: "config",
		Resource:  "dump",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			return config, nil
		},
	}
}

func configDeleteCommand() *core.Command {
	return &core.Command{
		Short:     `Allows the deletion of a profile from the config file`,
		Namespace: "config",
		Resource:  "delete",
		NoClient:  true,
	}
}

type configDeleteProfileArgs struct {
	Name string
}

// configDeleteProfileCommand deletes a profile from the config
func configDeleteProfileCommand() *core.Command {
	return &core.Command{
		Short:     `Delete a profile from the config file`,
		Namespace: "config",
		Resource:  "delete",
		Verb:      "profile",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(configDeleteProfileArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:     "name",
				Required: true,
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			profileName := argsI.(*configDeleteProfileArgs).Name
			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			if _, exists := config.Profiles[profileName]; exists {
				delete(config.Profiles, profileName)
			} else {
				return nil, unknownProfileError(profileName)
			}
			err = config.Save()
			if err != nil {
				return nil, err
			}

			return configDeleteProfileSuccess(profileName), nil
		},
	}
}

// configResetCommand resets the config
func configResetCommand() *core.Command {
	return &core.Command{
		Short:     `Reset the config`,
		Namespace: "config",
		Resource:  "reset",
		NoClient:  true,
		ArgsType:  reflect.TypeOf(args.RawArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			config, err := scw.LoadConfig()
			if err != nil {
				return nil, err
			}
			config = &scw.Config{}
			err = config.Save()
			if err != nil {
				return nil, err
			}
			return configResetSuccess(), nil
		},
	}
}

//
// Helper functions
//

// getProfile returns scw.Config.Profiles[profileName] or scw.Config.Profile if profileName is empty.
func getProfile(config *scw.Config, profileName string) (profile *scw.Profile, err error) {
	if profileName != "" {
		profile, err := config.GetProfile(profileName)
		if err != nil {
			return nil, err
		}
		return profile, nil
	} else {
		return &config.Profile, nil
	}
}

// splitProfileKey splits a "profile.key" string into ("profile", "key")
func splitProfileKey(arg string) (profileName string, key string, err error) {
	strs := strings.Split(arg, ".")
	if len(strs) == 1 {
		return "", strs[0], nil
	}
	if len(strs) == 2 {
		return strs[0], strs[1], nil
	}
	return "", "", invalidProfileKeyPairError(arg)
}

func getProfileValue(profile *scw.Profile, fieldName string) (string, error) {
	field := reflect.ValueOf(profile).Elem().FieldByName(strcase.ToPublicGoName(fieldName))
	if field.IsValid() == false {
		return "", invalidProfileKeyIdentifierError(fieldName)
	}
	if field.IsNil() {
		return "", nilFieldError(fieldName)
	}
	return fmt.Sprint(field.Elem().Interface()), nil
}

func setProfileValue(profile *scw.Profile, fieldName string, value string) error {
	field := reflect.ValueOf(profile).Elem().FieldByName(strcase.ToPublicGoName(fieldName))
	switch kind := field.Type().Elem().Kind(); kind {
	case reflect.String:
		field.Set(reflect.ValueOf(&value))
	case reflect.Bool:
		field.Set(reflect.ValueOf(scw.BoolPtr(value == "true")))
	default:
		return invalidKindForKeyError(kind, fieldName)
	}
	return nil
}

func unsetProfileValue(profile *scw.Profile, key string) error {
	field, err := getProfileAttribute(profile, key)
	if err != nil {
		return err
	}
	field.Set(reflect.Zero(field.Type()))
	return nil
}

func getProfileAttribute(profile *scw.Profile, key string) (reflect.Value, error) {
	field := reflect.ValueOf(profile).Elem().FieldByName(strcase.ToPublicGoName(key))
	if field.IsValid() == false {
		return reflect.ValueOf(nil), invalidProfileAttributeError(key)
	}
	return field, nil
}
