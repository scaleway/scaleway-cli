package messageq

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	messageq "github.com/scaleway/scaleway-sdk-go/api/messageq/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type deploymentCockpitArgs struct {
	DeploymentID string
	Region       scw.Region
}

func deploymentCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this MessageQ deployment`,
		ArgsType:  reflect.TypeFor[deploymentCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{Name: "deployment-id", Short: "Deployment ID", Required: true, Positional: true},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*deploymentCockpitArgs)

			api := messageq.NewAPI(core.ExtractClient(ctx))
			deployment, err := api.GetDeployment(&messageq.GetDeploymentRequest{
				Region:       args.Region,
				DeploymentID: args.DeploymentID,
			})
			if err != nil {
				return nil, err
			}

			return cockpit.OpenDashboard(
				ctx,
				"scw-messaging-nats-overview",
				map[string]string{
					"region":       string(deployment.Region),
					"nats_account": deployment.Name,
				},
			)
		},
	}
}
