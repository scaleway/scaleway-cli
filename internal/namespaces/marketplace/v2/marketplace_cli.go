// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package marketplace

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		marketplaceRoot(),
		marketplaceImage(),
		marketplaceLocalImage(),
		marketplaceVersion(),
		marketplaceCategory(),
		marketplaceImageList(),
		marketplaceImageGet(),
		marketplaceVersionList(),
		marketplaceVersionGet(),
		marketplaceLocalImageList(),
		marketplaceLocalImageGet(),
		marketplaceCategoryList(),
		marketplaceCategoryGet(),
	)
}
func marketplaceRoot() *core.Command {
	return &core.Command{
		Short:     `Marketplace API`,
		Long:      ``,
		Namespace: "marketplace",
	}
}

func marketplaceImage() *core.Command {
	return &core.Command{
		Short:     `Marketplace images management commands`,
		Long:      `Marketplace images management commands.`,
		Namespace: "marketplace",
		Resource:  "image",
	}
}

func marketplaceLocalImage() *core.Command {
	return &core.Command{
		Short:     `Marketplace local images management commands`,
		Long:      `Marketplace local images management commands.`,
		Namespace: "marketplace",
		Resource:  "local-image",
	}
}

func marketplaceVersion() *core.Command {
	return &core.Command{
		Short:     `Marketplace version management commands`,
		Long:      `Marketplace version management commands.`,
		Namespace: "marketplace",
		Resource:  "version",
	}
}

func marketplaceCategory() *core.Command {
	return &core.Command{
		Short:     `Marketplace category management commands`,
		Long:      `Marketplace category management commands.`,
		Namespace: "marketplace",
		Resource:  "category",
	}
}

func marketplaceImageList() *core.Command {
	return &core.Command{
		Short:     `List marketplace images`,
		Long:      `List all available images on the marketplace, their UUID, CPU architecture and description.`,
		Namespace: "marketplace",
		Resource:  "image",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.ListImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "order-by",
				Short:      `Ordering to use`,
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"name_asc", "name_desc", "created_at_asc", "created_at_desc", "updated_at_asc", "updated_at_desc"},
			},
			{
				Name:       "arch",
				Short:      `Choose for which machine architecture to return images`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "category",
				Short:      `Choose the category of images to get`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "include-eol",
				Short:      `Choose to include end-of-life images`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
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
				FieldName: "Label",
			},
			{
				FieldName: "Name",
			},
			{
				FieldName: "Categories",
			},
			{
				FieldName: "ValidUntil",
			},
			{
				FieldName: "Description",
			},
			{
				FieldName: "UpdatedAt",
			},
			{
				FieldName: "CreatedAt",
			},
			{
				FieldName: "Logo",
			},
		}},
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance image list",
				Short:   "List all images available in an account",
			},
		},
	}
}

func marketplaceImageGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific marketplace image`,
		Long:      `Get detailed information about a marketplace image, specified by its ` + "`" + `image_id` + "`" + ` (UUID format).`,
		Namespace: "marketplace",
		Resource:  "image",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.GetImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Short:      `Display the image name`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.GetImageRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			return api.GetImage(request)

		},
	}
}

func marketplaceVersionList() *core.Command {
	return &core.Command{
		Short:     `List versions of an Image`,
		Long:      `Get a list of all available version of an image, specified by its ` + "`" + `image_id` + "`" + ` (UUID format).`,
		Namespace: "marketplace",
		Resource:  "version",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.ListVersionsRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.ListVersionsRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListVersions(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Versions, nil

		},
	}
}

func marketplaceVersionGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific image version`,
		Long:      `Get information such as the name, creation date, last update and published date for an image version specified by its ` + "`" + `version_id` + "`" + ` (UUID format).`,
		Namespace: "marketplace",
		Resource:  "version",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.GetVersionRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "version-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.GetVersionRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			return api.GetVersion(request)

		},
	}
}

func marketplaceLocalImageList() *core.Command {
	return &core.Command{
		Short:     `List local images from a specific image or version`,
		Long:      `List information about local images in a specific Availability Zone, specified by its ` + "`" + `image_id` + "`" + ` (UUID format), ` + "`" + `version_id` + "`" + ` (UUID format) or ` + "`" + `image_label` + "`" + `. Only one of these three parameters may be set.`,
		Namespace: "marketplace",
		Resource:  "local-image",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.ListLocalImagesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "image-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "version-id",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "order-by",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"created_at_asc", "created_at_desc"},
			},
			{
				Name:       "image-label",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "zone",
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "type",
				Required:   false,
				Deprecated: false,
				Positional: false,
				EnumValues: []string{"unknown_type", "instance_local", "instance_sbs"},
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.ListLocalImagesRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListLocalImages(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.LocalImages, nil

		},
		View: &core.View{Fields: []*core.ViewField{
			{
				FieldName: "ID",
			},
			{
				FieldName: "Label",
			},
			{
				FieldName: "Arch",
			},
			{
				FieldName: "Zone",
			},
			{
				FieldName: "CompatibleCommercialTypes",
			},
		}},
	}
}

func marketplaceLocalImageGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific local image by ID`,
		Long:      `Get detailed information about a local image, including compatible commercial types, supported architecture, labels and the Availability Zone of the image, specified by its ` + "`" + `local_image_id` + "`" + ` (UUID format).`,
		Namespace: "marketplace",
		Resource:  "local-image",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.GetLocalImageRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "local-image-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.GetLocalImageRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			return api.GetLocalImage(request)

		},
	}
}

func marketplaceCategoryList() *core.Command {
	return &core.Command{
		Short:     `List existing image categories`,
		Long:      `Get a list of all existing categories. The output can be paginated.`,
		Namespace: "marketplace",
		Resource:  "category",
		Verb:      "list",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.ListCategoriesRequest{}),
		ArgSpecs: core.ArgSpecs{},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.ListCategoriesRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			opts := []scw.RequestOption{scw.WithAllPages()}
			resp, err := api.ListCategories(request, opts...)
			if err != nil {
				return nil, err
			}
			return resp.Categories, nil

		},
	}
}

func marketplaceCategoryGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific category`,
		Long:      `Get information about a specific category of the marketplace catalog, specified by its ` + "`" + `category_id` + "`" + ` (UUID format).`,
		Namespace: "marketplace",
		Resource:  "category",
		Verb:      "get",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(marketplace.GetCategoryRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "category-id",
				Required:   true,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.GetCategoryRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			return api.GetCategory(request)

		},
	}
}
