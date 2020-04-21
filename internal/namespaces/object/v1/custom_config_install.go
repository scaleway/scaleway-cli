package object

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func installCommand() *core.Command {
	type installRequest struct {
		Region scw.Region
		Type   string
	}
	return &core.Command{
		Namespace: "object",
		Resource:  "config",
		Verb:      "install",
		Short:     "Install a S3 related configuration file to its default location",
		Long:      "Install a S3 related configuration file.",
		ArgsType:  reflect.TypeOf(installRequest{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "type",
				Short:      "Type of tool supported",
				Required:   true,
				EnumValues: supportedTools,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short: "Install a s3cmd config file for Paris region",
				Raw:   "scw object config install type=s3cmd region=fr-par",
			},
			{
				Short: "Install a rclone config file for default region",
				Raw:   "scw object config install type=rclone",
			},

			{
				Short: "Install a mc (minio) config file for default region",
				Raw:   "scw object config install type=mc",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Generate a S3 tool configuration file",
				Command: "scw object config get",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			requestedType := argsI.(*installRequest)
			region := requestedType.Region.String()

			config, err := createS3Config(ctx, region)
			if err != nil {
				return "", err
			}
			var configPath string
			switch requestedType.Type {
			case "s3cmd":
				configPath, err = installS3cmd(config)
				if err != nil {
					return nil, err
				}

			case "rclone":
				configPath, err = installRclone(config)
				if err != nil {
					return nil, err
				}

			case "mc":
				configPath, err = installMc(config)
				if err != nil {
					return nil, err
				}

			default:
				return nil, &core.CliError{
					Message: "Unknown tool type",
					Details: fmt.Sprintf("%s is an unknown tool", requestedType.Type),
					Hint:    fmt.Sprintf("Try using on the following types: %s", supportedTools),
				}
			}
			return fmt.Sprintf("Configuration file successfully installed at %s", configPath), nil
		},
	}
}

func ensureFile(configPath string, newConfig string) error {
	// Ask whether to remove previous configuration file if it exists
	if _, err := os.Stat(configPath); err == nil {
		_, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Prompt:       "Do you want to overwrite the existing configuration file (" + configPath + ")?",
			DefaultValue: false,
		})
		if err != nil {
			return err
		}
		return ioutil.WriteFile(configPath, []byte(newConfig), 0644)
	}
	return ioutil.WriteFile(configPath, []byte(newConfig), 0644)
}

func installS3cmd(config s3config) (string, error) {
	newConfig, err := config.exportS3cmdConfig()
	if err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	s3cmdConfigPath := path.Join(homeDir, ".s3cfg")
	err = ensureFile(s3cmdConfigPath, newConfig)
	if err != nil {
		return "", err
	}
	return s3cmdConfigPath, nil
}

func installRclone(config s3config) (string, error) {
	newConfig, err := config.exportRcloneConfig()
	if err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// `rclone config file` returns the path of the configuration file
	rcloneConfigPath := path.Join(homeDir, ".config", "rclone", "rclone.conf")
	err = ensureFile(rcloneConfigPath, newConfig)
	if err != nil {
		return "", err
	}
	return rcloneConfigPath, nil
}

func installMc(config s3config) (string, error) {
	newConfig, err := config.exportMcConfig()
	if err != nil {
		return "", err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	mcConfigPath := path.Join(homeDir, ".mc", "config.json")
	err = ensureFile(mcConfigPath, newConfig)
	if err != nil {
		return "", err
	}
	return mcConfigPath, nil
}
