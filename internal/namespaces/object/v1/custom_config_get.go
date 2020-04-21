package object

import (
	"bytes"
	"context"
	"encoding/json"
	"reflect"
	"text/template"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func getCommand() *core.Command {
	type getRequest struct {
		Region scw.Region
		Type   string
	}

	return &core.Command{
		Namespace: "object",
		Resource:  "config",
		Verb:      "get",
		Short:     "Generate a S3 related configuration file",
		Long:      "Generate a S3 related configuration file.",
		ArgsType:  reflect.TypeOf(getRequest{}),
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
				Short: "Generate a s3cmd s3config file for Paris region",
				Raw:   "scw object config install type=s3cmd region=fr-par",
			},
			{
				Short: "Generate a rclone s3config file for default region",
				Raw:   "scw object config install type=rclone",
			},
			{
				Short: "Generate a mc (minio) s3config file for default region",
				Raw:   "scw object config install type=mc",
			},
		},
		SeeAlsos: []*core.SeeAlso{
			{
				Short:   "Install a S3 tool configuration file",
				Command: "scw object config install",
			},
		},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			requestedType := argsI.(*getRequest)
			region := requestedType.Region.String()

			config, err := createS3Config(ctx, region)
			if err != nil {
				return "", err
			}

			switch requestedType.Type {
			case "s3cmd":
				res, err := config.exportS3cmdConfig()
				if err != nil {
					return nil, err
				}
				return res, nil
			case "rclone":
				res, err := config.exportRcloneConfig()
				if err != nil {
					return nil, err
				}
				return res, nil
			case "mc":
				res, err := config.exportMcConfig()
				if err != nil {
					return nil, err
				}
				return res, nil
			default:
				return nil, &core.CliError{
					Message: "",
					Details: "",
					Hint:    "",
				}
			}
		},
	}
}

func createS3Config(ctx context.Context, region string) (s3config, error) {
	client := core.ExtractClient(ctx)
	accessKey, accessExists := client.GetAccessKey()
	if !accessExists {
		return s3config{}, &core.CliError{
			Err:     nil,
			Message: "",
			Details: "",
			Hint:    "",
		}
	}
	secretKey, secretExists := client.GetSecretKey()
	if !secretExists {
		return s3config{}, &core.CliError{
			Err:     nil,
			Message: "",
			Details: "",
			Hint:    "",
		}
	}
	if region == "" {
		defaultRegion, _ := client.GetDefaultRegion()
		region = defaultRegion.String()
	}
	config := s3config{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Region:    region,
	}
	return config, nil
}

type s3config struct {
	AccessKey string
	SecretKey string
	Region    string
}

func renderTemplate(configFileTemplate string, c s3config) (string, error) {
	tmpl, err := template.New("configuration").Parse(configFileTemplate)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, c)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (c s3config) exportS3cmdConfig() (string, error) {
	configFileTemplate := `# Generated by scaleway-cli command
# Configuration file for s3cmd https://s3tools.org/s3cmd
# Default location: $HOME/.s3cfg
[default]
access_key = {{ .AccessKey }}
bucket_location = {{ .Region }}
host_base = s3.{{ .Region }}.scw.cloud
host_bucket = %(bucket)s.s3.{{ .Region }}.scw.cloud
secret_key = {{ .SecretKey }}
use_https = True`

	return renderTemplate(configFileTemplate, c)
}

func (c s3config) exportRcloneConfig() (string, error) {
	configFileTemplate := `# Generated by scaleway-cli command
# Configuration file for rclone https://rclone.org/s3/#scaleway
# Default location: $HOME/.config/rclone/rclone.conf 
[scaleway_{{ .Region }}]
type = s3
env_auth = false
endpoint = s3.{{ .Region }}.scw.cloud
access_key_id = {{ .AccessKey }}
secret_access_key = {{ .SecretKey }}
region = {{ .Region }}
location_constraint =
acl = private
force_path_style = false
server_side_encryption =
storage_class =`

	return renderTemplate(configFileTemplate, c)
}

func (c s3config) exportMcConfig() (string, error) {
	type hostconfig struct {
		URL       string `json:"url"`
		AccessKey string `json:"accessKey"`
		SecretKey string `json:"secretKey"`
		API       string `json:"api"`
	}
	type mcconfig struct {
		Version string                `json:"version"`
		Hosts   map[string]hostconfig `json:"hosts"`
	}
	m := mcconfig{
		Version: "9",
		Hosts: map[string]hostconfig{
			"scaleway_" + c.Region: {
				URL:       "https://s3." + c.Region + ".scw.cloud",
				AccessKey: c.AccessKey,
				SecretKey: c.SecretKey,
				API:       "S3v2",
			},
		},
	}
	res, err := json.Marshal(m)
	if err != nil {
		return "", nil
	}
	return string(res), nil
}
