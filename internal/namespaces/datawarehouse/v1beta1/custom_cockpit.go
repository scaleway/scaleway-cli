package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type deploymentCockpitArgs struct {
	DeploymentID string
	Region       scw.Region
}

func deploymentCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this data warehouse`,
		ArgsType:  reflect.TypeFor[deploymentCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "deployment-id", Short: "Deployment ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*deploymentCockpitArgs)

			api := datawarehouse.NewAPI(core.ExtractClient(ctx))
			deployment, err := api.GetDeployment(&datawarehouse.GetDeploymentRequest{
				Region:       args.Region,
				DeploymentID: args.DeploymentID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-clickhouse-data-warehouse-overview",
				map[string]string{
					"region":         string(deployment.Region),
					"data_warehouse": deployment.Name,
				},
			)
		},
	}
}
