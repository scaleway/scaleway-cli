// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package searchdb

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	searchdb "github.com/scaleway/scaleway-sdk-go/api/searchdb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		searchdbRoot(),
		searchdbDeployment(),
		searchdbVersions(),
		searchdbNodeTypes(),
		searchdbUser(),
		searchdbEndpoint(),
		searchdbDeploymentCreate(),
		searchdbDeploymentUpdate(),
		searchdbDeploymentUpgrade(),
		searchdbDeploymentGet(),
		searchdbDeploymentDelete(),
		searchdbDeploymentList(),
		searchdbVersionsList(),
		searchdbNodeTypesList(),
		searchdbEndpointCreate(),
		searchdbEndpointDelete(),
		searchdbUserList(),
		searchdbUserCreate(),
		searchdbUserUpdate(),
		searchdbUserDelete(),
	)
}

func searchdbRoot() *core.Command {
	return &core.Command{
		Short:     `SearchDB API`,
		Long:      `SearchDB API.`,
		Namespace: "searchdb",
	}
}

func searchdbDeployment() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `Manage your SearchDB deployment.`,
		Namespace: "searchdb",
		Resource:  "deployment",
	}
}

func searchdbVersions() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `List your SearchDB versions.`,
		Namespace: "searchdb",
		Resource:  "versions",
	}
}

func searchdbNodeTypes() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `List available node types.`,
		Namespace: "searchdb",
		Resource:  "node-types",
	}
}

func searchdbUser() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `Manage your SearchDB deployment users.`,
		Namespace: "searchdb",
		Resource:  "user",
	}
}

func searchdbEndpoint() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `Manage your SearchDB deployment endpoint.`,
		Namespace: "searchdb",
		Resource:  "endpoint",
	}
}

func searchdbDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create searchdb resources`,
		Long:      `Create searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.CreateDeploymentRequest{}),
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
				Name:       "node-amount",
				Short:      `Number of nodes`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-type",
				Short:      `Node type`,
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
				Short:      `Define the type of the Volume`,
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
				Short:      `Define the size of the Volume`,
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
				Short:      `The Opensearch version to use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*searchdb.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.CreateDeployment(request)
		},
	}
}

func searchdbDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update searchdb resources`,
		Long:      `Update searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.UpdateDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment to update`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `Tags of a deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*searchdb.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.UpdateDeployment(request)
		},
	}
}

func searchdbDeploymentUpgrade() *core.Command {
	return &core.Command{
		Short:     `Upgrade searchdb resources`,
		Long:      `Upgrade searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "upgrade",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.UpgradeDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the Deployment to upgrade`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-amount",
				Short:      `Amount of node upgrade target`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "volume-size-bytes",
				Short:      `Volume size upgrade target`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*searchdb.UpgradeDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.UpgradeDeployment(request)
		},
	}
}

func searchdbDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Get searchdb resources`,
		Long:      `Get searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.GetDeploymentRequest{}),
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
			request := args.(*searchdb.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.GetDeployment(request)
		},
	}
}

func searchdbDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete searchdb resources`,
		Long:      `Delete searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.DeleteDeploymentRequest{}),
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
			request := args.(*searchdb.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.DeleteDeployment(request)
		},
	}
}

func searchdbDeploymentList() *core.Command {
	return &core.Command{
		Short:     `List searchdb resources`,
		Long:      `List searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.ListDeploymentsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `ID of the Project containing the deployments`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Define the order of the returned deployments`,
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
				Short:      `Filter by tag, only deployments with one or more matching tags will be returned`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Deployment name to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `Engine version to filter for`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `ID of the Organization containing the deployments`,
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
			request := args.(*searchdb.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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

func searchdbVersionsList() *core.Command {
	return &core.Command{
		Short:     `List searchdb resources`,
		Long:      `List searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "versions",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Define the order of the returned version`,
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
				Short:      `Filter by version`,
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
			request := args.(*searchdb.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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

func searchdbNodeTypesList() *core.Command {
	return &core.Command{
		Short:     `List searchdb resources`,
		Long:      `List searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "node-types",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Sort order of nodes in the response (name, vcpus or memory)`,
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
			request := args.(*searchdb.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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

func searchdbEndpointCreate() *core.Command {
	return &core.Command{
		Short:     `Create searchdb resources`,
		Long:      `Create searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "endpoint",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.CreateEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment for which to create an endpoint`,
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
			request := args.(*searchdb.CreateEndpointRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.CreateEndpoint(request)
		},
	}
}

func searchdbEndpointDelete() *core.Command {
	return &core.Command{
		Short:     `Delete searchdb resources`,
		Long:      `Delete searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "endpoint",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.DeleteEndpointRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "endpoint-id",
				Short:      `ID of the endpoint to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*searchdb.DeleteEndpointRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			e = api.DeleteEndpoint(request)
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

func searchdbUserList() *core.Command {
	return &core.Command{
		Short:     `List searchdb resources`,
		Long:      `List searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.ListUsersRequest{}),
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
			request := args.(*searchdb.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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

func searchdbUserCreate() *core.Command {
	return &core.Command{
		Short:     `Create searchdb resources`,
		Long:      `Create searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.CreateUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment in which to create the user`,
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
			request := args.(*searchdb.CreateUserRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.CreateUser(request)
		},
	}
}

func searchdbUserUpdate() *core.Command {
	return &core.Command{
		Short:     `Update searchdb resources`,
		Long:      `Update searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "user",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.UpdateUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment in which to create the user`,
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
			request := args.(*searchdb.UpdateUserRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)

			return api.UpdateUser(request)
		},
	}
}

func searchdbUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete searchdb resources`,
		Long:      `Delete searchdb resources.`,
		Namespace: "searchdb",
		Resource:  "user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(searchdb.DeleteUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `ID of the deployment in which to create the user`,
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
			request := args.(*searchdb.DeleteUserRequest)

			client := core.ExtractClient(ctx)
			api := searchdb.NewAPI(client)
			e = api.DeleteUser(request)
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
