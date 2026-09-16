// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package messageq

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	messageq "github.com/scaleway/scaleway-sdk-go/api/messageq/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		messageqRoot(),
		messageqDeployment(),
		messageqVersion(),
		messageqNodeType(),
		messageqUser(),
		messageqEndpoint(),
		messageqDeploymentCreate(),
		messageqDeploymentUpdate(),
		messageqDeploymentUpgrade(),
		messageqDeploymentGet(),
		messageqDeploymentDelete(),
		messageqDeploymentList(),
		messageqVersionList(),
		messageqNodeTypeList(),
		messageqEndpointCreate(),
		messageqEndpointDelete(),
		messageqUserList(),
		messageqUserCreate(),
		messageqUserUpdate(),
		messageqUserDelete(),
	)
}

func messageqRoot() *core.Command {
	return &core.Command{
		Short:     `MessageQ API`,
		Long:      `MessageQ API.`,
		Namespace: "messageq",
	}
}

func messageqDeployment() *core.Command {
	return &core.Command{
		Short:     `Manage your MessageQ deployment`,
		Long:      `Manage your MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
	}
}

func messageqVersion() *core.Command {
	return &core.Command{
		Short:     `List your MessageQ versions`,
		Long:      `List your MessageQ versions.`,
		Namespace: "messageq",
		Resource:  "version",
	}
}

func messageqNodeType() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List available node types.`,
		Namespace: "messageq",
		Resource:  "node-type",
	}
}

func messageqUser() *core.Command {
	return &core.Command{
		Short:     `Manage your MessageQ deployment users`,
		Long:      `Manage your MessageQ deployment users.`,
		Namespace: "messageq",
		Resource:  "user",
	}
}

func messageqEndpoint() *core.Command {
	return &core.Command{
		Short:     `Manage your MessageQ deployment endpoints`,
		Long:      `Manage your MessageQ deployment endpoints.`,
		Namespace: "messageq",
		Resource:  "endpoint",
	}
}

func messageqDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new MessageQ deployment`,
		Long:      `Create a new MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.CreateDeploymentRequest](),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `Name of the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-count",
				Short:      `Number of nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Node type to use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "user-name",
				Short:      `Username for the deployment user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the deployment user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume.type",
				Short:      `Type of the Volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"sbs_5k",
					"sbs_15k",
				},
			},
			{
				Name:       "volume.size-bytes",
				Short:      `Size of the Volume`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.public",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.private-network-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `The MessageQ version to use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.CreateDeployment(request, scw.WithContext(ctx))
		},
	}
}

func messageqDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a MessageQ deployment`,
		Long:      `Update a MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.UpdateDeploymentRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `New name for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags to update`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.UpdateDeployment(request, scw.WithContext(ctx))
		},
	}
}

func messageqDeploymentUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade a MessageQ deployment`,
		Long:      `Upgrade a MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "upgrade",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.UpgradeDeploymentRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-count",
				Short:      `Target number of nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-size-bytes",
				Short:      `Target volume size`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.UpgradeDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.UpgradeDeployment(request, scw.WithContext(ctx))
		},
	}
}

func messageqDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a specific MessageQ deployment`,
		Long:      `Retrieve a specific MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.GetDeploymentRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.GetDeployment(request, scw.WithContext(ctx))
		},
	}
}

func messageqDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a MessageQ deployment`,
		Long:      `Delete a MessageQ deployment.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.DeleteDeploymentRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.DeleteDeployment(request, scw.WithContext(ctx))
		},
	}
}

func messageqDeploymentList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of MessageQ deployments`,
		Long:      `Retrieve a list of MessageQ deployments.`,
		Namespace: "messageq",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.ListDeploymentsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `Project ID to filter for, only deployments from this Project will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Sort order for deployments in the response`,
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
				Name:       "tags.{index}",
				Short:      `Tags to filter for, only deployments with one or more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Deployment name to filter for, only deployments with this string within their name will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID to filter for, only deployments from this Organization will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDeployments(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Deployments, nil
		},
	}
}

func messageqVersionList() *core.Command {
	return &core.Command{
		Short:     `List available MessageQ versions`,
		Long:      `List available MessageQ versions.`,
		Namespace: "messageq",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.ListVersionsRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for versions in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"version_asc",
					"version_desc",
				},
			},
			{
				Name:       "version",
				Short:      `Engine version to filter for, only versions with this version will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Versions, nil
		},
	}
}

func messageqNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of available node types`,
		Long:      `Retrieve a list of available node types.`,
		Namespace: "messageq",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.ListNodeTypesRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order for versions in the response`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
					"vcpus_asc",
					"vcpus_desc",
					"memory_asc",
					"memory_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNodeTypes(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.NodeTypes, nil
		},
	}
}

func messageqEndpointCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new endpoint for a deployment`,
		Long:      `Create a new endpoint for a deployment.`,
		Namespace: "messageq",
		Resource:  "endpoint",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.CreateEndpointRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoint-spec.public",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoint-spec.private-network.private-network-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.CreateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.CreateEndpoint(request, scw.WithContext(ctx))
		},
	}
}

func messageqEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing endpoint`,
		Long:      `Delete an existing endpoint.`,
		Namespace: "messageq",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.DeleteEndpointRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `ID of the endpoint`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			e = api.DeleteEndpoint(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "endpoint",
				Verb:     "delete",
			}, nil
		},
	}
}

func messageqUserList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of deployment users`,
		Long:      `Retrieve a list of deployment users.`,
		Namespace: "messageq",
		Resource:  "user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.ListUsersRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "name",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "deployment-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages(), scw.WithContext(ctx)}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListUsers(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Users, nil
		},
	}
}

func messageqUserCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new user`,
		Long:      `Create a new user.`,
		Namespace: "messageq",
		Resource:  "user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.CreateUserRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username of the deployment user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the deployment user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.CreateUserRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.CreateUser(request, scw.WithContext(ctx))
		},
	}
}

func messageqUserUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing user`,
		Long:      `Update an existing user.`,
		Namespace: "messageq",
		Resource:  "user",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.UpdateUserRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username of the deployment user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password of the deployment user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.UpdateUserRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)

			return api.UpdateUser(request, scw.WithContext(ctx))
		},
	}
}

func messageqUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing user`,
		Long:      `Delete an existing user.`,
		Namespace: "messageq",
		Resource:  "user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[messageq.DeleteUserRequest](),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "username",
				Short:      `Username of the deployment user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*messageq.DeleteUserRequest)

			client := core.ExtractClient(ctx)
			api := messageq.NewAPI(client)
			e = api.DeleteUser(request, scw.WithContext(ctx))
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "user",
				Verb:     "delete",
			}, nil
		},
	}
}
