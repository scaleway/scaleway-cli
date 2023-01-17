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
		marketplaceImageList(),
		marketplaceImageGet(),
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

func marketplaceImageList() *core.Command {
	return &core.Command{
		Short:     `List marketplace images`,
		Long:      `List marketplace images.`,
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
		SeeAlsos: []*core.SeeAlso{
			{
				Command: "scw instance list images",
				Short:   "List all images available in an account",
			},
		},
	}
}

func marketplaceImageGet() *core.Command {
	return &core.Command{
		Short:     `Get a specific marketplace image`,
		Long:      `Get a specific marketplace image.`,
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
