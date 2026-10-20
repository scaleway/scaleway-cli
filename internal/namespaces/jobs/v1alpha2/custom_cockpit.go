package jobs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type definitionCockpitArgs struct {
	JobDefinitionID string
	Region          scw.Region
}

func definitionCockpitCommand() *core.Command {
	return &core.Command{
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "cockpit",
		Short:     `Open the Cockpit Grafana dashboard for this job definition`,
		ArgsType:  reflect.TypeFor[definitionCockpitArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      "Job Definition ID",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*definitionCockpitArgs)

			api := jobs.NewAPI(core.ExtractClient(ctx))
			d, err := api.GetJobDefinition(&jobs.GetJobDefinitionRequest{
				Region:          args.Region,
				JobDefinitionID: args.JobDefinitionID,
			})
			if err != nil {
				return nil, err
			}

			// This dashboard has no region variable (it scopes by datasource).
			return cockpit.OpenDashboard(
				ctx,
				"scw-serverless-jobs-overview",
				map[string]string{
					"job_definition_name": d.Name,
				},
			)
		},
	}
}
