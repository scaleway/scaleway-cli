package object

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type bucketChargeArgs struct {
	Name   string
	Region scw.Region
}

func bucketChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "object",
		Resource:  "bucket",
		Verb:      "charge",
		Short:     `List all the charges of this bucket since its creation`,
		ArgsType:  reflect.TypeFor[bucketChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "name", Short: "Bucket name", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*bucketChargeArgs)

			// Buckets are identified by name and have no created_at: query the full history.
			return billing.ListChargesForResource(ctx, args.Name, nil)
		},
	}
}
