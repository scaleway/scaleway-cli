package function

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type functionCockpitArgs struct {
	FunctionID string
	Region     scw.Region
}

func functionCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "function",
		Resource:  "function",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this function`,
		ArgsType:  reflect.TypeFor[functionCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "function-id", Short: "Function ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*functionCockpitArgs)

			api := function.NewAPI(core.ExtractClient(ctx))
			f, err := api.GetFunction(&function.GetFunctionRequest{
				Region:     args.Region,
				FunctionID: args.FunctionID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-serverless-functions-overview",
				map[string]string{
					"region":        string(f.Region),
					"function_name": f.Name,
				},
			)
		},
	}
}
