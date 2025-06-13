package jobs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func jobsRunWait() *core.Command {
	return &core.Command{
		Short:     `Wait for a job run to reach a stable state`,
		Long:      `Wait for a job run to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "wait",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.WaitForJobRunRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-run-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.WaitForJobRunRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.WaitForJobRun(request)
		},
	}
}
