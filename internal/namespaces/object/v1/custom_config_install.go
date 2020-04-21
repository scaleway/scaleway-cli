package object

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
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
		Short:     "Install a S3 related configuration file to its default location",
		Long:      "Install a S3 related configuration file.",
		ArgsType:  reflect.TypeOf(installArgs{}),
		ArgSpecs: []*core.ArgSpec{
			{
				Name:       "type",
				Short:      "Type of S3 tool you want to generate a config for",
				Required:   true,
				EnumValues: supportedTools,
			},
			{
				Name:     "name",
				Short:    "Name of the s3 remote you want to generate",
				Required: false,
				Default: func() (value string, doc string) {
					return "scaleway", "default value"
				},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:   "Install a s3cmd config file for Paris region",
				Request: `{"type": "s3cmd", "region": "fr-par"}`,
			},
			{
				Short:   "Install a rclone config file for default region",
				Request: `{"type": "rclone"}`,
			},

			{
				Short:   "Install a mc (minio) config file for default region",
				Request: `{"type": "mc"}`,
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
			region := args.Region
			name := args.Name

			config, err := newS3Config(ctx, region, name)
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
				_, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
					Prompt:       "Do you want to overwrite the existing configuration file (" + configPath + ")?",
					DefaultValue: false,
				})
				if err != nil {
					return nil, err
				}
			}
			err = ioutil.WriteFile(configPath, []byte(newConfig), 0644)
			if err != nil {
				return "", err
			}
			return fmt.Sprintf("Configuration file successfully installed at %s", configPath), nil
		},
	}
}
