package k8s

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
)

const (
	nodeActionTimeout = 10 * time.Minute
)

// nodeStatusFormat allows to override the displayed status color
var nodeStatusFormat = humanv2.EnumFormats[k8s.NodeStatus]{
	k8s.NodeStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
	k8s.NodeStatusRebooting:     humanv2.EnumFormat{Color: color.FgBlue},
	k8s.NodeStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
	k8s.NodeStatusNotReady:      humanv2.EnumFormat{Color: color.FgYellow},
	k8s.NodeStatusCreationError: humanv2.EnumFormat{Color: color.FgRed},
	k8s.NodeStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
}

const (
	nodeActionReboot = iota
)

func nodeRebootBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForNodeFunc(nodeActionReboot)

	return c
}

func waitForNodeFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI any) (any, error) {
		node, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForNode(&k8s.WaitForNodeRequest{
			Region:        respI.(*k8s.Node).Region,
			NodeID:        respI.(*k8s.Node).ID,
			Timeout:       new(nodeActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if action == nodeActionReboot {
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
		Groups:    []string{"workflow"},
		ArgsType:  reflect.TypeFor[k8s.WaitForNodeRequest](),
		Run: func(ctx context.Context, argsI any) (i any, err error) {
			api := k8s.NewAPI(core.ExtractClient(ctx))

			return api.WaitForNode(&k8s.WaitForNodeRequest{
				Region:        argsI.(*k8s.WaitForNodeRequest).Region,
				NodeID:        argsI.(*k8s.WaitForNodeRequest).NodeID,
				Timeout:       argsI.(*k8s.WaitForNodeRequest).Timeout,
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
			core.WaitTimeoutArgSpec(nodeActionTimeout),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a node to reach a stable state",
				ArgsJSON: `{"node_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
