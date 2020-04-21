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
				EnumValues: []string{"rclone", "s3cmd", "mc"},
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

			switch requestedType.Type {
			case "s3cmd":
				i, err2 := installS3cmd(config)
				if err2 != nil {
					return i, err2
				}

			case "rclone":
				i, err2 := installRclone(config)
				if err2 != nil {
					return i, err2
				}

			case "mc":
				i, err2 := installMc(config)
				if err2 != nil {
					return i, err2
				}

			default:
				fmt.Println("Unknown tool.")
			}
			return nil, nil
		},
	}
}

func ensureFile(mcConfigPath string, newConfig string) (interface{}, error, bool) {
	// Ask whether to remove previous configuration file if it exists
	if _, err := os.Stat(mcConfigPath); err == nil {
		_, err := interactive.PromptBoolWithConfig(&interactive.PromptBoolConfig{
			Prompt:       "Do you want to overwrite the existing configuration file (" + mcConfigPath + ")?",
			DefaultValue: false,
		})
		if err != nil {
			return "", err, true
		}
		return nil, ioutil.WriteFile(mcConfigPath, []byte(newConfig), 0644), true
	}
	return nil, nil, false
}

func installS3cmd(config s3config) (interface{}, error) {
	newConfig, err := config.exportS3cmdConfig()
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	s3cmdConfigPath := path.Join(homeDir, ".s3cfg")
	i, err2, done := ensureFile(s3cmdConfigPath, newConfig)
	if done {
		return i, err2
	}
	return nil, nil
}

func installRclone(config s3config) (interface{}, error) {
	newConfig, err := config.exportRcloneConfig()
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	// `rclone config file` returns the path of the configuration file
	rcloneConfigPath := path.Join(homeDir, ".config", "rclone", "rclone.conf")
	i, err2, done := ensureFile(rcloneConfigPath, newConfig)
	if done {
		return i, err2
	}

	return nil, nil
}

func installMc(config s3config) (interface{}, error) {
	newConfig, err := config.exportMcConfig()
	if err != nil {
		return nil, err
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	mcConfigPath := path.Join(homeDir, ".mc", "config.json")
	i, err2, done := ensureFile(mcConfigPath, newConfig)
	if done {
		return i, err2
	}
	return nil, nil
}
