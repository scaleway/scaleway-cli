package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type deploymentChargeArgs struct {
	DeploymentID string
	Region       scw.Region
}

func deploymentChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "charge",
		Short:     `List all the charges of this data warehouse since its creation`,
		ArgsType:  reflect.TypeFor[deploymentChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "deployment-id", Short: "Deployment ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*deploymentChargeArgs)

			api := datawarehouse.NewAPI(core.ExtractClient(ctx))
			deployment, err := api.GetDeployment(&datawarehouse.GetDeploymentRequest{
				Region:       args.Region,
				DeploymentID: args.DeploymentID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.DeploymentID, deployment.CreatedAt)
		},
	}
}
