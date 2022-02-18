// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package function

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		functionRoot(),
		functionNamespace(),
		functionFunction(),
		functionCron(),
		functionRuntime(),
		functionLogs(),
		functionNamespaceList(),
		functionNamespaceGet(),
		functionNamespaceCreate(),
		functionNamespaceUpdate(),
		functionNamespaceDelete(),
		functionFunctionList(),
		functionFunctionGet(),
		functionFunctionCreate(),
		functionFunctionUpdate(),
		functionFunctionDelete(),
		functionFunctionDeploy(),
		functionRuntimeList(),
		functionFunctionGetUploadURL(),
		functionFunctionGetDownloadURL(),
		functionCronList(),
		functionCronGet(),
		functionCronDelete(),
		functionLogsList(),
	)
}
func functionRoot() *core.Command {
	return &core.Command{
		Short:     `Function as a Service API`,
		Long:      ``,
		Namespace: "function",
	}
}

func functionNamespace() *core.Command {
	return &core.Command{
		Short:     `Function namespace management commands`,
		Long:      `Function namespace management commands.`,
		Namespace: "function",
		Resource:  "namespace",
	}
}

func functionFunction() *core.Command {
	return &core.Command{
		Short:     `Function management commands`,
		Long:      `Function management commands.`,
		Namespace: "function",
		Resource:  "function",
	}
}

func functionCron() *core.Command {
	return &core.Command{
		Short:     `Cron management commands`,
		Long:      `Cron management commands.`,
		Namespace: "function",
		Resource:  "cron",
	}
}

func functionRuntime() *core.Command {
	return &core.Command{
		Short:     `Runtime management commands`,
		Long:      `Runtime management commands.`,
		Namespace: "function",
		Resource:  "runtime",
	}
}

func functionLogs() *core.Command {
	return &core.Command{
		Short:     `Logs management commands`,
		Long:      `Logs management commands.`,
		Namespace: "function",
		Resource:  "logs",
	}
}

func functionNamespaceList() *core.Command {
	return &core.Command{
		Short:     `List all your namespaces`,
		Long:      `List all your namespaces.`,
		Namespace: "function",
		Resource:  "namespace",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.ListNamespacesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
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
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.ListNamespacesRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			resp, err := api.ListNamespaces(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Namespaces, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "RegistryNamespaceID",
			},
			{
				FieldName: "RegistryEndpoint",
			},
			{
				FieldName: "EnvironmentVariables",
			},
			{
				FieldName: "ErrorMessage",
			},
			{
				FieldName: "Description",
			},
			{
				FieldName: "Region",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "OrganizationID",
			},
		}},
	}
}

func functionNamespaceGet() *core.Command {
	return &core.Command{
		Short:     `Get a namespace`,
		Long:      `Get the namespace associated with the given id.`,
		Namespace: "function",
		Resource:  "namespace",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.GetNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.GetNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.GetNamespace(request)

		},
	}
}

func functionNamespaceCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new namespace`,
		Long:      `Create a new namespace.`,
		Namespace: "function",
		Resource:  "namespace",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.CreateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("ns"),
			},
			{
				Name:       "environment-variables.value.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.CreateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.CreateNamespace(request)

		},
	}
}

func functionNamespaceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing namespace`,
		Long:      `Update the space associated with the given id.`,
		Namespace: "function",
		Resource:  "namespace",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.UpdateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "environment-variables.value.{key}",
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
				Name:       "secret-environment-variables.{index}.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.UpdateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.UpdateNamespace(request)

		},
	}
}

func functionNamespaceDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing namespace`,
		Long:      `Delete the namespace associated with the given id.`,
		Namespace: "function",
		Resource:  "namespace",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.DeleteNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.DeleteNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.DeleteNamespace(request)

		},
	}
}

func functionFunctionList() *core.Command {
	return &core.Command{
		Short:     `List all your functions`,
		Long:      `List all your functions.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.ListFunctionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "namespace-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
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
				Name:       "organization-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.ListFunctionsRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			resp, err := api.ListFunctions(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Functions, nil

		},
	}
}

func functionFunctionGet() *core.Command {
	return &core.Command{
		Short:     `Get a function`,
		Long:      `Get the function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.GetFunctionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.GetFunctionRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.GetFunction(request)

		},
	}
}

func functionFunctionCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new function`,
		Long:      `Create a new function.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.CreateFunctionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("fn"),
			},
			{
				Name:       "namespace-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.value.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-scale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-scale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "runtime",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_runtime", "golang", "python", "python3", "node8", "node10", "node14"},
			},
			{
				Name:       "memory-limit",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "handler",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "privacy",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_privacy", "public", "private"},
			},
			{
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-option",
				Short:      `Configure how HTTP and HTTPS requests are handled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.CreateFunctionRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.CreateFunction(request)

		},
	}
}

func functionFunctionUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing function`,
		Long:      `Update the function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.UpdateFunctionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "environment-variables.value.{key}",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-scale",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-scale",
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
				Name:       "timeout.seconds",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout.nanos",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "redeploy",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "handler",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "privacy",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_privacy", "public", "private"},
			},
			{
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.key",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{index}.value",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "http-option",
				Short:      `Configure how HTTP and HTTPS requests are handled`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.UpdateFunctionRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.UpdateFunction(request)

		},
	}
}

func functionFunctionDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a function`,
		Long:      `Delete the function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.DeleteFunctionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.DeleteFunctionRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.DeleteFunction(request)

		},
	}
}

func functionFunctionDeploy() *core.Command {
	return &core.Command{
		Short:     `Deploy a function`,
		Long:      `Deploy a function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "deploy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.DeployFunctionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.DeployFunctionRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.DeployFunction(request)

		},
	}
}

func functionRuntimeList() *core.Command {
	return &core.Command{
		Short:     `List function runtimes`,
		Long:      `List available function runtimes.`,
		Namespace: "function",
		Resource:  "runtime",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.ListFunctionRuntimesRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.ListFunctionRuntimesRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.ListFunctionRuntimes(request)

		},
	}
}

func functionFunctionGetUploadURL() *core.Command {
	return &core.Command{
		Short:     `Get an upload URL of a function`,
		Long:      `Get an upload URL of a function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "get-upload-url",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.GetFunctionUploadURLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "content-length",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.GetFunctionUploadURLRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.GetFunctionUploadURL(request)

		},
	}
}

func functionFunctionGetDownloadURL() *core.Command {
	return &core.Command{
		Short:     `Get a download URL of a function`,
		Long:      `Get a download URL for a function associated with the given id.`,
		Namespace: "function",
		Resource:  "function",
		Verb:      "get-download-url",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.GetFunctionDownloadURLRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.GetFunctionDownloadURLRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.GetFunctionDownloadURL(request)

		},
	}
}

func functionCronList() *core.Command {
	return &core.Command{
		Short:     `List all your crons`,
		Long:      `List all your crons.`,
		Namespace: "function",
		Resource:  "cron",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.ListCronsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "function-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.ListCronsRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			resp, err := api.ListCrons(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Crons, nil

		},
	}
}

func functionCronGet() *core.Command {
	return &core.Command{
		Short:     `Get a cron`,
		Long:      `Get the cron associated with the given id.`,
		Namespace: "function",
		Resource:  "cron",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.GetCronRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cron-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.GetCronRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.GetCron(request)

		},
	}
}

func functionCronDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing cron`,
		Long:      `Delete the cron associated with the given id.`,
		Namespace: "function",
		Resource:  "cron",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.DeleteCronRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "cron-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.DeleteCronRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			return api.DeleteCron(request)

		},
	}
}

func functionLogsList() *core.Command {
	return &core.Command{
		Short:     `List your application logs`,
		Long:      `List your application logs.`,
		Namespace: "function",
		Resource:  "logs",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(function.ListLogsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "function-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"timestamp_desc", "timestamp_asc"},
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*function.ListLogsRequest)

			client := core.ExtractClient(ctx)
			api := function.NewAPI(client)
			resp, err := api.ListLogs(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Logs, nil

		},
	}
}
