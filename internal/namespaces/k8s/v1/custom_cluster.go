package k8s

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	clusterActionTimeout = 10 * time.Minute
)

//
// Marshalers
//

// clusterStatusMarshalerFunc marshals a k8s.ClusterStatus.
var (
	clusterStatusMarshalSpecs = human.EnumMarshalSpecs{
		k8s.ClusterStatusCreating:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
		k8s.ClusterStatusReady:        &human.EnumMarshalSpec{Attribute: color.FgGreen},
		k8s.ClusterStatusPoolRequired: &human.EnumMarshalSpec{Attribute: color.FgRed},
		k8s.ClusterStatusLocked:       &human.EnumMarshalSpec{Attribute: color.FgRed},
		k8s.ClusterStatusUpdating:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

const (
	clusterActionCreate = iota
	clusterActionUpdate
	clusterActionUpgrade
	clusterActionDelete
)

func clusterAvailableVersionsListBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		originalRes, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}

		listClusterAvailableVersionsResponse := originalRes.(*k8s.ListClusterAvailableVersionsResponse)
		return listClusterAvailableVersionsResponse.Versions, nil
	})

	return c
}

func clusterCreateBuilder(c *core.Command) *core.Command {
	type customCreateClusterRequest struct {
		*k8s.CreateClusterRequest
		Pools []*struct {
			*k8s.CreateClusterRequestPoolConfig
			Size *uint32
		}
	}

	c.ArgsType = reflect.TypeOf(customCreateClusterRequest{})

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (i interface{}, err error) {
		args := argsI.(*customCreateClusterRequest)

		request := args.CreateClusterRequest
		for _, pool := range args.Pools {
			poolReq := pool.CreateClusterRequestPoolConfig
			poolReq.Size = *pool.Size
			request.Pools = append(request.Pools, poolReq)
		}

		return runner(ctx, request)
	})

	c.WaitFunc = waitForClusterFunc(clusterActionCreate)
	return c
}

func clusterDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc(clusterActionDelete)
	return c
}

func clusterUpgradeBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc(clusterActionUpgrade)
	return c
}

func clusterUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForClusterFunc(clusterActionUpdate)
	return c
}

func waitForClusterFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		cluster, err := k8s.NewAPI(core.ExtractClient(ctx)).WaitForCluster(&k8s.WaitForClusterRequest{
			Region:    respI.(*k8s.Cluster).Region,
			ClusterID: respI.(*k8s.Cluster).ID,
			Timeout:   scw.TimeDurationPtr(clusterActionTimeout),
		})
		switch action {
		case clusterActionCreate:
			return cluster, err
		case clusterActionUpdate:
			return cluster, err
		case clusterActionUpgrade:
			return cluster, err
		case clusterActionDelete:
			if err != nil {
				// if we get a 404 here, it means the resource was successfully deleted
				notFoundError := &scw.ResourceNotFoundError{}
				responseError := &scw.ResponseError{}
				if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound || errors.As(err, &notFoundError) {
					return fmt.Sprintf("Cluster %s successfully deleted.", respI.(*k8s.Cluster).ID), nil
				}
			}
		}
		return nil, err
	}
}
