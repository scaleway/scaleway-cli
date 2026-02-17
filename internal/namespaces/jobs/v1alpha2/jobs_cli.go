// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package jobs

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	jobs "github.com/scaleway/scaleway-sdk-go/api/jobs/v1alpha2"
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
		jobsSecret(),
		jobsDefinitionCreate(),
		jobsDefinitionGet(),
		jobsDefinitionList(),
		jobsDefinitionUpdate(),
		jobsDefinitionDelete(),
		jobsDefinitionStart(),
		jobsRunGet(),
		jobsRunList(),
		jobsRunStop(),
		jobsSecretCreate(),
		jobsSecretGet(),
		jobsSecretList(),
		jobsSecretUpdate(),
		jobsSecretDelete(),
	)
}

func jobsRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Serverless Jobs`,
		Long:      `This API allows you to manage your Serverless Jobs.`,
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

func jobsSecret() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "jobs",
		Resource:  "secret",
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
				Short:      `CPU limit of the job (in mvCPU)`,
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
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command",
				Short:      `Startup command. If empty or not defined, the image's default command is used.`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "startup-command.{index}",
				Short:      `Job startup command. Overrides the default defined in the job image.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "args.{index}",
				Short:      `Job arguments. Overrides the default arguments defined in the job image.`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
				Short:      `CPU limit of the job (in mvCPU)`,
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
				Short:      `Startup command. If empty or not defined, the image's default command is used.`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "startup-command.{index}",
				Short:      `Job startup command. Overrides the default defined in the job image.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "args.{index}",
				Short:      `Job arguments. Overrides the default arguments defined in the job image.`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.UpdateJobDefinitionRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.UpdateJobDefinition(request)
		},
	}
}

func jobsDefinitionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing job definition by its unique identifier`,
		Long:      `Delete an existing job definition by its unique identifier.`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
		Short:     `Run an existing job definition using its unique identifier and create a new job run`,
		Long:      `Run an existing job definition using its unique identifier and create a new job run.`,
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
				Short:      `Contextual startup command for this specific job run.`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			{
				Name:       "startup-command.{index}",
				Short:      `Contextual startup command for this specific job run.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "args.{index}",
				Short:      `Contextual arguments for this specific job run.`,
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.GetJobRunRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.GetJobRun(request)
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
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
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
			{
				Name:       "state",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_state",
					"initialized",
					"validated",
					"queued",
					"running",
					"succeeded",
					"failed",
					"interrupting",
					"interrupted",
					"retrying",
				},
			},
			{
				Name:       "states.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_state",
					"initialized",
					"validated",
					"queued",
					"running",
					"succeeded",
					"failed",
					"interrupting",
					"interrupted",
					"retrying",
				},
			},
			{
				Name:       "reasons.{index}",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_reason",
					"invalid_request",
					"timeout",
					"cancellation",
					"technical_error",
					"image_not_found",
					"invalid_image",
					"memory_usage_exceeded",
					"storage_usage_exceeded",
					"exited_with_error",
					"secret_disabled",
					"secret_not_found",
					"quota_exceeded",
				},
			},
			{
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
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

func jobsRunStop() *core.Command {
	return &core.Command{
		Short:     `Stop a job run using its unique identifier`,
		Long:      `Stop a job run using its unique identifier.`,
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
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.StopJobRunRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.StopJobRun(request)
		},
	}
}

func jobsSecretCreate() *core.Command {
	return &core.Command{
		Short:     `Create a secret reference within a job definition`,
		Long:      `Create a secret reference within a job definition.`,
		Namespace: "jobs",
		Resource:  "secret",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.CreateSecretsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secrets.{index}.secret-manager-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secrets.{index}.secret-manager-version",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secrets.{index}.path",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secrets.{index}.env-var-name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.CreateSecretsRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.CreateSecrets(request)
		},
	}
}

func jobsSecretGet() *core.Command {
	return &core.Command{
		Short:     `Get a secret references within a job definition`,
		Long:      `Get a secret references within a job definition.`,
		Namespace: "jobs",
		Resource:  "secret",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.GetSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `UUID of the secret reference within the job`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.GetSecretRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.GetSecret(request)
		},
	}
}

func jobsSecretList() *core.Command {
	return &core.Command{
		Short:     `List secrets references within a job definition`,
		Long:      `List secrets references within a job definition.`,
		Namespace: "jobs",
		Resource:  "secret",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.ListSecretsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.ListSecretsRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.ListSecrets(request)
		},
	}
}

func jobsSecretUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a secret reference within a job definition`,
		Long:      `Update a secret reference within a job definition.`,
		Namespace: "jobs",
		Resource:  "secret",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.UpdateSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `UUID of the secret reference within the job`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-manager-version",
				Short:      `Version of the secret in Secret Manager`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "path",
				Short:      `Path of the secret to mount inside the job (either ` + "`" + `path` + "`" + ` or ` + "`" + `env_var_name` + "`" + ` must be set)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "env-var-name",
				Short:      `Environment variable name used to expose the secret inside the job (either ` + "`" + `path` + "`" + ` or ` + "`" + `env_var_name` + "`" + ` must be set)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.UpdateSecretRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.UpdateSecret(request)
		},
	}
}

func jobsSecretDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a secret reference within a job definition`,
		Long:      `Delete a secret reference within a job definition.`,
		Namespace: "jobs",
		Resource:  "secret",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(jobs.DeleteSecretRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "secret-id",
				Short:      `UUID of the secret reference within the job`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.DeleteSecretRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			e = api.DeleteSecret(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "secret",
				Verb:     "delete",
			}, nil
		},
	}
}
