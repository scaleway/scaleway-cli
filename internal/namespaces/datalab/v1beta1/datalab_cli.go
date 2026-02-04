// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package datalab

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	datalab "github.com/scaleway/scaleway-sdk-go/api/datalab/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		datalabRoot(),
		datalabDatalab(),
		datalabNodeType(),
		datalabNotebookVersion(),
		datalabClusterVersion(),
		datalabDatalabCreate(),
		datalabDatalabGet(),
		datalabDatalabList(),
		datalabDatalabUpdate(),
		datalabDatalabDelete(),
		datalabNodeTypeList(),
		datalabNotebookVersionList(),
		datalabClusterVersionList(),
	)
}

func datalabRoot() *core.Command {
	return &core.Command{
		Short:     `Data Lab API for Apache Spark™`,
		Long:      `Data Lab API.`,
		Namespace: "datalab",
	}
}

func datalabDatalab() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `Manage your Data Labs.`,
		Namespace: "datalab",
		Resource:  "datalab",
	}
}

func datalabNodeType() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `List available node types.`,
		Namespace: "datalab",
		Resource:  "node-type",
	}
}

func datalabNotebookVersion() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `List available notebook versions.`,
		Namespace: "datalab",
		Resource:  "notebook-version",
	}
}

func datalabClusterVersion() *core.Command {
	return &core.Command{
		Short:     ``,
		Long:      `Lists the Spark versions available for Data Lab creation.`,
		Namespace: "datalab",
		Resource:  "cluster-version",
	}
}

func datalabDatalabCreate() *core.Command {
	return &core.Command{
		Short:     `Create datalab resources`,
		Long:      `Create datalab resources.`,
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.CreateDatalabRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ProjectIDArgSpec(),
			{
				Name:       "name",
				Short:      `The name of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The description of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "main.node-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "worker.node-type",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "worker.node-count",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "has-notebook",
				Short:      `Select this option to include a notebook as part of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "spark-version",
				Short:      `The version of Spark running inside the Data Lab, available options can be viewed at ListClusterVersions.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "total-storage.type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_type",
					"sbs_5k",
				},
			},
			{
				Name:       "total-storage.size",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "private-network-id",
				Short:      `The unique identifier of the private network the Data Lab will be attached to.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datalab.CreateDatalabRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)

			return api.CreateDatalab(request)
		},
	}
}

func datalabDatalabGet() *core.Command {
	return &core.Command{
		Short:     `Get datalab resources`,
		Long:      `Get datalab resources.`,
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.GetDatalabRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "datalab-id",
				Short:      `The unique identifier of the Data Lab`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datalab.GetDatalabRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)

			return api.GetDatalab(request)
		},
	}
}

func datalabDatalabList() *core.Command {
	return &core.Command{
		Short:     `List datalab resources`,
		Long:      `List datalab resources.`,
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.ListDatalabsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      `The unique identifier of the project whose Data Labs you want to list.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The name of the Data Lab you want to list.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The tags associated with the Data Lab you want to list.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `The order by field, available options are ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + `, ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `, ` + "`" + `updated_at_asc` + "`" + `, ` + "`" + `updated_at_desc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
					"created_at_asc",
					"created_at_desc",
					"updated_at_asc",
					"updated_at_desc",
				},
			},
			{
				Name:       "organization-id",
				Short:      `The unique identifier of the organization whose Data Labs you want to list.`,
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
			request := args.(*datalab.ListDatalabsRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListDatalabs(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Datalabs, nil
		},
	}
}

func datalabDatalabUpdate() *core.Command {
	return &core.Command{
		Short:     `Update datalab resources`,
		Long:      `Update datalab resources.`,
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.UpdateDatalabRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "datalab-id",
				Short:      `The unique identifier of the Data Lab.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `The updated name of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "description",
				Short:      `The updated description of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "tags.{index}",
				Short:      `The updated tags of the Data Lab.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "node-count",
				Short:      `The updated node count of the Data Lab. Scale up or down the number of worker nodes.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datalab.UpdateDatalabRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)

			return api.UpdateDatalab(request)
		},
	}
}

func datalabDatalabDelete() *core.Command {
	return &core.Command{
		Short:     `Delete datalab resources`,
		Long:      `Delete datalab resources.`,
		Namespace: "datalab",
		Resource:  "datalab",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.DeleteDatalabRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "datalab-id",
				Short:      `The unique identifier of the Data Lab.`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datalab.DeleteDatalabRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)

			return api.DeleteDatalab(request)
		},
	}
}

func datalabNodeTypeList() *core.Command {
	return &core.Command{
		Short:     `List datalab resources`,
		Long:      `List datalab resources.`,
		Namespace: "datalab",
		Resource:  "node-type",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.ListNodeTypesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The order by field. Available fields are ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + `, ` + "`" + `vcpus_asc` + "`" + `, ` + "`" + `vcpus_desc` + "`" + `, ` + "`" + `memory_gigabytes_asc` + "`" + `, ` + "`" + `memory_gigabytes_desc` + "`" + `, ` + "`" + `vram_bytes_asc` + "`" + `, ` + "`" + `vram_bytes_desc` + "`" + `, ` + "`" + `gpus_asc` + "`" + `, ` + "`" + `gpus_desc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"name_asc",
					"name_desc",
					"vcpus_asc",
					"vcpus_desc",
					"memory_gigabytes_asc",
					"memory_gigabytes_desc",
					"vram_bytes_asc",
					"vram_bytes_desc",
					"gpus_asc",
					"gpus_desc",
				},
			},
			{
				Name:       "targets.{index}",
				Short:      `Filter based on the target of the nodes. Allows to filter the nodes based on their purpose which can be main or worker node.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"unknown_target",
					"notebook",
					"worker",
				},
			},
			{
				Name:       "resource-type",
				Short:      `Filter based on node type ( ` + "`" + `cpu` + "`" + `/` + "`" + `gpu` + "`" + `/` + "`" + `all` + "`" + ` ).`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{
					"all",
					"gpu",
					"cpu",
				},
			},
			core.RegionArgSpec(
				scw.RegionFrPar,
				scw.Region(core.AllLocalities),
			),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*datalab.ListNodeTypesRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)
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

func datalabNotebookVersionList() *core.Command {
	return &core.Command{
		Short:     `List datalab resources`,
		Long:      `List datalab resources.`,
		Namespace: "datalab",
		Resource:  "notebook-version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.ListNotebookVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The order by field. Available options are ` + "`" + `name_asc` + "`" + ` and ` + "`" + `name_desc` + "`" + `.`,
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
			request := args.(*datalab.ListNotebookVersionsRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListNotebookVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Notebooks, nil
		},
	}
}

func datalabClusterVersionList() *core.Command {
	return &core.Command{
		Short:     `List datalab resources`,
		Long:      `List datalab resources.`,
		Namespace: "datalab",
		Resource:  "cluster-version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(datalab.ListClusterVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `The order by field.`,
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
			request := args.(*datalab.ListClusterVersionsRequest)

			client := core.ExtractClient(ctx)
			api := datalab.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListClusterVersions(request, opts...)
			if err != nil {
				return nil, err
			}

			return resp.Clusters, nil
		},
	}
}
