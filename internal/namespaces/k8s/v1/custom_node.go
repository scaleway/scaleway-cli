package k8s

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	nodeActionTimeout = 10 * time.Minute
)

var (
	// nodeStatusMarshalSpecs allows to override the displayed status color
	nodeStatusMarshalSpecs = human.EnumMarshalSpecs{
		k8s.NodeStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.NodeStatusRebooting:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.NodeStatusReady:         &human.EnumMarshalSpec{Attribute: color.FgGreen},
		k8s.NodeStatusNotReady:      &human.EnumMarshalSpec{Attribute: color.FgYellow},
		k8s.NodeStatusCreationError: &human.EnumMarshalSpec{Attribute: color.FgRed},
		k8s.NodeStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)

const (
	nodeActionReboot = iota
)

func nodeRebootBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForNodeFunc(nodeActionReboot)
	return c
}

func waitForNodeFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		node, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForNode(&k8s.WaitForNodeRequest{
			Region:        respI.(*k8s.Node).Region,
			NodeID:        respI.(*k8s.Node).ID,
			Timeout:       scw.TimeDurationPtr(nodeActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		switch action {
		case nodeActionReboot:
			return node, err
		}
		return nil, err
	}
}

func k8sNodeWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a node to reach a stable state`,
		Long:      `Wait for a node to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the node.`,
		Namespace: "k8s",
		Resource:  "node",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(k8s.WaitForNodeRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := k8s.NewAPI(core.ExtractClient(ctx))
			return api.WaitForNode(&k8s.WaitForNodeRequest{
				Region:        argsI.(*k8s.WaitForNodeRequest).Region,
				NodeID:        argsI.(*k8s.WaitForNodeRequest).NodeID,
				Timeout:       scw.TimeDurationPtr(nodeActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "node-id",
				Short:      `ID of the node.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a node to reach a stable state",
				ArgsJSON: `{"node_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
