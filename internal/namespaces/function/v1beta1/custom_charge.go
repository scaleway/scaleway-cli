package function

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type functionChargeArgs struct {
	FunctionID string
	Region     scw.Region
}

func functionChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "function",
		Resource:  "function",
		Verb:      "charge",
		Short:     `List all the charges of this function since its creation`,
		ArgsType:  reflect.TypeFor[functionChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "function-id", Short: "Function ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*functionChargeArgs)

			api := function.NewAPI(core.ExtractClient(ctx))
			f, err := api.GetFunction(&function.GetFunctionRequest{
				Region:     args.Region,
				FunctionID: args.FunctionID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.FunctionID, f.CreatedAt)
		},
	}
}
