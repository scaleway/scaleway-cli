package kafka

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	kafka "github.com/scaleway/scaleway-sdk-go/api/kafka/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterCockpitArgs struct {
	ClusterID string
	Region    scw.Region
}

func clusterCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "kafka",
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

			api := kafka.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&kafka.GetClusterRequest{
				Region:    args.Region,
				ClusterID: args.ClusterID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(ctx, "scw-apache-kafka-overview", map[string]string{
				"region":        string(cluster.Region),
				"resource_name": cluster.Name,
			})
		},
	}
}
