// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package dataorchestrator

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/dataorchestrator/v1alpha2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		workflowsRoot(),
		workflowsDefinition(),
		workflowsVersion(),
		workflowsRun(),
		workflowsDefinitionCreate(),
		workflowsDefinitionGet(),
		workflowsDefinitionList(),
		workflowsDefinitionUpdate(),
		workflowsDefinitionDelete(),
		workflowsVersionUpdate(),
		workflowsDefinitionStart(),
		workflowsRunGet(),
		workflowsRunList(),
		workflowsRunTerminate(),
		workflowsRunPause(),
		workflowsRunResume(),
		workflowsRunDelete(),
	)
}

func workflowsRoot() *core.Command {
	return &core.Command{
		Short:     `This API is a boilerplate for products`,
		Long:      `This API is a boilerplate for products.`,
		Namespace: "workflows",
	}
}

func workflowsDefinition() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "workflows",
		Resource:  "definition",
	}
}

func workflowsVersion() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "workflows",
		Resource:  "version",
	}
}

func workflowsRun() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "workflows",
		Resource:  "run",
	}
}

func workflowsDefinitionCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new workflow definition in a specified Project`,
		Long:      `Create a new workflow definition in a specified Project.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.CreateWorkflowDefinitionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the workflow definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("workflow"),
			},
			{
				Name:       "description",
				Short:      `Description of the workflow`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the workflow`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "workflow-timeout",
				Short:      `Timeout of the workflow in seconds`,
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
				Name:       "cron-schedule.timezone-identifier",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.overlap",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_schedule_overlap_policy",
					"forbid",
					"buffer_one",
					"buffer_all",
					"allow_all",
				},
			},
			{
				Name:       "workflow-version.name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "workflow-version.description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "workflow-version.workflow-dsl",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.CreateWorkflowDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.CreateWorkflowDefinition(request, scw.WithContext(ctx))
		},
	}
}

func workflowsDefinitionGet() *core.Command {
	return &core.Command{
		Short:     `Get a workflow definition by its unique identifier`,
		Long:      `Get a workflow definition by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.GetWorkflowDefinitionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-definition-id",
				Short:      `UUID of the workflow definition to get`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.GetWorkflowDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.GetWorkflowDefinition(request, scw.WithContext(ctx))
		},
	}
}

func workflowsDefinitionList() *core.Command {
	return &core.Command{
		Short:     `List all your workflow definitions with filters`,
		Long:      `List all your workflow definitions with filters.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.ListWorkflowDefinitionsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"name_asc",
					"name_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "statuses.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_run_status",
					"initialized",
					"validated",
					"queued",
					"running",
					"succeeded",
					"failed",
					"interrupting",
					"interrupted",
					"pausing",
					"paused",
				},
			},
			{
				Name:       "created-after",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "created-before",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "trigger-types.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_trigger_type",
					"manual",
					"scheduled",
				},
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.ListWorkflowDefinitionsRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListWorkflowDefinitions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.WorkflowDefinitions, nil
		},
	}
}

func workflowsDefinitionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing workflow definition by its unique identifier`,
		Long:      `Update an existing workflow definition by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.UpdateWorkflowDefinitionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-definition-id",
				Short:      `UUID of the workflow definition to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the workflow definition`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the workflow`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the workflow`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "workflow-timeout",
				Short:      `Timeout of the workflow in seconds`,
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
				Name:       "cron-schedule.timezone-identifier",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-schedule.overlap",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_schedule_overlap_policy",
					"forbid",
					"buffer_one",
					"buffer_all",
					"allow_all",
				},
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.UpdateWorkflowDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.UpdateWorkflowDefinition(request, scw.WithContext(ctx))
		},
	}
}

func workflowsDefinitionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing workflow definition by its unique identifier`,
		Long:      `Delete an existing workflow definition by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.DeleteWorkflowDefinitionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-definition-id",
				Short:      `UUID of the workflow definition to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.DeleteWorkflowDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)
			e = api.DeleteWorkflowDefinition(request, scw.WithContext(ctx))
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

func workflowsVersionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing workflow version by its unique identifier`,
		Long:      `Update an existing workflow version by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "version",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.UpdateWorkflowVersionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-version-id",
				Short:      `UUID of the workflow version to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the workflow version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the workflow version`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "workflow-dsl",
				Short:      `DSL of workflow in YAML format based from https://open-workflow-specification.org/schemas/1.0.3/workflow.yaml`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.UpdateWorkflowVersionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.UpdateWorkflowVersion(request, scw.WithContext(ctx))
		},
	}
}

func workflowsDefinitionStart() *core.Command {
	return &core.Command{
		Short:     `Run an existing workflow definition using its unique identifier and create a new workflow run`,
		Long:      `Run an existing workflow definition using its unique identifier and create a new workflow run.`,
		Namespace: "workflows",
		Resource:  "definition",
		Verb:      "start",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.StartWorkflowDefinitionRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-definition-id",
				Short:      `UUID of the workflow definition to start`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Contextual environment variables for this specific workflow run`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.StartWorkflowDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.StartWorkflowDefinition(request, scw.WithContext(ctx))
		},
	}
}

func workflowsRunGet() *core.Command {
	return &core.Command{
		Short:     `Get a workflow run by its unique identifier`,
		Long:      `Get a workflow run by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.GetWorkflowRunRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-run-id",
				Short:      `UUID of the workflow run to get`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.GetWorkflowRunRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.GetWorkflowRun(request, scw.WithContext(ctx))
		},
	}
}

func workflowsRunList() *core.Command {
	return &core.Command{
		Short:     `List all your workflow runs with filters`,
		Long:      `List all your workflow runs with filters.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.ListWorkflowRunsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"terminated_at_asc",
					"terminated_at_desc",
					"duration_asc",
					"duration_desc",
					"workflow_definition_name_asc",
					"workflow_definition_name_desc",
				},
			},
			{
				Name:       "workflow-definition-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "statuses.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_run_status",
					"initialized",
					"validated",
					"queued",
					"running",
					"succeeded",
					"failed",
					"interrupting",
					"interrupted",
					"pausing",
					"paused",
				},
			},
			{
				Name:       "created-after",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "created-before",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "trigger-types.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_trigger_type",
					"manual",
					"scheduled",
				},
			},
			{
				Name:       "workflow-definition-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.ListWorkflowRunsRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListWorkflowRuns(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.WorkflowRuns, nil
		},
	}
}

func workflowsRunTerminate() *core.Command {
	return &core.Command{
		Short:     `Terminate a workflow run using its unique identifier`,
		Long:      `Terminate a workflow run using its unique identifier.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "terminate",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.TerminateWorkflowRunRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-run-id",
				Short:      `UUID of the workflow run to terminate`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.TerminateWorkflowRunRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.TerminateWorkflowRun(request, scw.WithContext(ctx))
		},
	}
}

func workflowsRunPause() *core.Command {
	return &core.Command{
		Short:     `Pause a workflow run using its unique identifier`,
		Long:      `Pause a workflow run using its unique identifier.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "pause",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.PauseWorkflowRunRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-run-id",
				Short:      `UUID of the workflow run to pause`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.PauseWorkflowRunRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.PauseWorkflowRun(request, scw.WithContext(ctx))
		},
	}
}

func workflowsRunResume() *core.Command {
	return &core.Command{
		Short:     `Resume a workflow run using its unique identifier`,
		Long:      `Resume a workflow run using its unique identifier.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "resume",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.ResumeWorkflowRunRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-run-id",
				Short:      `UUID of the workflow run to resume`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.ResumeWorkflowRunRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)

			return api.ResumeWorkflowRun(request, scw.WithContext(ctx))
		},
	}
}

func workflowsRunDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing workflow run by its unique identifier`,
		Long:      `Delete an existing workflow run by its unique identifier.`,
		Namespace: "workflows",
		Resource:  "run",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataorchestrator.DeleteWorkflowRunRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "workflow-run-id",
				Short:      `UUID of the workflow run to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataorchestrator.DeleteWorkflowRunRequest)

			client := core.ExtractClient(ctx)
			api := dataorchestrator.NewAPI(client)
			e = api.DeleteWorkflowRun(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "run",
				Verb:     "delete",
			}, nil
		},
	}
}
