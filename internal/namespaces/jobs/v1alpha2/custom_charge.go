package jobs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	billing "github.com/scaleway/scaleway-cli/v2/internal/namespaces/billing/v2beta1"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type definitionChargeArgs struct {
	JobDefinitionID string
	Region          scw.Region
}

func definitionChargeCommand() *core.Command {
	return &core.Command{
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "charge",
		Short:     `List all the charges of this job definition since its creation`,
		ArgsType:  reflect.TypeFor[definitionChargeArgs](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      "Job Definition ID",
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(),
		},
		View: billing.ChargesView(),
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*definitionChargeArgs)

			api := jobs.NewAPI(core.ExtractClient(ctx))
			d, err := api.GetJobDefinition(&jobs.GetJobDefinitionRequest{
				Region:          args.Region,
				JobDefinitionID: args.JobDefinitionID,
			})
			if err != nil {
				return nil, err
			}

			return billing.ListChargesForResource(ctx, args.JobDefinitionID, d.CreatedAt)
		},
	}
}
