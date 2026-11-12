package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterCockpitArgs struct {
	ClusterID string
	Region    scw.Region
}

func clusterCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
		Resource:  "cluster",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this cluster`,
		ArgsType:  reflect.TypeFor[clusterCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "cluster-id", Short: "Cluster ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterCockpitArgs)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&k8s.GetClusterRequest{
				ClusterID: args.ClusterID,
				Region:    args.Region,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-kubernetes-cluster-overview",
				map[string]string{
					"region":       string(cluster.Region),
					"cluster_name": cluster.Name,
				},
			)
		},
	}
}
