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

func clusterMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	type tmp k8s.Cluster
	cluster := tmp(i.(k8s.Cluster))

	// Sections
	opt.Sections = []*human.MarshalSection{
		{
			FieldName: "AutoscalerConfig",
			Title:     "Autoscaler configuration",
		},
		{
			FieldName: "AutoUpgrade",
			Title:     "Auto-upgrade settings",
		},
		{
			FieldName: "OpenIDConnectConfig",
			Title:     "Open ID Connect configuration",
		},
	}

	str, err := human.Marshal(cluster, opt)
	if err != nil {
		return "", err
	}

	return str, nil
}

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
	c.WaitFunc = waitForClusterFunc(clusterActionCreate)

	c.ArgSpecs.GetByName("cni").Default = core.DefaultValueSetter("cilium")
	c.ArgSpecs.GetByName("version").Default = core.DefaultValueSetter("latest")

	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		args := argsI.(*k8s.CreateClusterRequest)

		// Handle default latest version for k8s cluster
		if args.Version == "latest" {
			latestVersion, err := getLatestK8SVersion(core.ExtractClient(ctx))
			if err != nil {
				return nil, fmt.Errorf("could not retrieve latest K8S version")
			}
			args.Version = latestVersion
		}

		return runner(ctx, args)
	}

	return c
}

func clusterGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		res, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		cluster := res.(*k8s.Cluster)

		args := argsI.(*k8s.GetClusterRequest)
		pools, err := k8s.NewAPI(core.ExtractClient(ctx)).ListPools(&k8s.ListPoolsRequest{
			Region:    args.Region,
			ClusterID: args.ClusterID,
		})
		if err != nil {
			return res, err
		}

		type customPool struct {
			ID          string         `json:"id"`
			Name        string         `json:"name"`
			Status      k8s.PoolStatus `json:"status"`
			Version     string         `json:"version"`
			NodeType    string         `json:"node_type"`
			MinSize     uint32         `json:"min_size"`
			Size        uint32         `json:"size"`
			MaxSize     uint32         `json:"max_size"`
			Autoscaling bool           `json:"autoscaling"`
			Autohealing bool           `json:"autohealing"`
		}

		customPools := []customPool{}

		for _, pool := range pools.Pools {
			customPools = append(customPools, customPool{
				ID:          pool.ID,
				Name:        pool.Name,
				Status:      pool.Status,
				Version:     pool.Version,
				NodeType:    pool.NodeType,
				MinSize:     pool.MinSize,
				Size:        pool.Size,
				MaxSize:     pool.MaxSize,
				Autoscaling: pool.Autoscaling,
				Autohealing: pool.Autohealing,
			})
		}

		return struct {
			*k8s.Cluster
			Pools []customPool `json:"pools"`
		}{
			cluster,
			customPools,
		}, nil
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				FieldName: "AutoscalerConfig",
				Title:     "Autoscaler configuration",
			},
			{
				FieldName: "AutoUpgrade",
				Title:     "Auto-upgrade settings",
			},
			{
				FieldName: "OpenIDConnectConfig",
				Title:     "Open ID Connect configuration",
			},
			{
				FieldName: "Pools",
				Title:     "Pools",
			},
		},
	}

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
			Region:        respI.(*k8s.Cluster).Region,
			ClusterID:     respI.(*k8s.Cluster).ID,
			Timeout:       scw.TimeDurationPtr(clusterActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
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

func k8sClusterWaitCommand() *core.Command {
	type customClusterWaitArgs struct {
		k8s.WaitForClusterRequest
		WaitForPools bool
	}
	return &core.Command{
		Short:     `Wait for a cluster to reach a stable state`,
		Long:      `Wait for server to reach a stable state. This is similar to using --wait flag on other action commands, but without requiring a new action on the server.`,
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(customClusterWaitArgs{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			args := argsI.(*customClusterWaitArgs)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.WaitForCluster(&k8s.WaitForClusterRequest{
				Region:        args.Region,
				ClusterID:     args.ClusterID,
				Timeout:       scw.TimeDurationPtr(clusterActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
			if err != nil {
				return nil, err
			}

			if args.WaitForPools {
				pools, err := api.ListPools(&k8s.ListPoolsRequest{
					Region:    cluster.Region,
					ClusterID: cluster.ID,
				}, scw.WithAllPages())
				if err != nil {
					return cluster, err
				}
				for _, pool := range pools.Pools {
					_, err := api.WaitForPool(&k8s.WaitForPoolRequest{
						Region:        pool.Region,
						PoolID:        pool.ID,
						Timeout:       scw.TimeDurationPtr(poolActionTimeout),
						RetryInterval: core.DefaultRetryInterval,
					})
					if err != nil {
						return cluster, err
					}
				}
			}

			return cluster, nil
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cluster-id",
				Short:      `ID of the cluster.`,
				Required:   true,
				Positional: true,
			},
			{
				Name:  "wait-for-pools",
				Short: "Wait for pools to be ready.",
			},
			core.RegionArgSpec(),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a cluster to reach a stable state",
				ArgsJSON: `{"cluster_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}
