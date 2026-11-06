// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package dataviz

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/dataviz/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		datavizRoot(),
		datavizDeployment(),
		datavizVersion(),
		datavizNodeType(),
		datavizUser(),
		datavizEndpoint(),
		datavizDeploymentCreate(),
		datavizDeploymentUpdate(),
		datavizDeploymentUpgrade(),
		datavizDeploymentGet(),
		datavizDeploymentDelete(),
		datavizDeploymentList(),
		datavizVersionList(),
		datavizNodeTypeList(),
		datavizEndpointCreate(),
		datavizEndpointDelete(),
		datavizUserList(),
		datavizUserCreate(),
		datavizUserUpdate(),
		datavizUserDelete(),
	)
}

func datavizRoot() *core.Command {
	return &core.Command{
		Short:     `DataViz API`,
		Long:      `DataViz API.`,
		Namespace: "dataviz",
	}
}

func datavizDeployment() *core.Command {
	return &core.Command{
		Short:     `Manage your DataViz deployment`,
		Long:      `Manage your DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
	}
}

func datavizVersion() *core.Command {
	return &core.Command{
		Short:     `List your DataViz versions`,
		Long:      `List your DataViz versions.`,
		Namespace: "dataviz",
		Resource:  "version",
	}
}

func datavizNodeType() *core.Command {
	return &core.Command{
		Short:     `List available node types`,
		Long:      `List available node types.`,
		Namespace: "dataviz",
		Resource:  "node-type",
	}
}

func datavizUser() *core.Command {
	return &core.Command{
		Short:     `Manage your DataViz deployment users`,
		Long:      `Manage your DataViz deployment users.`,
		Namespace: "dataviz",
		Resource:  "user",
	}
}

func datavizEndpoint() *core.Command {
	return &core.Command{
		Short:     `Manage your DataViz deployment endpoints`,
		Long:      `Manage your DataViz deployment endpoints.`,
		Namespace: "dataviz",
		Resource:  "endpoint",
	}
}

func datavizDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new DataViz deployment`,
		Long:      `Create a new DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.CreateDeploymentRequest](),
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
				Short:      `The DataViz version to use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*dataviz.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.CreateDeployment(request, scw.WithContext(ctx))
		},
	}
}

func datavizDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a DataViz deployment`,
		Long:      `Update a DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.UpdateDeploymentRequest](),
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
			request := args.(*dataviz.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.UpdateDeployment(request, scw.WithContext(ctx))
		},
	}
}

func datavizDeploymentUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade a DataViz deployment`,
		Long:      `Upgrade a DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "upgrade",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.UpgradeDeploymentRequest](),
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
			request := args.(*dataviz.UpgradeDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.UpgradeDeployment(request, scw.WithContext(ctx))
		},
	}
}

func datavizDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Retrieve a specific DataViz deployment`,
		Long:      `Retrieve a specific DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.GetDeploymentRequest](),
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
			request := args.(*dataviz.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.GetDeployment(request, scw.WithContext(ctx))
		},
	}
}

func datavizDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a DataViz deployment`,
		Long:      `Delete a DataViz deployment.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.DeleteDeploymentRequest](),
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
			request := args.(*dataviz.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.DeleteDeployment(request, scw.WithContext(ctx))
		},
	}
}

func datavizDeploymentList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of DataViz deployments`,
		Long:      `Retrieve a list of DataViz deployments.`,
		Namespace: "dataviz",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.ListDeploymentsRequest](),
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
			request := args.(*dataviz.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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

func datavizVersionList() *core.Command {
	return &core.Command{
		Short:     `List available DataViz versions`,
		Long:      `List available DataViz versions.`,
		Namespace: "dataviz",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.ListVersionsRequest](),
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
			request := args.(*dataviz.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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

func datavizNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of available node types`,
		Long:      `Retrieve a list of available node types.`,
		Namespace: "dataviz",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.ListNodeTypesRequest](),
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
			request := args.(*dataviz.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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

func datavizEndpointCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new endpoint for a deployment`,
		Long:      `Create a new endpoint for a deployment.`,
		Namespace: "dataviz",
		Resource:  "endpoint",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.CreateEndpointRequest](),
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
			request := args.(*dataviz.CreateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.CreateEndpoint(request, scw.WithContext(ctx))
		},
	}
}

func datavizEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing endpoint`,
		Long:      `Delete an existing endpoint.`,
		Namespace: "dataviz",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.DeleteEndpointRequest](),
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
			request := args.(*dataviz.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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

func datavizUserList() *core.Command {
	return &core.Command{
		Short:     `Retrieve a list of deployment users`,
		Long:      `Retrieve a list of deployment users.`,
		Namespace: "dataviz",
		Resource:  "user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.ListUsersRequest](),
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
			request := args.(*dataviz.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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

func datavizUserCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new user`,
		Long:      `Create a new user.`,
		Namespace: "dataviz",
		Resource:  "user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.CreateUserRequest](),
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
			request := args.(*dataviz.CreateUserRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.CreateUser(request, scw.WithContext(ctx))
		},
	}
}

func datavizUserUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing user`,
		Long:      `Update an existing user.`,
		Namespace: "dataviz",
		Resource:  "user",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.UpdateUserRequest](),
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
			request := args.(*dataviz.UpdateUserRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)

			return api.UpdateUser(request, scw.WithContext(ctx))
		},
	}
}

func datavizUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an existing user`,
		Long:      `Delete an existing user.`,
		Namespace: "dataviz",
		Resource:  "user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeFor[dataviz.DeleteUserRequest](),
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
			request := args.(*dataviz.DeleteUserRequest)

			client := core.ExtractClient(ctx)
			api := dataviz.NewAPI(client)
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
