// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package jobs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		jobsRoot(),
		jobsRun(),
		jobsDefinition(),
		jobsDefinitionCreate(),
		jobsDefinitionGet(),
		jobsDefinitionList(),
		jobsDefinitionUpdate(),
		jobsDefinitionDelete(),
		jobsDefinitionStart(),
		jobsRunGet(),
		jobsRunStop(),
		jobsRunList(),
	)
}
func jobsRoot() *core.Command {
	return &core.Command{
		Short:     `Serverless Jobs API`,
		Long:      `Serverless Jobs API.`,
		Namespace: "jobs",
	}
}

func jobsRun() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "jobs",
		Resource:  "run",
	}
}

func jobsDefinition() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "jobs",
		Resource:  "definition",
	}
}

func jobsDefinitionCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new job definition in a specified Project`,
		Long:      `Create a new job definition in a specified Project.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.CreateJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the job definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("job"),
			},
			{
				Name:       "cpu-limit",
				Short:      `CPU limit of the job`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit",
				Short:      `Memory limit of the job (in MiB)`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "local-storage-capacity",
				Short:      `Local storage capacity of the job (in MiB)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image-uri",
				Short:      `Image to use for the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command",
				Short:      `Startup command. If empty or not defined, the image's default command is used.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "job-timeout",
				Short:      `Timeout of the job in seconds`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.schedule",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.timezone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.CreateJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.CreateJobDefinition(request)

		},
	}
}

func jobsDefinitionGet() *core.Command {
	return &core.Command{
		Short:     `Get a job definition by its unique identifier`,
		Long:      `Get a job definition by its unique identifier.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.GetJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.GetJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.GetJobDefinition(request)

		},
	}
}

func jobsDefinitionList() *core.Command {
	return &core.Command{
		Short:     `List all your job definitions with filters`,
		Long:      `List all your job definitions with filters.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.ListJobDefinitionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.ListJobDefinitionsRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListJobDefinitions(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.JobDefinitions, nil

		},
	}
}

func jobsDefinitionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing job definition associated with the specified unique identifier`,
		Long:      `Update an existing job definition associated with the specified unique identifier.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.UpdateJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the job definition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-limit",
				Short:      `CPU limit of the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit",
				Short:      `Memory limit of the job (in MiB)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "local-storage-capacity",
				Short:      `Local storage capacity of the job (in MiB)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image-uri",
				Short:      `Image to use for the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command",
				Short:      `Startup command`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "job-timeout",
				Short:      `Timeout of the job in seconds`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.schedule",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.timezone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.UpdateJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.UpdateJobDefinition(request)

		},
	}
}

func jobsDefinitionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an exsisting job definition by its unique identifier`,
		Long:      `Delete an exsisting job definition by its unique identifier.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.DeleteJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.DeleteJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			e = api.DeleteJobDefinition(request)
			if e != nil {
				return nil, e
			}
			return &core.SuccessResult{
				Resource: "definition",
				Verb:     "delete",
			}, nil
		},
	}
}

func jobsDefinitionStart() *core.Command {
	return &core.Command{
		Short:     `Run an existing job definition by its unique identifier. This will create a new job run`,
		Long:      `Run an existing job definition by its unique identifier. This will create a new job run.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.StartJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition to start`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "command",
				Short:      `Contextual startup command for this specific job run`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Contextual environment variables for this specific job run`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "replicas",
				Short:      `Number of jobs to run`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.StartJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.StartJobDefinition(request)

		},
	}
}

func jobsRunGet() *core.Command {
	return &core.Command{
		Short:     `Get a job run by its unique identifier`,
		Long:      `Get a job run by its unique identifier.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.GetJobRunRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-run-id",
				Short:      `UUID of the job run to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.GetJobRunRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.GetJobRun(request)

		},
	}
}

func jobsRunStop() *core.Command {
	return &core.Command{
		Short:     `Stop a job run by its unique identifier`,
		Long:      `Stop a job run by its unique identifier.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.StopJobRunRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-run-id",
				Short:      `UUID of the job run to stop`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.StopJobRunRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			return api.StopJobRun(request)

		},
	}
}

func jobsRunList() *core.Command {
	return &core.Command{
		Short:     `List all job runs with filters`,
		Long:      `List all job runs with filters.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.ListJobRunsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "job-definition-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*jobs.ListJobRunsRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListJobRuns(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.JobRuns, nil

		},
	}
}
