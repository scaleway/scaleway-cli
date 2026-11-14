package kafka

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	kafka "github.com/scaleway/scaleway-sdk-go/api/kafka/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterChargeArgs struct {
	ClusterID string
	Region    scw.Region
}

func clusterChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "kafka",
		Resource:  "cluster",
		Verb:      "charge",
		Short:     `List all the charges of this cluster since its creation`,
		ArgsType:  reflect.TypeFor[clusterChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "cluster-id", Short: "Cluster ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterChargeArgs)

			api := kafka.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&kafka.GetClusterRequest{
				Region:    args.Region,
				ClusterID: args.ClusterID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ClusterID, cluster.CreatedAt)
		},
	}
}
