package k8s

import (
	"context"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	billingv2beta1 "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type poolChargeArgs struct {
	PoolID string
	Region scw.Region
}

// Kapsule pool nodes are billed as instance servers, recorded under the
// underlying instance server ID (embedded in the node's ProviderID), so the
// pool's charges are aggregated from each node's server ID.
//
// ponytail: servers of auto-replaced nodes are no longer reachable from the
// k8s API, so their historical charges are missed. Full history would need a
// billing-side "pool" resource key.
func poolChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "k8s",
		Resource:  "pool",
		Verb:      "charge",
		Short:     `List all the charges of this pool since its creation`,
		ArgsType:  reflect.TypeFor[poolChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "pool-id", Short: "Pool ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*poolChargeArgs)

			api := k8s.NewAPI(core.ExtractClient(ctx))
			pool, err := api.GetPool(&k8s.GetPoolRequest{
				PoolID: args.PoolID,
				Region: args.Region,
			})
			if err != nil {
				return nil, err
			}

			nodes, err := api.ListNodes(&k8s.ListNodesRequest{
				Region:    args.Region,
				ClusterID: pool.ClusterID,
				PoolID:    &args.PoolID,
			}, scw.WithAllPages(), scw.WithContext(ctx))
			if err != nil {
				return nil, err
			}

			charges := []*billingv2beta1.Charge(nil)
			for _, node := range nodes.Nodes {
				serverID := serverIDFromProviderID(node.ProviderID)
				if serverID == "" {
					continue
				}
				nodeCharges, err := billing.ListChargesForResource(ctx, serverID, node.CreatedAt)
				if err != nil {
					return nil, err
				}
				charges = append(charges, nodeCharges...)
			}

			return charges, nil
		},
	}
}

// serverIDFromProviderID extracts the instance server ID from a node
// ProviderID of the form "scaleway://instance/<zone>/<server-id>".
func serverIDFromProviderID(providerID string) string {
	parts := strings.Split(strings.TrimPrefix(providerID, "scaleway://"), "/")
	if len(parts) == 3 && parts[0] == "instance" {
		return parts[2]
	}

	return ""
}
