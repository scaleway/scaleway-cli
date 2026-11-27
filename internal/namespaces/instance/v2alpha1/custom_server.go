package instance

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var serverStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerStatusStarted:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.ServerStatusStarting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusStopping:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusPausing:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusRebooting:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	instance.ServerStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	instance.ServerStatusStopped:       &human.EnumMarshalSpec{Attribute: color.Faint},
	instance.ServerStatusPaused:        &human.EnumMarshalSpec{Attribute: color.Faint},
	instance.ServerStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
}

var serverPublicNetworkInterfaceStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerPublicNetworkInterfaceStatusAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.ServerPublicNetworkInterfaceStatusSyncing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPublicNetworkInterfaceStatusUnknownStatus: &human.EnumMarshalSpec{
		Attribute: color.Faint,
	},
}

var serverPrivateNetworkInterfaceStatusMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerPrivateNetworkInterfaceStatusAvailable: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
	},
	instance.ServerPrivateNetworkInterfaceStatusAttaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDetaching: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusSyncing: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDeleting: &human.EnumMarshalSpec{
		Attribute: color.FgBlue,
	},
	instance.ServerPrivateNetworkInterfaceStatusDetachError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerPrivateNetworkInterfaceStatusDeleteError: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
	},
	instance.ServerPrivateNetworkInterfaceStatusUnknownStatus: &human.EnumMarshalSpec{
		Attribute: color.Faint,
	},
}

// serverMarshalerFunc marshals a v2 instance.Server
func serverMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp instance.Server
	server := tmp(i.(instance.Server))

	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "Volumes",
			Title:     "Volumes",
		},
		{
			FieldName: "PublicNetworkInterface",
			Title:     "Public Network Interface",
		},
		{
			FieldName:   "PrivateNetworkInterfaces",
			Title:       "Private Network Interfaces",
			HideIfEmpty: true,
		},
		{
			FieldName:   "Filesystems",
			Title:       "Server Filesystems",
			HideIfEmpty: true,
		},
		{
			FieldName:   "Tags",
			Title:       "Tags",
			HideIfEmpty: true,
		},
		{
			FieldName: "RuntimeInfo",
			Title:     "Runtime Info",
		},
	}

	// Replace empty lists by nil for HideIfEmpty to properly detect them
	if len(server.PrivateNetworkInterfaces) == 0 {
		server.PrivateNetworkInterfaces = nil
	}

	if len(server.Filesystems) == 0 {
		server.Filesystems = nil
	}

	if len(server.Tags) == 0 {
		server.Tags = nil
	}

	str, err := human.Marshal(server, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

//
// Builders
//

func serverStartSpotBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		return instance.NewAPI(core.ExtractClient(ctx)).
			WaitForServer(&instance.WaitForServerRequest{
				Zone:          respI.(*instance.Server).Zone,
				ServerID:      respI.(*instance.Server).ID,
				Timeout:       new(serverActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
	}

	return c
}

//
// Custom commands
//

type serverWaitRequest struct {
	Zone     scw.Zone
	ServerID string
	Timeout  time.Duration
}

func serverWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for server to reach a stable state`,
		Long:      `Wait for server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "instance",
		Resource:  "server",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeFor[serverWaitRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*serverWaitRequest)

			return instance.NewAPI(core.ExtractClient(ctx)).
				WaitForServer(&instance.WaitForServerRequest{
					Zone:          args.Zone,
					ServerID:      args.ServerID,
					Timeout:       new(args.Timeout),
					RetryInterval: core.DefaultRetryInterval,
				})
		},
		ArgSpecs: core.ArgSpecs{
			core.WaitTimeoutArgSpec(serverActionTimeout),
			{
				Name:       "server-id",
				Short:      `ID of the server affected by the action.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec((*instance.API)(nil).Zones()...),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a server to reach a stable state",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverWaitFunc(ctx context.Context, argsI, respI any) (any, error) {
	timeout := new(serverActionTimeout)
	retryInterval := core.DefaultRetryInterval

	switch waitReq := argsI.(type) {
	case *serverWaitRequest:
		timeout = &waitReq.Timeout
	case *instance.WaitForServerRequest:
		if waitReq.Timeout != nil {
			timeout = waitReq.Timeout
		}

		if waitReq.RetryInterval != nil {
			retryInterval = waitReq.RetryInterval
		}
	}

	resp := respI.(*instance.Server)

	return waitForServer(
		ctx,
		resp.ID,
		resp.Zone,
		timeout,
		retryInterval,
	)
}

func waitForServer(
	ctx context.Context,
	serverID string,
	zone scw.Zone,
	timeout, retryInterval *time.Duration,
) (*instance.Server, error) {
	req := &instance.WaitForServerRequest{
		ServerID:      serverID,
		Zone:          zone,
		Timeout:       timeout,
		RetryInterval: retryInterval,
	}

	return instance.NewAPI(core.ExtractClient(ctx)).
		WaitForServer(req, scw.WithContext(ctx))
}
