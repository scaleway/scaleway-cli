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
		Short:     `Create jobs resources`,
		Long:      `Create jobs resources.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.CreateJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("job"),
			},
			{
				Name:       "cpu-limit",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image-uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "environment-variables.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "job-timeout",
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
		Short:     `Get jobs resources`,
		Long:      `Get jobs resources.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.GetJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
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
		Short:     `List jobs resources`,
		Long:      `List jobs resources.`,
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
		Short:     `Update jobs resources`,
		Long:      `Update jobs resources.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.UpdateJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-limit",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image-uri",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "job-timeout",
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
		Short:     `Delete jobs resources`,
		Long:      `Delete jobs resources.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.DeleteJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
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
		Short:     `Start jobs resources`,
		Long:      `Start jobs resources.`,
		Namespace: "jobs",
		Resource:  "definition",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.StartJobDefinitionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
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
		Short:     `Get jobs resources`,
		Long:      `Get jobs resources.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.GetJobRunRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-run-id",
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
		Short:     `Stop jobs resources`,
		Long:      `Stop jobs resources.`,
		Namespace: "jobs",
		Resource:  "run",
		Verb:      "stop",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.StopJobRunRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-run-id",
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
		Short:     `List jobs resources`,
		Long:      `List jobs resources.`,
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
