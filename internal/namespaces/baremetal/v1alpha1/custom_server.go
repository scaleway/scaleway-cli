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

func waitForServerFunc() core.WaitFunc {
	return func(ctx context.Context, argsI, _ interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			ServerID: argsI.(*baremetalActionRequest).ServerID,
			Zone:     argsI.(*baremetalActionRequest).Zone,
			Timeout:  serverActionTimeout,
		})
	}
}
