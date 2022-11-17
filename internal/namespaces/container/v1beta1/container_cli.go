// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package container

import (
	"context"
	"fmt"
	"reflect"
	"strconv"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
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
		containerContext(),
		containerCron(),
		containerDomain(),
		containerToken(),
		containerNamespaceList(),
		containerNamespaceGet(),
		containerNamespaceCreate(),
		containerNamespaceUpdate(),
		containerNamespaceDelete(),
		containerContainerList(),
		containerContainerGet(),
		containerContainerCreate(),
		containerContainerUpdate(),
		containerContainerDelete(),
		containerContainerDeploy(),
		containerContextCreate(),
		containerContextStart(),
		containerContextStop(),
		containerContextDelete(),
		containerCronList(),
		containerCronGet(),
		containerCronDelete(),
		containerDomainList(),
		containerDomainGet(),
		containerDomainCreate(),
		containerDomainDelete(),
		containerTokenCreate(),
		containerTokenGet(),
		containerTokenList(),
		containerTokenDelete(),
	)
}
func containerRoot() *core.Command {
	return &core.Command{
		Short:     `Container as a Service API`,
		Long:      ``,
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

func containerContext() *core.Command {
	return &core.Command{
		Short:     `Context management commands`,
		Long:      `Context management commands.`,
		Namespace: "container",
		Resource:  "context",
	}
}

func containerCron() *core.Command {
	return &core.Command{
		Short:     `Cron management commands`,
		Long:      `Cron management commands.`,
		Namespace: "container",
		Resource:  "cron",
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

func containerToken() *core.Command {
	return &core.Command{
		Short:     `Token management commands`,
		Long:      `Token management commands.`,
		Namespace: "container",
		Resource:  "token",
	}
}

func containerNamespaceList() *core.Command {
	return &core.Command{
		Short:     `List all your namespaces`,
		Long:      `List all your namespaces.`,
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
			request := args.(*container.ListNamespacesRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			resp, err := api.ListNamespaces(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Namespaces, nil

		},
	}
}

func containerNamespaceGet() *core.Command {
	return &core.Command{
		Short:     `Get a namespace`,
		Long:      `Get the namespace associated with the given id.`,
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
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.GetNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.GetNamespace(request)

		},
	}
}

func containerNamespaceCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new namespace`,
		Long:      `Create a new namespace.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("cns"),
			},
			{
				Name:       "environment-variables.{key}",
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
			request := args.(*container.CreateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.CreateNamespace(request)

		},
	}
}

func containerNamespaceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing namespace`,
		Long:      `Update the space associated with the given id.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
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
			request := args.(*container.UpdateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.UpdateNamespace(request)

		},
	}
}

func containerNamespaceDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing namespace`,
		Long:      `Delete the namespace associated with the given id.`,
		Namespace: "container",
		Resource:  "namespace",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteNamespaceRequest{}),
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
			request := args.(*container.DeleteNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeleteNamespace(request)

		},
	}
}

func containerContainerList() *core.Command {
	return &core.Command{
		Short:     `List all your containers`,
		Long:      `List all your containers.`,
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
			request := args.(*container.ListContainersRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			resp, err := api.ListContainers(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Containers, nil

		},
	}
}

func containerContainerGet() *core.Command {
	return &core.Command{
		Short:     `Get a container`,
		Long:      `Get the container associated with the given id.`,
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
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.GetContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.GetContainer(request)

		},
	}
}

func containerContainerCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new container`,
		Long:      `Create a new container.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
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
				Default:    core.RandomValueGenerator("ctnr"),
			},
			{
				Name:       "environment-variables.{key}",
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
				Name:       "registry-image",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-concurrency",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_protocol", "http1", "h2c"},
			},
			{
				Name:       "port",
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
				Deprecated: true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.CreateContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.CreateContainer(request)

		},
	}
}

func containerContainerUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing container`,
		Long:      `Update the container associated with the given id.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.UpdateContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "environment-variables.{key}",
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
				Name:       "registry-image",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "max-concurrency",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "protocol",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_protocol", "http1", "h2c"},
			},
			{
				Name:       "port",
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
				Deprecated: true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.UpdateContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.UpdateContainer(request)

		},
	}
}

func containerContainerDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a container`,
		Long:      `Delete the container associated with the given id.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.DeleteContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeleteContainer(request)

		},
	}
}

func containerContainerDeploy() *core.Command {
	return &core.Command{
		Short:     `Deploy a container`,
		Long:      `Deploy a container associated with the given id.`,
		Namespace: "container",
		Resource:  "container",
		Verb:      "deploy",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeployContainerRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.DeployContainerRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeployContainer(request)

		},
	}
}

type createContextRequest struct {
	Name string   `json:"-"`
	Zone scw.Zone `json:"-"`
	Size uint64   `json:"-"`
}

func containerContextCreate() *core.Command {
	return &core.Command{
		Short:     `Create context storage`,
		Long:      `Create block storage that one can attach to a container context.`,
		Namespace: "container",
		Resource:  "context",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(createContextRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name: "size",
				// TODO: shove off with: docker builder prune --keep-storage 20000000000 -f
				Short: "Size of block storage in GB",
				ValidateFunc: func(argSpec *core.ArgSpec, value interface{}) error {
					if _, err := strconv.ParseUint(value.(string), 10, 64); err != nil {
						return &core.CliError{
							Err:  fmt.Errorf("invalid size %s", value),
							Hint: "Size should be an integer of gigabytes",
						}
					}
					return nil
				},
			},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*createContextRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)
			return api.CreateVolume(&instance.CreateVolumeRequest{
				Zone:       request.Zone,
				Name:       request.Name,
				Tags:       []string{"builder", "b-a-a-s", request.Name},
				VolumeType: instance.VolumeVolumeTypeBSSD,
				Size:       scw.SizePtr(scw.Size(request.Size) * scw.GB),
			})
		},
	}
}

func containerContextStart() *core.Command { return nil }
func containerContextStop() *core.Command  { return nil }

type deleteContextRequest struct {
	Name string   `json:"-"`
	Zone scw.Zone `json:"-"`
	Size uint64   `json:"-"`
}

func containerContextDelete() *core.Command {
	return &core.Command{
		Short:     `Delete context`,
		Long:      `Stop a context and shutdown its compute resources.`,
		Namespace: "container",
		Resource:  "context",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(deleteContextRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.ZoneArgSpec(),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*deleteContextRequest)

			client := core.ExtractClient(ctx)
			api := instance.NewAPI(client)

			x := instance.VolumeVolumeTypeBSSD
			response, err := api.ListVolumes(&instance.ListVolumesRequest{
				Zone:       request.Zone,
				VolumeType: &x,
				Tags:       []string{"builder", "b-a-a-s", request.Name},
				Name:       scw.StringPtr(request.Name),
			})
			if err != nil {
				return nil, err
			}
			if response.TotalCount != 1 {
				return nil, fmt.Errorf("Could not find volume named %q", request.Name)
			}

			err = api.DeleteVolume(&instance.DeleteVolumeRequest{
				Zone:     request.Zone,
				VolumeID: response.Volumes[0].ID,
			})
			return nil, err
		},
	}
}
func containerCronList() *core.Command {
	return &core.Command{
		Short:     `List all your crons`,
		Long:      `List all your crons.`,
		Namespace: "container",
		Resource:  "cron",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListCronsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "container-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.ListCronsRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			resp, err := api.ListCrons(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Crons, nil

		},
	}
}

func containerCronGet() *core.Command {
	return &core.Command{
		Short:     `Get a cron`,
		Long:      `Get the cron associated with the given id.`,
		Namespace: "container",
		Resource:  "cron",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetCronRequest{}),
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
			request := args.(*container.GetCronRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.GetCron(request)

		},
	}
}

func containerCronDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing cron`,
		Long:      `Delete the cron associated with the given id.`,
		Namespace: "container",
		Resource:  "cron",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteCronRequest{}),
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
			request := args.(*container.DeleteCronRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeleteCron(request)

		},
	}
}

func containerDomainList() *core.Command {
	return &core.Command{
		Short:     `List all domain name bindings`,
		Long:      `List all domain name bindings.`,
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
				EnumValues: []string{"created_at_asc", "created_at_desc", "hostname_asc", "hostname_desc"},
			},
			{
				Name:       "container-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.ListDomainsRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			resp, err := api.ListDomains(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Domains, nil

		},
	}
}

func containerDomainGet() *core.Command {
	return &core.Command{
		Short:     `Get a domain name binding`,
		Long:      `Get a domain name binding.`,
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
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.GetDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.GetDomain(request)

		},
	}
}

func containerDomainCreate() *core.Command {
	return &core.Command{
		Short:     `Create a domain name binding`,
		Long:      `Create a domain name binding.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "hostname",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.CreateDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.CreateDomain(request)

		},
	}
}

func containerDomainDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a domain name binding`,
		Long:      `Delete a domain name binding.`,
		Namespace: "container",
		Resource:  "domain",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteDomainRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "domain-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.DeleteDomainRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeleteDomain(request)

		},
	}
}

func containerTokenCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new revocable token`,
		Long:      `Create a new revocable token.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.CreateTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "container-id",
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
				Name:       "description",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "expires-at",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.CreateTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.CreateToken(request)

		},
	}
}

func containerTokenGet() *core.Command {
	return &core.Command{
		Short:     `Get a token`,
		Long:      `Get a token.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.GetTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.GetTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.GetToken(request)

		},
	}
}

func containerTokenList() *core.Command {
	return &core.Command{
		Short:     `List all tokens`,
		Long:      `List all tokens.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.ListTokensRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "container-id",
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
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.ListTokensRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			resp, err := api.ListTokens(request, scw.WithAllPages())
			if err != nil {
				return nil, err
			}
			return resp.Tokens, nil

		},
	}
}

func containerTokenDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a token`,
		Long:      `Delete a token.`,
		Namespace: "container",
		Resource:  "token",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(container.DeleteTokenRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "token-id",
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*container.DeleteTokenRequest)

			client := core.ExtractClient(ctx)
			api := container.NewAPI(client)
			return api.DeleteToken(request)

		},
	}
}
