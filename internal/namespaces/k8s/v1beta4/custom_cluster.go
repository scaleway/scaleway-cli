package k8s

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"golang.org/x/xerrors"
)

const (
	clusterActionTimeout = 10 * time.Minute
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

func clusterCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc()
	return c
}

func clusterDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc()
	return c
}

func clusterUpgradeBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc()
	return c
}

func clusterUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc()
	return c
}

func waitForClusterFunc() core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		cluster, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForCluster(&k8s.WaitForClusterRequest{
			Region:    respI.(*k8s.Cluster).Region,
			ClusterID: respI.(*k8s.Cluster).ID,
			Timeout:   scw.DurationPtr(clusterActionTimeout),
		})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if xerrors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || xerrors.As(err, &notFoundError) {
				return fmt.Sprintf("Cluster %s successfully deleted.", respI.(*k8s.Cluster).ID), nil
			}
		}
		return cluster, nil
	}
}
