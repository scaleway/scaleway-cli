package baremetal

import (
	"context"
	"errors"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	fip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type serverAddFlexibleIPRequest struct {
	ServerID    string
	Description string
	IPType      string
	Tags        []string
	Zone        scw.Zone
}

func serverAddFlexibleIP() *core.Command {
	return &core.Command{
		Short:     "Attach a new flexible IP to a server",
		Long:      "Create and attach a new flexible IP to a server",
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "add-flexible-ip",
		Groups:    []string{"utility"},
		ArgsType:  reflect.TypeOf(serverAddFlexibleIPRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server to which the newly created flexible IP will be attached.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Flexible IP description (max. of 255 characters)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ip-type",
				Short:      `Define whether the flexible IP is an IPv4 or IPv6`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: ipTypeOption,
				ValidateFunc: func(_ *core.ArgSpec, value any) error {
					if value == "IPv4" || value == "IPv6" || value == "" {
						return nil
					}

					return &core.CliError{
						Err:  errors.New("error looks like the IP type isn't correct"),
						Hint: "Two type available : IPv6 or IPv4",
					}
				},
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to associate to the flexible IP`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ZoneArgSpec(scw.ZoneFrPar1, scw.ZoneFrPar2, scw.ZoneNlAms1),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			request := argsI.(*serverAddFlexibleIPRequest)
			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)
			server, err := api.GetServer(&baremetal.GetServerRequest{
				Zone:     request.Zone,
				ServerID: request.ServerID,
			}, scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}
			args := request
			if args.IPType == "" {
				args, err = promptIPFlexibleServer(ctx, request)
			}
			if err != nil {
				return nil, err
			}
			apiFip := fip.NewAPI(client)
			IsIPv6 := args.IPType == "IPv6"
			flexibleIP, err := apiFip.CreateFlexibleIP(&fip.CreateFlexibleIPRequest{
				ServerID:    &server.ID,
				Zone:        server.Zone,
				ProjectID:   server.ProjectID,
				Description: args.Description,
				Tags:        args.Tags,
				IsIPv6:      IsIPv6,
			})
			if err != nil {
				return nil, err
			}

			return apiFip.WaitForFlexibleIP(&fip.WaitForFlexibleIPRequest{
				FipID:         flexibleIP.ID,
				Zone:          flexibleIP.Zone,
				Timeout:       scw.TimeDurationPtr(flexibleip.FlexibleIPTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
	}
}
