package registry

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	registry "github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type namespaceChargeArgs struct {
	NamespaceID string
	Region      scw.Region
}

func namespaceChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "charge",
		Short:     `List all the charges of this registry namespace since its creation`,
		ArgsType:  reflect.TypeFor[namespaceChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "namespace-id", Short: "Namespace ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*namespaceChargeArgs)

			api := registry.NewAPI(core.ExtractClient(ctx))
			namespace, err := api.GetNamespace(&registry.GetNamespaceRequest{
				Region:      args.Region,
				NamespaceID: args.NamespaceID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.NamespaceID, namespace.CreatedAt)
		},
	}
}
