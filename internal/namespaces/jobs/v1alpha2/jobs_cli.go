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
		jobsTrigger(),
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
		jobsTriggerCreate(),
		jobsTriggerGet(),
		jobsTriggerList(),
		jobsTriggerUpdate(),
		jobsTriggerDelete(),
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

func jobsTrigger() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      ``,
		Namespace: "jobs",
		Resource:  "trigger",
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
		ArgsType: reflect.TypeFor[jobs.CreateJobDefinitionRequest](),
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
				Required:   true,
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
			{
				Name:       "retry-policy.max-retries",
				Short:      `Maximum number of retries upon a job failure.`,
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

			return api.CreateJobDefinition(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.GetJobDefinitionRequest](),
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

			return api.GetJobDefinition(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.ListJobDefinitionsRequest](),
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
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
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
		ArgsType: reflect.TypeFor[jobs.UpdateJobDefinitionRequest](),
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
			{
				Name:       "retry-policy.max-retries",
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

			return api.UpdateJobDefinition(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.DeleteJobDefinitionRequest](),
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
			e = api.DeleteJobDefinition(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.StartJobDefinitionRequest](),
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

			return api.StartJobDefinition(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.GetJobRunRequest](),
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

			return api.GetJobRun(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.ListJobRunsRequest](),
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
					"application_not_started",
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
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
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
		ArgsType: reflect.TypeFor[jobs.StopJobRunRequest](),
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

			return api.StopJobRun(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.CreateSecretsRequest](),
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

			return api.CreateSecrets(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.GetSecretRequest](),
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

			return api.GetSecret(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.ListSecretsRequest](),
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

			return api.ListSecrets(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.UpdateSecretRequest](),
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

			return api.UpdateSecret(request, scw.WithContext(ctx))
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
		ArgsType: reflect.TypeFor[jobs.DeleteSecretRequest](),
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
			e = api.DeleteSecret(request, scw.WithContext(ctx))
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

func jobsTriggerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a trigger`,
		Long:      `Create a trigger.`,
		Namespace: "jobs",
		Resource:  "trigger",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[jobs.CreateTriggerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the trigger`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.schedule",
				Short:      `CRON schedule in UNIX format`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.timezone",
				Short:      `Time zone for the CRON schedule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.startup-command.{index}",
				Short:      `Startup command that will be used by the triggered job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.args.{index}",
				Short:      `Arguments passed to the startup command used by the triggered job`,
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
			request := args.(*jobs.CreateTriggerRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.CreateTrigger(request, scw.WithContext(ctx))
		},
	}
}

func jobsTriggerGet() *core.Command {
	return &core.Command{
		Short:     `Get a trigger`,
		Long:      `Get a trigger.`,
		Namespace: "jobs",
		Resource:  "trigger",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[jobs.GetTriggerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
				Short:      `UUID of the trigger`,
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
			request := args.(*jobs.GetTriggerRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.GetTrigger(request, scw.WithContext(ctx))
		},
	}
}

func jobsTriggerList() *core.Command {
	return &core.Command{
		Short:     `List triggers of a job definition`,
		Long:      `List triggers of a job definition.`,
		Namespace: "jobs",
		Resource:  "trigger",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[jobs.ListTriggersRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "job-definition-id",
				Short:      `UUID of the job definition`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sorting order of triggers`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.RegionNlAms,
				scw.RegionPlWaw,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*jobs.ListTriggersRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListTriggers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Triggers, nil
		},
	}
}

func jobsTriggerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a trigger`,
		Long:      `Update a trigger.`,
		Namespace: "jobs",
		Resource:  "trigger",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[jobs.UpdateTriggerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
				Short:      `UUID of the trigger`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the trigger`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.schedule",
				Short:      `CRON schedule in UNIX format`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.timezone",
				Short:      `Time zone for the CRON schedule`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.startup-command.{index}",
				Short:      `Startup command that will be used by the triggered job`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.args.{index}",
				Short:      `Arguments passed to the startup command used by the triggered job`,
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
			request := args.(*jobs.UpdateTriggerRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)

			return api.UpdateTrigger(request, scw.WithContext(ctx))
		},
	}
}

func jobsTriggerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a trigger`,
		Long:      `Delete a trigger.`,
		Namespace: "jobs",
		Resource:  "trigger",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[jobs.DeleteTriggerRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
				Short:      `UUID of the trigger`,
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
			request := args.(*jobs.DeleteTriggerRequest)

			client := core.ExtractClient(ctx)
			api := jobs.NewAPI(client)
			e = api.DeleteTrigger(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "trigger",
				Verb:     "delete",
			}, nil
		},
	}
}
