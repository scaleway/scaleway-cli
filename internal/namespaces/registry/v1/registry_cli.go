// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package registry

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		registryRoot(),
		registryNamespace(),
		registryImage(),
		registryTag(),
		registryNamespaceList(),
		registryNamespaceGet(),
		registryNamespaceCreate(),
		registryNamespaceUpdate(),
		registryNamespaceDelete(),
		registryImageList(),
		registryImageGet(),
		registryImageUpdate(),
		registryImageDelete(),
		registryTagList(),
		registryTagGet(),
		registryTagDelete(),
	)
}
func registryRoot() *core.Command {
	return &core.Command{
		Short:     `Container Registry API`,
		Long:      `Container Registry API.`,
		Namespace: "registry",
	}
}

func registryNamespace() *core.Command {
	return &core.Command{
		Short: `Namespace management commands`,
		Long: `A namespace is a collection of container images, each bearing the unique identifier of that namespace. A namespace can be either public or private, by default.

Each namespace must have a globally unique name within its region. This means no namespaces in the same region can bear the same name.

You can use namespace privacy policies to specify whether everyone has the right to pull an image from a namespace or not. When an image is in a public namespace, anyone is able to pull it. You can set your namespace to private if you want to restrict access.`,
		Namespace: "registry",
		Resource:  "namespace",
	}
}

func registryImage() *core.Command {
	return &core.Command{
		Short: `Image management commands`,
		Long: `An image represents a container image. A container image is a file that includes all the requirements and instructions of a complete and executable version of an application. When running, it becomes one or multiple instances of that application.

The visibility of an image can be public - when anyone can pull it, private - when only users within your organization can pull it, or inherited from the namespace visibility - which is the default. The visibility of your image can be changed using the [update image endpoit](#path-images-update-an-image).`,
		Namespace: "registry",
		Resource:  "image",
	}
}

func registryTag() *core.Command {
	return &core.Command{
		Short:     `Tag management commands`,
		Long:      `Tags allow you to organize your container images. This gives you the possibility of sorting and filtering your images in any organizational pattern of your choice, which in turn helps you arrange, control and monitor your cloud resources. You can assign as many tags as you want to each image.`,
		Namespace: "registry",
		Resource:  "tag",
	}
}

func registryNamespaceList() *core.Command {
	return &core.Command{
		Short:     `List namespaces`,
		Long:      `List all namespaces in a specified region. By default, the namespaces listed are ordered by creation date in ascending order. This can be modified via the order_by field. You can also define additional parameters for your query, such as the ` + "`" + `instance_id` + "`" + ` and ` + "`" + `project_id` + "`" + ` parameters.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListNamespacesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering namespace listings. Possible values are ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `, ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + `, ` + "`" + `region` + "`" + `, ` + "`" + `status_asc` + "`" + ` and ` + "`" + `status_desc` + "`" + `. The default value is ` + "`" + `created_at_asc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "description_asc", "description_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "project-id",
				Short:      `Filter by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by the namespace name (exact match)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListNamespacesRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
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
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Region",
			},
			{
				FieldName: "Endpoint",
			},
			{
				FieldName: "IsPublic",
			},
			{
				FieldName: "Size",
			},
			{
				FieldName: "ImageCount",
			},
			{
				FieldName: "OrganizationID",
			},
			{
				FieldName: "ProjectID",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "StatusMessage",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Description",
			},
		}},
	}
}

func registryNamespaceGet() *core.Command {
	return &core.Command{
		Short:     `Get a namespace`,
		Long:      `Retrieve information about a given namespace, specified by its ` + "`" + `namespace_id` + "`" + ` and region. Full details about the namespace, such as ` + "`" + `description` + "`" + `, ` + "`" + `project_id` + "`" + `, ` + "`" + `status` + "`" + `, ` + "`" + `endpoint` + "`" + `, ` + "`" + `is_public` + "`" + `, ` + "`" + `size` + "`" + `, and ` + "`" + `image_count` + "`" + ` are returned in the response.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.GetNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.GetNamespace(request)

		},
	}
}

func registryNamespaceCreate() *core.Command {
	return &core.Command{
		Short:     `Create a namespace`,
		Long:      `Create a new Container Registry namespace. You must specify the namespace name and region in which you want it to be created. Optionally, you can specify the ` + "`" + `project_id` + "`" + ` and ` + "`" + `is_public` + "`" + ` in the request payload.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.CreateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Name of the namespace`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("ns"),
			},
			{
				Name:       "description",
				Short:      `Description of the namespace`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "is-public",
				Short:      `Defines whether or not namespace is public`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.OrganizationIDArgSpec(),
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.CreateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.CreateNamespace(request)

		},
	}
}

func registryNamespaceUpdate() *core.Command {
	return &core.Command{
		Short:     `Update a namespace`,
		Long:      `Update the parameters of a given namespace, specified by its ` + "`" + `namespace_id` + "`" + ` and ` + "`" + `region` + "`" + `. You can update the ` + "`" + `description` + "`" + ` and ` + "`" + `is_public` + "`" + ` parameters.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.UpdateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `ID of the namespace to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Namespace description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-public",
				Short:      `Defines whether or not the namespace is public`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.UpdateNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.UpdateNamespace(request)

		},
	}
}

func registryNamespaceDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a namespace`,
		Long:      `Delete a given namespace. You must specify, in the endpoint, the ` + "`" + `region` + "`" + ` and ` + "`" + `namespace_id` + "`" + ` parameters of the namespace you want to delete.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `UUID of the namespace`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.DeleteNamespaceRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.DeleteNamespace(request)

		},
	}
}

func registryImageList() *core.Command {
	return &core.Command{
		Short:     `List images`,
		Long:      `List all images in a specified region. By default, the images listed are ordered by creation date in ascending order. This can be modified via the order_by field. You can also define additional parameters for your query, such as the ` + "`" + `namespace_id` + "`" + ` and ` + "`" + `project_id` + "`" + ` parameters.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering image listings. Possible values are ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `, ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + `, ` + "`" + `region` + "`" + `, ` + "`" + `status_asc` + "`" + ` and ` + "`" + `status_desc` + "`" + `. The default value is ` + "`" + `created_at_asc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "namespace-id",
				Short:      `Filter by the namespace ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by the image name (exact match)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "project-id",
				Short:      `Filter by Project ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "organization-id",
				Short:      `Filter by Organization ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListImages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Images, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Size",
			},
			{
				FieldName: "Visibility",
			},
			{
				FieldName: "NamespaceID",
			},
			{
				FieldName: "Status",
			},
			{
				FieldName: "StatusMessage",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "Tags",
			},
		}},
	}
}

func registryImageGet() *core.Command {
	return &core.Command{
		Short:     `Get an image`,
		Long:      `Retrieve information about a given container image, specified by its ` + "`" + `image_id` + "`" + ` and region. Full details about the image, such as ` + "`" + `name` + "`" + `, ` + "`" + `namespace_id` + "`" + `, ` + "`" + `status` + "`" + `, ` + "`" + `visibility` + "`" + `, and ` + "`" + `size` + "`" + ` are returned in the response.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `UUID of the image`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.GetImageRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.GetImage(request)

		},
	}
}

func registryImageUpdate() *core.Command {
	return &core.Command{
		Short:     `Update an image`,
		Long:      `Update the parameters of a given image, specified by its ` + "`" + `image_id` + "`" + ` and ` + "`" + `region` + "`" + `. You can update the ` + "`" + `visibility` + "`" + ` parameter.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.UpdateImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `ID of the image to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "visibility",
				Short:      `Set to ` + "`" + `public` + "`" + ` to allow the image to be pulled without authentication. Else, set to  ` + "`" + `private` + "`" + `. Set to ` + "`" + `inherit` + "`" + ` to keep the same visibility configuration as the namespace`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"visibility_unknown", "inherit", "public", "private"},
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.UpdateImageRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.UpdateImage(request)

		},
	}
}

func registryImageDelete() *core.Command {
	return &core.Command{
		Short:     `Delete an image`,
		Long:      `Delete a given image. You must specify, in the endpoint, the ` + "`" + `region` + "`" + ` and ` + "`" + `image_id` + "`" + ` parameters of the image you want to delete.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `UUID of the image`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.DeleteImageRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.DeleteImage(request)

		},
	}
}

func registryTagList() *core.Command {
	return &core.Command{
		Short:     `List tags`,
		Long:      `List all tags for a given image, specified by region. By default, the tags listed are ordered by creation date in ascending order. This can be modified via the order_by field. You can also define additional parameters for your query, such as the ` + "`" + `name` + "`" + `.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListTagsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `UUID of the image`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Short:      `Criteria to use when ordering tag listings. Possible values are ` + "`" + `created_at_asc` + "`" + `, ` + "`" + `created_at_desc` + "`" + `, ` + "`" + `name_asc` + "`" + `, ` + "`" + `name_desc` + "`" + `, ` + "`" + `region` + "`" + `, ` + "`" + `status_asc` + "`" + ` and ` + "`" + `status_desc` + "`" + `. The default value is ` + "`" + `created_at_asc` + "`" + `.`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "name",
				Short:      `Filter by the tag name (exact match)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw, scw.Region(core.AllLocalities)),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListTagsRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			if request.Region == scw.Region(core.AllLocalities) {
				opts = append(opts, scw.WithRegions(api.Regions()...))
				request.Region = ""
			}
			resp, err := api.ListTags(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Tags, nil

		},
	}
}

func registryTagGet() *core.Command {
	return &core.Command{
		Short:     `Get a tag`,
		Long:      `Retrieve information about a given image tag, specified by its ` + "`" + `tag_id` + "`" + ` and region. Full details about the tag, such as ` + "`" + `name` + "`" + `, ` + "`" + `image_id` + "`" + `, ` + "`" + `status` + "`" + `, and ` + "`" + `digest` + "`" + ` are returned in the response.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetTagRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tag-id",
				Short:      `UUID of the tag`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.GetTagRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.GetTag(request)

		},
	}
}

func registryTagDelete() *core.Command {
	return &core.Command{
		Short:     `Delete a tag`,
		Long:      `Delete a given image tag. You must specify, in the endpoint, the ` + "`" + `region` + "`" + ` and ` + "`" + `tag_id` + "`" + ` parameters of the tag you want to delete.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteTagRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tag-id",
				Short:      `UUID of the tag`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "force",
				Short:      `If two tags share the same digest the deletion will fail unless this parameter is set to true (deprecated)`,
				Required:   false,
				Deprecated: true,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.DeleteTagRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			return api.DeleteTag(request)

		},
	}
}
