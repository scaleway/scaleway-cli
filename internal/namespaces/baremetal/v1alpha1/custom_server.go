package baremetal

import (
	"context"
	"reflect"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	serverActionTimeout = 20 * time.Minute
)

type baremetalActionRequest struct {
	ServerID string
	Zone     scw.Zone
}

var serverActionArgSpecs = core.ArgSpecs{
	{
		Name:     "server-id",
		Short:    `ID of the server affected by the action.`,
		Required: true,
	},
	core.ZoneArgSpec(),
}

func serverWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a server to reach a stable state`,
		Long:      `Wait for a server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(baremetalActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := baremetal.NewAPI(core.ExtractClient(ctx))
			return api.WaitForServer(&baremetal.WaitForServerRequest{
				ServerID: argsI.(*baremetalActionRequest).ServerID,
				Zone:     argsI.(*baremetalActionRequest).Zone,
				Timeout:  serverActionTimeout,
			})
		},
		ArgSpecs: serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Wait for a server to reach a stable state",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

// serverStartBuilder overrides the baremetalServerStart command
func serverStartBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			Zone:     argsI.(*baremetal.StartServerRequest).Zone,
			ServerID: respI.(*baremetal.StartServerRequest).ServerID,
			Timeout:  serverActionTimeout,
		})
	}

	return c
}

// serverStopBuilder overrides the baremetalServerStop command
func serverStopBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			Zone:     argsI.(*baremetal.StopServerRequest).Zone,
			ServerID: respI.(*baremetal.StopServerRequest).ServerID,
			Timeout:  serverActionTimeout,
		})
	}

	return c
}

// serverRebootBuilder overrides the baremetalServerReboot command
func serverRebootBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			Zone:     argsI.(*baremetal.RebootServerRequest).Zone,
			ServerID: respI.(*baremetal.RebootServerRequest).ServerID,
			Timeout:  serverActionTimeout,
		})
	}

	return c
}
