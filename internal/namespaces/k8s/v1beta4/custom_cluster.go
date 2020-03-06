package k8s

import (
	"context"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
)

//
// Marshalers
//

// clusterStatusMarshalerFunc marshals a k8s.ClusterStatus.
var (
	clusterStatusAttributes = human.Attributes{
		k8s.ClusterStatusCreating: color.FgBlue,
		k8s.ClusterStatusReady:    color.FgGreen,
		k8s.ClusterStatusError:    color.FgRed,
		k8s.ClusterStatusLocked:   color.FgRed,
		k8s.ClusterStatusUpdating: color.FgBlue,
		k8s.ClusterStatusWarning:  color.FgHiYellow,
	}
)

func clusterAvailableVersionsListBuilder(c *core.Command) *core.Command {
	originalRun := c.Run

	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		originalRes, err := originalRun(ctx, argsI)
		if err != nil {
			return nil, err
		}

		listClusterAvailableVersionsResponse := originalRes.(*k8s.ListClusterAvailableVersionsResponse)
		return listClusterAvailableVersionsResponse.Versions, nil
	}

	return c
}
