package baremetal

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	ServerActionTimeout = 20 * time.Minute
)

var serverPingStatusMarshalSpecs = human.EnumMarshalSpecs{
	baremetal.ServerPingStatusPingStatusDown: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "down",
	},
	baremetal.ServerPingStatusPingStatusUp: &human.EnumMarshalSpec{
		Attribute: color.FgGreen,
		Value:     "up",
	},
	baremetal.ServerPingStatusPingStatusUnknown: &human.EnumMarshalSpec{
		Attribute: color.Faint,
		Value:     "unknown",
	},
}

func serverWaitCommand() *core.Command {
	type serverWaitRequest struct {
		ServerID string
		Zone     scw.Zone
		Timeout  time.Duration
	}

	return &core.Command{
		Short:     `Wait for a server to reach a stable state (delivery and installation)`,
		Long:      `Wait for a server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "wait",
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeOf(serverWaitRequest{}),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			args := argsI.(*serverWaitRequest)

			api := baremetal.NewAPI(core.ExtractClient(ctx))
			logger.Debugf("starting to wait for server to reach a stable delivery status")
			server, err := api.WaitForServer(&baremetal.WaitForServerRequest{
				ServerID:      args.ServerID,
				Zone:          args.Zone,
				Timeout:       scw.TimeDurationPtr(args.Timeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}
			if server.Status != baremetal.ServerStatusReady {
				return nil, &core.CliError{
					Err:     errors.New("server did not reach a stable delivery status"),
					Details: fmt.Sprintf("server %s is in %s status", server.ID, server.Status),
				}
			}
			if server.Install == nil {
				return server, nil
			}

			logger.Debugf(
				"server reached a stable delivery status. Will now starting to wait for server to reach a stable installation status",
			)
			server, err = api.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
				ServerID:      args.ServerID,
				Zone:          args.Zone,
				Timeout:       scw.TimeDurationPtr(args.Timeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}
			if server.Install.Status != baremetal.ServerInstallStatusCompleted {
				return nil, &core.CliError{
					Err: fmt.Errorf(
						"server %s did not reach a stable installation status",
						server.ID,
					),
					Details: fmt.Sprintf("server %s is in %s status", server.ID, server.Status),
				}
			}
			logger.Debugf("server reached a stable installation status")

			return server, nil
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "server-id",
				Short:      `ID of the server affected by the action.`,
				Required:   true,
				Positional: true,
			},
			core.ZoneArgSpec(),
			core.WaitTimeoutArgSpec(ServerActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a server to reach a stable state",
				ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

// serverStartBuilder overrides the baremetalServerStart command
func serverStartBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))

		return api.WaitForServer(&baremetal.WaitForServerRequest{
			Zone:          argsI.(*baremetal.StartServerRequest).Zone,
			ServerID:      respI.(*baremetal.Server).ID,
			Timeout:       scw.TimeDurationPtr(ServerActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

// serverStopBuilder overrides the baremetalServerStop command
func serverStopBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))

		return api.WaitForServer(&baremetal.WaitForServerRequest{
			Zone:          argsI.(*baremetal.StopServerRequest).Zone,
			ServerID:      respI.(*baremetal.Server).ID,
			Timeout:       scw.TimeDurationPtr(ServerActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

// serverRebootBuilder overrides the baremetalServerReboot command
func serverRebootBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("boot-type").Default = core.DefaultValueSetter("normal")

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))

		return api.WaitForServer(&baremetal.WaitForServerRequest{
			Zone:          argsI.(*baremetal.RebootServerRequest).Zone,
			ServerID:      respI.(*baremetal.Server).ID,
			Timeout:       scw.TimeDurationPtr(ServerActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}

type customServer struct {
	baremetal.Server
	OfferName string `json:"offer_name"`
}

func serverListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "Name",
				FieldName: "Name",
			},
			{
				Label:     "Offer Name",
				FieldName: "OfferName",
			},
			{
				Label:     "Status",
				FieldName: "Status",
			},
			{
				Label:     "Tags",
				FieldName: "Tags",
			},
			{
				Label:     "Ping",
				FieldName: "PingStatus",
			},
		},
	}

	// We add the server type to the list sent to the user
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		listServerResp, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		client := core.ExtractClient(ctx)
		api := baremetal.NewAPI(client)
		listOffers, err := api.ListOffers(&baremetal.ListOffersRequest{
			Zone: argsI.(*baremetal.ListServersRequest).Zone,
		}, scw.WithAllPages())
		if err != nil {
			return listServerResp, err
		}
		offerNameByID := make(map[string]string)
		for _, offer := range listOffers.Offers {
			offerNameByID[offer.ID] = offer.Name
		}

		var customRes []customServer
		for _, server := range listServerResp.([]*baremetal.Server) {
			customRes = append(customRes, customServer{
				Server:    *server,
				OfferName: offerNameByID[server.OfferID],
			})
		}

		return customRes, nil
	}

	return c
}

func serverMarshalerFunc(i any, opt *human.MarshalOpt) (string, error) {
	type tmp baremetal.Server
	baremetalServer := tmp(i.(baremetal.Server))
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "IPs",
			Title:     "IPs",
		},
		{
			FieldName:   "Options",
			Title:       "Options",
			HideIfEmpty: true,
		},
		{
			FieldName:   "Install",
			Title:       "Install",
			HideIfEmpty: true,
		},
	}
	str, err := human.Marshal(baremetalServer, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}
