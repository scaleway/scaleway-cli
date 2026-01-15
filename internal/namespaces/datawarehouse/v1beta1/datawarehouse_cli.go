// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package datawarehouse

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	datawarehouse "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		datawarehouseRoot(),
		datawarehouseDeployment(),
		datawarehousePreset(),
		datawarehouseVersion(),
		datawarehouseUser(),
		datawarehouseDatabase(),
		datawarehouseEndpoint(),
		datawarehousePresetList(),
		datawarehouseVersionList(),
		datawarehouseDeploymentList(),
		datawarehouseDeploymentGet(),
		datawarehouseDeploymentCreate(),
		datawarehouseDeploymentUpdate(),
		datawarehouseDeploymentDelete(),
		datawarehouseUserList(),
		datawarehouseUserCreate(),
		datawarehouseUserUpdate(),
		datawarehouseUserDelete(),
		datawarehouseDatabaseList(),
		datawarehouseDatabaseCreate(),
		datawarehouseDatabaseDelete(),
	)
}

func datawarehouseRoot() *core.Command {
	return &core.Command{
		Short:     `This API allows you to manage your Data Warehouse`,
		Long:      `Data Warehouse API.`,
		Namespace: "datawarehouse",
	}
}

func datawarehouseDeployment() *core.Command {
	return &core.Command{
		Short:     `Deployment management commands`,
		Long:      `A deployment is composed of one or multiple replicas.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
	}
}

func datawarehousePreset() *core.Command {
	return &core.Command{
		Short:     `List available presets`,
		Long:      `Data Warehouse preset to help you choose the best configuration.`,
		Namespace: "datawarehouse",
		Resource:  "preset",
	}
}

func datawarehouseVersion() *core.Command {
	return &core.Command{
		Short:     `List available Clickhouse® versions`,
		Long:      `ClickHouse® versions powering your deployment.`,
		Namespace: "datawarehouse",
		Resource:  "version",
	}
}

func datawarehouseUser() *core.Command {
	return &core.Command{
		Short:     `User management commands`,
		Long:      `Manage users associated with a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "user",
	}
}

func datawarehouseDatabase() *core.Command {
	return &core.Command{
		Short:     `Database management commands`,
		Long:      `Manage databases within a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "database",
	}
}

func datawarehouseEndpoint() *core.Command {
	return &core.Command{
		Short:     `Endpoint management commands`,
		Long:      `Manage endpoints associated with a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "endpoint",
	}
}

func datawarehousePresetList() *core.Command {
	return &core.Command{
		Short:     `List available presets`,
		Long:      `List available presets.`,
		Namespace: "datawarehouse",
		Resource:  "preset",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.ListPresetsRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.ListPresetsRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListPresets(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Presets, nil
		},
	}
}

func datawarehouseVersionList() *core.Command {
	return &core.Command{
		Short:     `List available ClickHouse® versions`,
		Long:      `List available ClickHouse® versions.`,
		Namespace: "datawarehouse",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version",
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
			request := args.(*datawarehouse.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
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

func datawarehouseDeploymentList() *core.Command {
	return &core.Command{
		Short:     `List deployments`,
		Long:      `List all deployments in the specified region, for a given Scaleway Project. By default, the deployments returned in the list are ordered by creation date in ascending order, though this can be modified via the order_by field. You can define additional parameters for your query, such as ` + "`" + `tags` + "`" + ` and ` + "`" + `name` + "`" + `. For the ` + "`" + `name` + "`" + ` parameter, the value you provide will be checked against the whole name string to see if it includes the string you put in the parameter.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.ListDeploymentsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tags.{index}",
				Short:      `List deployments with a given tag`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Lists deployments that match a name pattern`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering deployment listings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"created_at_desc",
					"created_at_asc",
					"name_asc",
					"name_desc",
				},
			},
			{
				Name:       "project-id",
				Short:      `Project ID the deployment belongs to`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Organization ID the deployment belongs to`,
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
			request := args.(*datawarehouse.ListDeploymentsRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
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

func datawarehouseDeploymentGet() *core.Command {
	return &core.Command{
		Short:     `Get a deployment`,
		Long:      `Retrieve information about a given deployment, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `deployment_id` + "`" + ` parameters. Its full details, including name, status are returned in the response object.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.GetDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.GetDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.GetDeployment(request)
		},
	}
}

func datawarehouseDeploymentCreate() *core.Command {
	return &core.Command{
		Short:     `Create a deployment`,
		Long:      `Create a new deployment.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.CreateDeploymentRequest{}),
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
				Short:      `Tags to apply to the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version",
				Short:      `ClickHouse® version to use for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "replica-count",
				Short:      `Number of replicas for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the initial user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-min",
				Short:      `Minimum CPU count for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-max",
				Short:      `Maximum CPU count for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "endpoints.{index}.private-network.private-network-id",
				Short:      `UUID of the Private Network`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "ram-per-cpu",
				Short:      `RAM per CPU count for the deployment (in GB)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.CreateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.CreateDeployment(request)
		},
	}
}

func datawarehouseDeploymentUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a deployment`,
		Long:      `Update the parameters of a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.UpdateDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
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
			{
				Name:       "cpu-min",
				Short:      `Minimum CPU count for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "cpu-max",
				Short:      `Maximum CPU count for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "replica-count",
				Short:      `Number of replicas for the deployment`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.UpdateDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.UpdateDeployment(request)
		},
	}
}

func datawarehouseDeploymentDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a deployment`,
		Long:      `Delete a given deployment, specified by the ` + "`" + `region` + "`" + ` and ` + "`" + `deployment_id` + "`" + ` parameters. Deleting a deployment is permanent, and cannot be undone. Upon deletion, all your data will be lost.`,
		Namespace: "datawarehouse",
		Resource:  "deployment",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.DeleteDeploymentRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment to delete`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.DeleteDeploymentRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.DeleteDeployment(request)
		},
	}
}

func datawarehouseUserList() *core.Command {
	return &core.Command{
		Short:     `List users associated with a deployment`,
		Long:      `List users associated with a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "user",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.ListUsersRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the user to filter by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering user listings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.ListUsersRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
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

func datawarehouseUserCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new user for a deployment`,
		Long:      `Create a new user for a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "user",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.CreateUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `Password for the user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-admin",
				Short:      `Indicates if the user is an administrator`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.CreateUserRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.CreateUser(request)
		},
	}
}

func datawarehouseUserUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an existing user for a deployment`,
		Long:      `Update an existing user for a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "user",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.UpdateUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the user`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "password",
				Short:      `New password for the user`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-admin",
				Short:      `Updates the user administrator permissions`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.UpdateUserRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.UpdateUser(request)
		},
	}
}

func datawarehouseUserDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a user from a deployment`,
		Long:      `Delete a user from a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "user",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.DeleteUserRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the user to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.DeleteUserRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
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

func datawarehouseDatabaseList() *core.Command {
	return &core.Command{
		Short:     `List databases within a deployment`,
		Long:      `List databases within a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "database",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.ListDatabasesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the database to filter by`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering database listings`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
					"size_asc",
					"size_desc",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.ListDatabasesRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatabases(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Databases, nil
		},
	}
}

func datawarehouseDatabaseCreate() *core.Command {
	return &core.Command{
		Short:     `Create a new database within a deployment`,
		Long:      `Create a new database within a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "database",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.CreateDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the database`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.CreateDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)

			return api.CreateDatabase(request)
		},
	}
}

func datawarehouseDatabaseDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a database from a deployment`,
		Long:      `Delete a database from a deployment.`,
		Namespace: "datawarehouse",
		Resource:  "database",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datawarehouse.DeleteDatabaseRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "deployment-id",
				Short:      `UUID of the deployment`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Name of the database to delete`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datawarehouse.DeleteDatabaseRequest)

			client := core.ExtractClient(ctx)
			api := datawarehouse.NewAPI(client)
			e = api.DeleteDatabase(request)
			if e != nil {
				return nil, e
			}

			return &core.SuccessResult{
				Resource: "database",
				Verb:     "delete",
			}, nil
		},
	}
}
