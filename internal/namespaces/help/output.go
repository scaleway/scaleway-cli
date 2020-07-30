package help

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func helpOutput() *core.Command {
	longText := `Output formatting in the CLI

Human output

- Standard human output

	scw instance server list                       
	ID                                    NAME               TYPE    STATE    ZONE      PUBLIC IP
	088b01da-9ba7-40d2-bc55-eb3170f42185  scw-cool-franklin  DEV1-S  running  fr-par-1  51.15.251.251

- Human with column selection

You can select the columns that you want to print with commands that return a list 

	scw instance server list -o human=Name,PublicIP
	NAME                                            PUBLIC IP
	scw-cool-franklin                               51.15.251.251

JSON output

- Standard JSON output

	scw config dump -o json
	{"access_key":"SCWXXXXXXXXXXXXXXXXX","secret_key":"11111111-1111-1111-1111-111111111111","default_organization_id":"11111111-1111-1111-1111-111111111111","default_region":"fr-par","default_zone":"fr-par-1","send_telemetry":true}

- Pretty JSON output

	scw config dump -o json=pretty
	{
	  "access_key": "SCWXXXXXXXXXXXXXXXXX",
	  "secret_key": "11111111-1111-1111-1111-111111111111",
	  "default_organization_id": "11111111-1111-1111-1111-111111111111",
	  "default_region": "fr-par",
	  "default_zone": "fr-par-1",
	  "send_telemetry": true
	}

Template output

You can use Go template to manipulate the output of a command and create a custom rendering of your resources. 
Visit https://golang.org/pkg/text/template/ to learn more about Go template format.

	scw instance server list -o template="{{ .Name }}---{{ .Status }}"
	TODO, show example here

`
	return &core.Command{
		Short:                "Get help about how the CLI output works",
		Long:                 longText,
		Namespace:            "help",
		Resource:             "output",
		AllowAnonymousClient: true,
		ArgsType:             reflect.TypeOf(struct{}{}),
		ArgSpecs:             core.ArgSpecs{},
		Run: func(ctx context.Context, argsI interface{}) (interface{}, error) {
			return longText, nil
		},
	}
}
