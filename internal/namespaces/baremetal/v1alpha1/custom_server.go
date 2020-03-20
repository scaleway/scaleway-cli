package baremetal

import (
	"context"
	"fmt"
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
	Zone     scw.Zone
	ServerID string
}

var serverActionArgSpecs = core.ArgSpecs{
	{
		Name:     "server-id",
		Short:    `ID of the server affected by the action.`,
		Required: true,
	},
	core.ZoneArgSpec(),
}

func waitForServerFunc() core.WaitFunc {
	return func(ctx context.Context, argsI, _ interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			ServerID: argsI.(*baremetalActionRequest).ServerID,
			Zone:     argsI.(*baremetalActionRequest).Zone,
			Timeout:  serverActionTimeout,
		})
	}
}

func serverWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for server to reach a stable state`,
		Long:      `Wait for server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(baremetalActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			return waitForServerFunc()(ctx, argsI, nil)
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

func serverStartCommand() *core.Command {
	return &core.Command{
		Short:     `Power on server`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "start",
		ArgsType:  reflect.TypeOf(baremetalActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*baremetalActionRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)

			_, err := api.StartServer(&baremetal.StartServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			return &core.SuccessResult{Message: fmt.Sprintf("%s successfully started", args.ServerID)}, err
		},
		WaitFunc: waitForServerFunc(),
		ArgSpecs: serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Start a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverStopCommand() *core.Command {
	return &core.Command{
		Short:     `Power off server`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "stop",
		ArgsType:  reflect.TypeOf(baremetalActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*baremetalActionRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)

			_, err := api.StopServer(&baremetal.StopServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			return &core.SuccessResult{Message: fmt.Sprintf("%s successfully stopped", args.ServerID)}, err
		},
		WaitFunc: waitForServerFunc(),
		ArgSpecs: serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Stop a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func serverRebootCommand() *core.Command {
	return &core.Command{
		Short:     `Reboot server`,
		Namespace: "baremetal",
		Resource:  "server",
		Verb:      "reboot",
		ArgsType:  reflect.TypeOf(baremetalActionRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
			args := argsI.(*baremetalActionRequest)

			client := core.ExtractClient(ctx)
			api := baremetal.NewAPI(client)

			_, err := api.RebootServer(&baremetal.RebootServerRequest{
				Zone:     args.Zone,
				ServerID: args.ServerID,
			})
			return &core.SuccessResult{Message: fmt.Sprintf("%s successfully rebooted", args.ServerID)}, err
		},
		WaitFunc: waitForServerFunc(),
		ArgSpecs: serverActionArgSpecs,
		Examples: []*core.Example{
			{
				Short:   "Reboot a server in the default zone with a given id",
				Request: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
