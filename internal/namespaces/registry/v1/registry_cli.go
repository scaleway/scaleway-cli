// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package registry

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
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
		Short:     `Container registry API`,
		Long:      ``,
		Namespace: "registry",
	}
}

func registryNamespace() *core.Command {
	return &core.Command{
		Short: `Namespace management commands`,
		Long: `A namespace is for images what a folder is for files

To use our services, the first step is to create a namespace.

A namespace is for images what a folder is for files. Every push or pull must mention the namespace :
` + "`" + `` + "`" + `` + "`" + `docker pull rg.nl-ams.scw.cloud/<namespace_name>/<image_name>:<tag_name>` + "`" + `` + "`" + `` + "`" + `

Note that a namespace name is unique on a region. Thus, if another client already has created "test", you can't have it as a namespace

A namespace can be either public or private (default), which determines who can pull images.
`,
		Namespace: "registry",
		Resource:  "namespace",
	}
}

func registryImage() *core.Command {
	return &core.Command{
		Short: `Image management commands`,
		Long: `An image represents a container image.

The visibility of an image can be public (everyone can pull it), private (only your organization can pull it) or inherit from the namespace visibility (default)
It can be changed with an update on the image via the registry API.
`,
		Namespace: "registry",
		Resource:  "image",
	}
}

func registryTag() *core.Command {
	return &core.Command{
		Short: `Tag management commands`,
		Long: `A tag represents a container tag of an image.
`,
		Namespace: "registry",
		Resource:  "tag",
	}
}

func registryNamespaceList() *core.Command {
	return &core.Command{
		Short:     `List all your namespaces`,
		Long:      `List all your namespaces.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListNamespacesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Field by which to order the display of Images`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListNamespacesRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
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
		Long:      `Get the namespace associated with the given id.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `The unique ID of the Namespace`,
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
		Short:     `Create a new namespace`,
		Long:      `Create a new namespace.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "create",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.CreateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      `Define a namespace name`,
				Required:   true,
				Deprecated: false,
				Positional: false,
				Default:    core.RandomValueGenerator("ns"),
			},
			{
				Name:       "description",
				Short:      `Define a description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.ProjectIDArgSpec(),
			{
				Name:       "is-public",
				Short:      `Define the default visibility policy`,
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
		Short:     `Update an existing namespace`,
		Long:      `Update the namespace associated with the given id.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.UpdateNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `Namespace ID to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "description",
				Short:      `Define a description`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "is-public",
				Short:      `Define the default visibility policy`,
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
		Short:     `Delete an existing namespace`,
		Long:      `Delete the namespace associated with the given id.`,
		Namespace: "registry",
		Resource:  "namespace",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteNamespaceRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace-id",
				Short:      `The unique ID of the Namespace`,
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
		Short:     `List all your images`,
		Long:      `List all your images.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Field by which to order the display of Images`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "namespace-id",
				Short:      `Filter by the Namespace ID`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by the Image name (exact match)`,
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
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			resp, err := api.ListImages(request, scw.WithAllPages())
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
		Short:     `Get a image`,
		Long:      `Get the image associated with the given id.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `The unique ID of the Image`,
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
		Short:     `Update an existing image`,
		Long:      `Update the image associated with the given id.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "update",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.UpdateImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `Image ID to update`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "visibility",
				Short:      `A ` + "`" + `public` + "`" + ` image is pullable from internet without authentication, opposed to a ` + "`" + `private` + "`" + ` image. ` + "`" + `inherit` + "`" + ` will use the namespace ` + "`" + `is_public` + "`" + ` parameter`,
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
		Long:      `Delete the image associated with the given id.`,
		Namespace: "registry",
		Resource:  "image",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `The unique ID of the Image`,
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
		Short:     `List all your tags`,
		Long:      `List all your tags.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.ListTagsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Field by which to order the display of Images`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc", "name_asc", "name_desc"},
			},
			{
				Name:       "image-id",
				Short:      `The unique ID of the image`,
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "name",
				Short:      `Filter by the tag name (exact match)`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms, scw.RegionPlWaw),
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*registry.ListTagsRequest)

			client := core.ExtractClient(ctx)
			api := registry.NewAPI(client)
			resp, err := api.ListTags(request, scw.WithAllPages())
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
		Long:      `Get the tag associated with the given id.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.GetTagRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tag-id",
				Short:      `The unique ID of the Tag`,
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
		Long:      `Delete the tag associated with the given id.`,
		Namespace: "registry",
		Resource:  "tag",
		Verb:      "delete",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(registry.DeleteTagRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "tag-id",
				Short:      `The unique ID of the tag`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			{
				Name:       "force",
				Short:      `If two tags share the same digest the deletion will fail unless this parameter is set to true`,
				Required:   false,
				Deprecated: false,
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
