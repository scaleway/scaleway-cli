package redis

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterChargeArgs struct {
	ClusterID string
	Zone      scw.Zone
}

func clusterChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "redis",
		Resource:  "cluster",
		Verb:      "charge",
		Short:     `List all the charges of this cluster since its creation`,
		ArgsType:  reflect.TypeFor[clusterChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "cluster-id", Short: "Cluster ID", Required: true, Positional: true},
			core.ZoneArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*clusterChargeArgs)

			api := redis.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&redis.GetClusterRequest{
				ClusterID: args.ClusterID,
				Zone:      args.Zone,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ClusterID, cluster.CreatedAt)
		},
	}
}
