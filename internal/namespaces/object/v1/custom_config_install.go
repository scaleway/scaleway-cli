//go:build darwin || linux || windows

package object

import (
	"context"
	"path/filepath"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
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
		Run: func(ctx context.Context, argsI any) (any, error) {
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

			// Extract the home directory and relative path
			homeDir := core.ExtractUserHomeDir(ctx)
			relPath, err := filepath.Rel(homeDir, configPath)
			if err != nil {
				return "", err
			}

			// Create options for WriteFile
			opts := &interactive.WriteFileOptions{
				Confirm: true,
			}

			// Construct the full path
			fullPath := filepath.Join(homeDir, relPath)

			// Write the configuration file using the utility function
			err = interactive.WriteFile(ctx, fullPath, []byte(newConfig), opts)
			if err != nil {
				return "", err
			}

			return &core.SuccessResult{
				Message: "Configuration file successfully installed at " + configPath,
			}, nil
		},
	}
}
