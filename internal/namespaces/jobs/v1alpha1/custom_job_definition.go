package jobs

import (
	"context"
	"errors"
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func definitionStartBuilder(c *core.Command) *core.Command {
	c.WaitUsage = "Wait until the job reach a stable state, use job definition timeout"
	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := jobs.NewAPI(core.ExtractClient(ctx))
		args := argsI.(*jobs.StartJobDefinitionRequest)
		resp := respI.(*jobs.StartJobDefinitionResponse)

		jobDefinition, err := api.GetJobDefinition(&jobs.GetJobDefinitionRequest{
			Region:          args.Region,
			JobDefinitionID: args.JobDefinitionID,
		}, scw.WithContext(ctx))
		if err != nil {
			return nil, fmt.Errorf("failed to fetch job definition for timeout: %w", err)
		}

		if len(resp.JobRuns) == 0 {
			return nil, errors.New("no job run found")
		}

		return api.WaitForJobRun(&jobs.WaitForJobRunRequest{
			Region:        args.Region,
			JobRunID:      resp.JobRuns[0].ID,
			Timeout:       jobDefinition.JobTimeout.ToTimeDuration(),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
