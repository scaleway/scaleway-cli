package k8s

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type clusterChargeArgs struct {
	ClusterID string
	Region    scw.Region
}

func clusterChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
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

			api := k8s.NewAPI(core.ExtractClient(ctx))
			cluster, err := api.GetCluster(&k8s.GetClusterRequest{
				ClusterID: args.ClusterID,
				Region:    args.Region,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.ClusterID, cluster.CreatedAt)
		},
	}
}
