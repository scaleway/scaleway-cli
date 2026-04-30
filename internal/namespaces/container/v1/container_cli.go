// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package container

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		containerRoot(),
		containerNamespace(),
		containerContainer(),
		containerDomain(),
		containerTrigger(),
		containerNamespaceCreate(),
		containerNamespaceGet(),
		containerNamespaceList(),
		containerNamespaceUpdate(),
		containerNamespaceDelete(),
		containerContainerCreate(),
		containerContainerGet(),
		containerContainerList(),
		containerContainerUpdate(),
		containerContainerDelete(),
		containerDomainCreate(),
		containerDomainGet(),
		containerDomainList(),
		containerDomainUpdate(),
		containerDomainDelete(),
		containerContainerRedeploy(),
		containerTriggerCreate(),
		containerTriggerGet(),
		containerTriggerList(),
		containerTriggerUpdate(),
		containerTriggerDelete(),
	)
}

func containerRoot() *core.Command {
	return &core.Command{
		Short:     `Easily run containers on the cloud with a single command`,
		Long:      `Easily run containers on the cloud with a single command.`,
		Namespace: "container",
	}
}

func containerNamespace() *core.Command {
	return &core.Command{
		Short:     `Namespace management commands`,
		Long:      `Namespace management commands.`,
		Namespace: "container",
		Resource:  "namespace",
	}
}

func containerContainer() *core.Command {
	return &core.Command{
		Short:     `Container management commands`,
		Long:      `Container management commands.`,
		Namespace: "container",
		Resource:  "container",
	}
}

func containerDomain() *core.Command {
	return &core.Command{
		Short:     `Domain management commands`,
		Long:      `Domain management commands.`,
		Namespace: "container",
		Resource:  "domain",
	}
}

func containerTrigger() *core.Command {
	return &core.Command{
		Short:     `Trigger management commands`,
		Long:      `Trigger management commands.`,
		Namespace: "container",
		Resource:  "trigger",
	}
}

func containerNamespaceCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new namespace.`,
		Long:      `Namespace name must be unique inside a project.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Namespace name.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Namespace description.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Namespace environment variables.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{key}",
				Short:      `Namespace secret environment variables.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `A list of arbitrary tags associated with the namespace.`,
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
			request := args.(*container.CreateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.CreateNamespace(request)
		},
	}
}

func containerNamespaceGet() *core.Command {
	return &core.Command{
		Short:     `Get the namespace associated with the specified ID.`,
		Long:      `Get the namespace associated with the specified ID.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
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
			request := args.(*container.GetNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.GetNamespace(request)
		},
	}
}

func containerNamespaceList() *core.Command {
	return &core.Command{
		Short: `List all namespaces the caller can access (read permission).`,
		Long: `By default, the namespaces listed are ordered by creation date in ascending order. This can be modified via the ` + "`" + `order_by` + "`" + ` field.

Additional parameters can be set in the query to filter, such as ` + "`" + `organization_id` + "`" + `, ` + "`" + `project_id` + "`" + `, and ` + "`" + `name` + "`" + `.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListNamespacesRequest{}),
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
				},
			},
			{
				Name:       "project-id",
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
			request := args.(*container.ListNamespacesRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNamespaces(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Namespaces, nil
		},
	}
}

func containerNamespaceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update the namespace associated with the specified ID.`,
		Long:      `Only fields present in the request are updated; others are left untouched.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace to update.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Namespace description.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Namespace environment variables.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{key}",
				Short:      `Namespace secret environment variables.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `A list of arbitrary tags associated with the namespace.`,
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
			request := args.(*container.UpdateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.UpdateNamespace(request)
		},
	}
}

func containerNamespaceDelete() *core.Command {
	return &core.Command{
		Short: `Delete the namespace associated with the specified ID.`,
		Long: `It also deletes in cascade any resource inside the namespace.

This action **cannot** be undone.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace to delete.`,
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
			request := args.(*container.DeleteNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.DeleteNamespace(request)
		},
	}
}

func containerContainerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new container in a namespace.`,
		Long:      `Name must be unique inside the given namespace.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `Unique ID of the namespace the container belongs to.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Container name.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{key}",
				Short:      `Secret environment variables of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-scale",
				Short:      `Minimum number of instances to scale the container to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-scale",
				Short:      `Maximum number of instances to scale the container to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit-bytes",
				Short:      `Memory limit of the container in bytes.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mvcpu-limit",
				Short:      `CPU limit of the container in mvCPU.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout",
				Short:      `Processing time limit for the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "privacy",
				Short:      `Privacy policy of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("public"),
				EnumValues: []string{
					"unknown_privacy",
					"public",
					"private",
				},
			},
			{
				Name:       "description",
				Short:      `Container description.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image",
				Short:      `Image reference (e.g. "rg.fr-par.scw.cloud/my-registry-namespace/image:tag").`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Protocol the container uses.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("http1"),
				EnumValues: []string{
					"unknown_protocol",
					"http1",
					"h2c",
				},
			},
			{
				Name:       "port",
				Short:      `Port the container listens on.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("8080"),
			},
			{
				Name:       "https-connections-only",
				Short:      `If true, it will allow only HTTPS connections to access your container to prevent it from being triggered by insecure connections (HTTP).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "sandbox",
				Short:      `Execution environment of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.DefaultValueSetter("v2"),
				EnumValues: []string{
					"unknown_sandbox",
					"v1",
					"v2",
				},
			},
			{
				Name:       "local-storage-limit-bytes",
				Short:      `Local storage limit of the container (in bytes).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.concurrent-requests-threshold",
				Short:      `Scale depending on the number of concurrent requests being processed per container instance. The threshold value is the number of concurrent requests above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.cpu-usage-threshold",
				Short:      `Scale depending on the CPU usage of a container instance. The threshold value is the percentage of CPU usage above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.memory-usage-threshold",
				Short:      `Scale depending on the memory usage of a container instance. The threshold value is the percentage of memory usage above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.failure-threshold",
				Short:      `Number of consecutive failures before considering the container as unhealthy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.interval",
				Short:      `Time interval between checks.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.timeout",
				Short:      `Duration before the check times out.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.http.path",
				Short:      `HTTP path to perform the check on.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.failure-threshold",
				Short:      `Number of consecutive failures before considering the container as unhealthy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.interval",
				Short:      `Time interval between checks.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.timeout",
				Short:      `Duration before the check times out.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.http.path",
				Short:      `HTTP path to perform the check on.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the Serverless Container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network the container is connected to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command.{index}",
				Short:      `Container command`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "args.{index}",
				Short:      `Container arguments`,
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
			request := args.(*container.CreateContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.CreateContainer(request)
		},
	}
}

func containerContainerGet() *core.Command {
	return &core.Command{
		Short:     `Get the container associated with the specified ID.`,
		Long:      `Get the container associated with the specified ID.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
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
			request := args.(*container.GetContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.GetContainer(request)
		},
	}
}

func containerContainerList() *core.Command {
	return &core.Command{
		Short: `List all containers the caller can access (read permission).`,
		Long: `By default, the containers listed are ordered by creation date in ascending order. This can be modified via the ` + "`" + `order_by` + "`" + ` field.

Additional parameters can be set in the query to filter, such as ` + "`" + `organization_id` + "`" + `, ` + "`" + `project_id` + "`" + `, and ` + "`" + `name` + "`" + `.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListContainersRequest{}),
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
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
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
			request := args.(*container.ListContainersRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListContainers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Containers, nil
		},
	}
}

func containerContainerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update the container associated with the specified ID.`,
		Long:      `Only fields present in the request are updated; others are left untouched.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `UUID of the container to update.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "environment-variables.{key}",
				Short:      `Environment variables of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "secret-environment-variables.{key}",
				Short:      `Secret environment variables of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "min-scale",
				Short:      `Minimum number of instances to scale the container to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-scale",
				Short:      `Maximum number of instances to scale the container to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "memory-limit-bytes",
				Short:      `Memory limit of the container in bytes.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "mvcpu-limit",
				Short:      `CPU limit of the container in mvCPU.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "timeout",
				Short:      `Processing time limit for the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "privacy",
				Short:      `Privacy policy of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_privacy",
					"public",
					"private",
				},
			},
			{
				Name:       "description",
				Short:      `Container description.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "image",
				Short:      `Image reference (e.g. "rg.fr-par.scw.cloud/my-registry-namespace/image:tag").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Short:      `Protocol the container uses.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_protocol",
					"http1",
					"h2c",
				},
			},
			{
				Name:       "port",
				Short:      `Port the container listens on.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "https-connection-only",
				Short:      `If true, it will allow only HTTPS connections to access your container to prevent it from being triggered by insecure connections (HTTP).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sandbox",
				Short:      `Execution environment of the container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_sandbox",
					"v1",
					"v2",
				},
			},
			{
				Name:       "local-storage-limit-bytes",
				Short:      `Local storage limit of the container (in bytes).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.concurrent-requests-threshold",
				Short:      `Scale depending on the number of concurrent requests being processed per container instance. The threshold value is the number of concurrent requests above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.cpu-usage-threshold",
				Short:      `Scale depending on the CPU usage of a container instance. The threshold value is the percentage of CPU usage above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "scaling-option.memory-usage-threshold",
				Short:      `Scale depending on the memory usage of a container instance. The threshold value is the percentage of memory usage above which the container will be scaled up.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.failure-threshold",
				Short:      `Number of consecutive failures before considering the container as unhealthy`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.interval",
				Short:      `Time interval between checks.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.timeout",
				Short:      `Duration before the check times out.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "liveness-probe.http.path",
				Short:      `HTTP path to perform the check on.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.failure-threshold",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.interval",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.timeout",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "startup-probe.http.path",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the Serverless Container.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `ID of the Private Network the container is connected to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "command.{index}",
				Short:      `Container command`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "args.{index}",
				Short:      `Container arguments`,
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
			request := args.(*container.UpdateContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.UpdateContainer(request)
		},
	}
}

func containerContainerDelete() *core.Command {
	return &core.Command{
		Short: `Delete the container associated with the specified ID.`,
		Long: `It also deletes in cascade any resource linked to the container (crons, tokens, etc.).

This action **cannot** be undone.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `UUID of the container to delete.`,
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
			request := args.(*container.DeleteContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.DeleteContainer(request)
		},
	}
}

func containerDomainCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new custom domain for the container with the specified ID.`,
		Long:      `Create a new custom domain for the container with the specified ID.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `Unique ID of the container the domain will be assigned to.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "hostname",
				Short:      `Domain assigned to the container.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `A list of arbitrary tags associated with the domain.`,
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
			request := args.(*container.CreateDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.CreateDomain(request)
		},
	}
}

func containerDomainGet() *core.Command {
	return &core.Command{
		Short:     `Get the custom domain associated with the specified ID.`,
		Long:      `Get the custom domain associated with the specified ID.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
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
			request := args.(*container.GetDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.GetDomain(request)
		},
	}
}

func containerDomainList() *core.Command {
	return &core.Command{
		Short: `List all custom domains the caller can access (read permission).`,
		Long: `By default, the custom domains listed are ordered by creation date in ascending order. This can be modified via the ` + "`" + `order_by` + "`" + ` field.
    
Additional parameters can be set in the query to filter the output, such as ` + "`" + `organization_id` + "`" + `, ` + "`" + `project_id` + "`" + `, ` + "`" + `namespace_id` + "`" + `, or ` + "`" + `container_id` + "`" + `.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListDomainsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_asc",
					"created_at_desc",
					"hostname_asc",
					"hostname_desc",
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "namespace-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "container-id",
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
			request := args.(*container.ListDomainsRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDomains(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Domains, nil
		},
	}
}

func containerDomainUpdate() *core.Command {
	return &core.Command{
		Short:     `Update the domain associated with the specified ID.`,
		Long:      `Only fields present in the request are updated; others are left untouched.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `UUID of the domain to update.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "tags.{index}",
				Short:      `A list of arbitrary tags associated with the domain.`,
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
			request := args.(*container.UpdateDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.UpdateDomain(request)
		},
	}
}

func containerDomainDelete() *core.Command {
	return &core.Command{
		Short:     `Delete the custom domain associated with the specified ID.`,
		Long:      `Delete the custom domain associated with the specified ID.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Short:      `UUID of the domain to delete.`,
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
			request := args.(*container.DeleteDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.DeleteDomain(request)
		},
	}
}

func containerContainerRedeploy() *core.Command {
	return &core.Command{
		Short: `Redeploy a container`,
		Long: `Performs a rollout of the container by creating new instances with the latest image version and terminating the old instances.
When using mutable registry image references (e.g. ` + "`" + `my-registry-namespace/image:tag` + "`" + `), this endpoint can be used to force the container to use
the most recent image version available in the registry.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "redeploy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.RedeployContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `ID of the container to redeploy.`,
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
			request := args.(*container.RedeployContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.RedeployContainer(request)
		},
	}
}

func containerTriggerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new trigger for the container with the specified ID.`,
		Long:      `Create a new trigger for the container with the specified ID.`,
		Namespace: "container",
		Resource:  "trigger",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateTriggerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Short:      `ID of the container to trigger.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the trigger.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the trigger.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the trigger.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-config.http-path",
				Short:      `The HTTP path to send the request to (e.g., "/my-webhook-endpoint").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-config.http-method",
				Short:      `The HTTP method to use when sending the request (e.g., get, post, put, patch, delete). Must be specified as lowercase.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_http_method",
					"get",
					"post",
					"put",
					"patch",
					"delete",
				},
			},
			{
				Name:       "cron-config.schedule",
				Short:      `UNIX cron schedule to run job (e.g., "* * * * *").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.timezone",
				Short:      `Timezone for the cron schedule, in tz database format (e.g., "Europe/Paris").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.body",
				Short:      `Body to send to the container when the trigger is invoked.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.headers.{key}",
				Short:      `Additional headers to send to the container when the trigger is invoked.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.region",
				Short:      `The region where the SQS queue is hosted (e.g., "fr-par", "nl-ams").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.endpoint",
				Short:      `Endpoint URL to use to access SQS (e.g., "https://sqs.mnq.fr-par.scaleway.com").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.access-key-id",
				Short:      `The access key for accessing the SQS queue.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.secret-access-key",
				Short:      `The secret key for accessing the SQS queue.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.queue-url",
				Short:      `The URL of the SQS queue to monitor for messages.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.server-urls.{index}",
				Short:      `The URLs of the NATS server (e.g., "nats://nats.mnq.fr-par.scaleway.com:4222").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.subject",
				Short:      `NATS subject to subscribe to (e.g., "my-subject").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.credentials-file-content",
				Short:      `The content of the NATS credentials file.`,
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
			request := args.(*container.CreateTriggerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.CreateTrigger(request)
		},
	}
}

func containerTriggerGet() *core.Command {
	return &core.Command{
		Short:     `Get the trigger associated with the specified ID.`,
		Long:      `Get the trigger associated with the specified ID.`,
		Namespace: "container",
		Resource:  "trigger",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetTriggerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
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
			request := args.(*container.GetTriggerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.GetTrigger(request)
		},
	}
}

func containerTriggerList() *core.Command {
	return &core.Command{
		Short: `List all triggers the caller can access (read permission).`,
		Long: `By default, the triggers listed are ordered by creation date in ascending order. This can be modified via the ` + "`" + `order_by` + "`" + ` field.

Additional parameters can be set in the query to filter, such as ` + "`" + `organization_id` + "`" + `, ` + "`" + `project_id` + "`" + `, ` + "`" + `namespace_id` + "`" + `, ` + "`" + `container_id` + "`" + ` or ` + "`" + `trigger_type` + "`" + `.`,
		Namespace: "container",
		Resource:  "trigger",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListTriggersRequest{}),
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
				},
			},
			{
				Name:       "project-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "namespace-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "container-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "trigger-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_source_type",
					"cron",
					"sqs",
					"nats",
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
			request := args.(*container.ListTriggersRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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

func containerTriggerUpdate() *core.Command {
	return &core.Command{
		Short: `Update the trigger associated with the specified ID.`,
		Long: `When updating a trigger, you cannot specify a different source type than the one already set.
Only fields present in the request are updated; others are left untouched.`,
		Namespace: "container",
		Resource:  "trigger",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateTriggerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
				Short:      `ID of the trigger to update.`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "name",
				Short:      `Name of the trigger.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `Description of the trigger.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of the trigger.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-config.http-path",
				Short:      `The HTTP path to send the request to (e.g., "/my-webhook-endpoint").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "destination-config.http-method",
				Short:      `The HTTP method to use when sending the request (e.g., get, post, put, patch, delete). Must be specified as lowercase.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_http_method",
					"get",
					"post",
					"put",
					"patch",
					"delete",
				},
			},
			{
				Name:       "cron-config.schedule",
				Short:      `UNIX cron schedule to run job (e.g., "* * * * *").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.timezone",
				Short:      `Timezone for the cron schedule, in tz database format (e.g., "Europe/Paris").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.body",
				Short:      `Body to send to the container when the trigger is invoked.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cron-config.headers.{key}",
				Short:      `Additional headers to send to the container when the trigger is invoked.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.region",
				Short:      `The region where the SQS queue is hosted (e.g., "fr-par", "nl-ams").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.endpoint",
				Short:      `Endpoint URL to use to access SQS (e.g., "https://sqs.mnq.fr-par.scaleway.com").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.access-key-id",
				Short:      `The access key for accessing the SQS queue.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.secret-access-key",
				Short:      `The secret key for accessing the SQS queue.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "sqs-config.queue-url",
				Short:      `The URL of the SQS queue to monitor for messages.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.server-urls.{index}",
				Short:      `The URLs of the NATS server (e.g., "nats://nats.mnq.fr-par.scaleway.com:4222").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.subject",
				Short:      `NATS subject to subscribe to (e.g., "my-subject").`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "nats-config.credentials-file-content",
				Short:      `The content of the NATS credentials file.`,
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
			request := args.(*container.UpdateTriggerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.UpdateTrigger(request)
		},
	}
}

func containerTriggerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete the trigger associated with the specified ID.`,
		Long:      `This action **cannot** be undone.`,
		Namespace: "container",
		Resource:  "trigger",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteTriggerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "trigger-id",
				Short:      `ID of the trigger to delete.`,
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
			request := args.(*container.DeleteTriggerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)

			return api.DeleteTrigger(request)
		},
	}
}
