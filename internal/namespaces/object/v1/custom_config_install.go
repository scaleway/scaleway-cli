package object

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func configInstallCommand() *core.Command {
	type installArgs struct {
		Region scw.Region
		Type   s3tool
		Name   string
	}
	return &core.Command{
		Namespace: "object",
		Resource:  "config",
		Verb:      "install",
		Short:     "Install a S3 tool configuration file to its default location",
		Long:      "Install a S3 tool configuration file to its default location.",
		ArgsType:  reflect.TypeOf(installArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "type",
				Short:      "Type of S3 tool you want to generate a config for",
				Required:   true,
				EnumValues: supportedTools.ToStringArray(),
			},
			{
				Name:     "name",
				Short:    "Name of the s3 remote you want to generate",
				Required: false,
				Default:  core.DefaultValueSetter("scaleway"),
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Install a s3cmd config file for Paris region",
				ArgsJSON: `{"type": "s3cmd", "region": "fr-par"}`,
			},
			{
				Short:    "Install a rclone config file for default region",
				ArgsJSON: `{"type": "rclone"}`,
			},

			{
				Short:    "Install a mc (minio) config file for default region",
				ArgsJSON: `{"type": "mc"}`,
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Generate a S3 tool configuration file",
				Command: "scw object config get",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			args := argsI.(*installArgs)

			config, err := newS3Config(ctx, args.Region, args.Name)
			if err != nil {
				return "", err
			}
			newConfig, err := config.getConfigFile(args.Type)
			if err != nil {
				return "", err
			}
			configPath, err := config.getPath(args.Type)
			if err != nil {
				return "", err
			}

			// Ask whether to remove previous configuration file if it exists
			if _, err := os.Stat(configPath); err == nil {
				doIt, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
					Ctx:          ctx,
					Prompt:       "Do you want to overwrite the existing configuration file (" + configPath + ")?",
					DefaultValue: false,
				})
				if err != nil {
					return nil, err
				}
				if !doIt {
					return nil, fmt.Errorf("installation aborted by user")
				}
			}

			// Ensure the subfolders for the configuration files are all created
			err = os.MkdirAll(filepath.Dir(configPath), 0755)
			if err != nil {
				return "", err
			}

			// Write the configuration file
			err = ioutil.WriteFile(configPath, []byte(newConfig), 0600)
			if err != nil {
				return "", err
			}
			return &core.SuccessResult{
				Message: fmt.Sprintf("Configuration file successfully installed at %s", configPath),
			}, nil
		},
	}
}
