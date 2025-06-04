package instance

import (
	"context"
	"fmt"
	"net"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/validation"
)

// Builders

func ipCreateBuilder(c *core.Command) *core.Command {
	type customCreateIPRequest struct {
		*instance.CreateIPRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customCreateIPRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customCreateIPRequest)

			if args.CreateIPRequest == nil {
				args.CreateIPRequest = &instance.CreateIPRequest{}
			}
			request := args.CreateIPRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func ipListBuilder(c *core.Command) *core.Command {
	type customListIPsRequest struct {
		*instance.ListIPsRequest
		OrganizationID *string
		ProjectID      *string
	}

	renameOrganizationIDArgSpec(c.ArgSpecs)
	renameProjectIDArgSpec(c.ArgSpecs)

	c.ArgsType = reflect.TypeOf(customListIPsRequest{})

	c.AddInterceptors(
		func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
			args := argsI.(*customListIPsRequest)

			if args.ListIPsRequest == nil {
				args.ListIPsRequest = &instance.ListIPsRequest{}
			}
			request := args.ListIPsRequest
			request.Organization = args.OrganizationID
			request.Project = args.ProjectID

			return runner(ctx, request)
		},
	)

	return c
}

func ipAttachCommand() *core.Command {
	type customIPAttachRequest struct {
		OrganizationID *string
		ProjectID      *string
		// Server: UUID of the server you want to attach the IP to
		ServerID string   `json:"server,omitempty"`
		IP       string   `json:"-"`
		Zone     scw.Zone `json:"zone"`
	}

	return &core.Command{
		Short:     `Attach an IP to a given server`,
		Long:      `Attach an IP to a given server.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "attach",
		ArgsType:  reflect.TypeOf(customIPAttachRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := instance.NewAPI(core.ExtractClient(ctx))
			args := argsI.(*customIPAttachRequest)

			return api.UpdateIP(&instance.UpdateIPRequest{
				IP: args.IP,
				Server: &instance.NullableStringValue{
					Value: args.ServerID,
				},
				Zone: args.Zone,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip",
				Short:      `IP or UUID of the IP.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:     "server-id",
				Short:    "UUID of the server to attach the IP to",
				Required: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Attach an IP to the given server",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111", "ip": "1.2.3.4"}`,
			},
		},
	}
}

func ipDetachCommand() *core.Command {
	type customIPDetachRequest struct {
		OrganizationID *string
		ProjectID      *string
		IP             string   `json:"-"`
		Zone           scw.Zone `json:"zone"`
	}

	return &core.Command{
		Short:     `Detach an ip from its server`,
		Long:      `Detach an ip from its server.`,
		Namespace: "instance",
		Resource:  "ip",
		Verb:      "detach",
		ArgsType:  reflect.TypeOf(customIPDetachRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*customIPDetachRequest)
			api := instance.NewAPI(core.ExtractClient(ctx))

			var ipID string
			switch {
			case validation.IsUUID(args.IP):
				ipID = args.IP
			case net.ParseIP(args.IP) != nil:
				// Find the corresponding flexible IP UUID.
				logger.Debugf("finding public IP UUID from address: %s", args.IP)
				res, err := api.GetIP(&instance.GetIPRequest{
					Zone: args.Zone,
					IP:   args.IP,
				})
				if err != nil { // FIXME: isNotFoundError
					return nil, fmt.Errorf("%s does not belong to you", args.IP)
				}
				ipID = res.IP.ID
			default:
				return nil, fmt.Errorf(
					`invalid IP "%s", should be either an IP address ID or a reserved flexible IP address`,
					args.IP,
				)
			}

			return api.UpdateIP(&instance.UpdateIPRequest{
				Zone: args.Zone,
				// We detach an ip by specifying no server
				Server: &instance.NullableStringValue{
					Null: true,
				},
				IP: ipID,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "ip",
				Short:      `IP or UUID of the IP.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Detach an IP by using its UUID",
				ArgsJSON: `{"ip": "11111111-1111-1111-1111-111111111111"}`,
			},
			{
				Short:    "Detach an IP by using its IP address",
				ArgsJSON: `{"ip": "1.2.3.4"}`,
			},
		},
	}
}
