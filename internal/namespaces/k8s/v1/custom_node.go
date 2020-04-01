package k8s

import (
	"context"
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
	// nodeStatusAttributes allows to override the displayed status color
	nodeStatusAttributes = human.Attributes{
		k8s.NodeStatusCreating:      color.FgBlue,
		k8s.NodeStatusRebooting:     color.FgBlue,
		k8s.NodeStatusReady:         color.FgGreen,
		k8s.NodeStatusNotReady:      color.FgYellow,
		k8s.NodeStatusCreationError: color.FgRed,
		k8s.NodeStatusLocked:        color.FgRed,
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
			Region:  respI.(*k8s.Node).Region,
			NodeID:  respI.(*k8s.Node).ID,
			Timeout: scw.TimeDurationPtr(nodeActionTimeout),
		})
		switch action {
		case nodeActionReboot:
			return node, err
		}
		return nil, err
	}
}
