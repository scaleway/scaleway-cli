// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package marketplace

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v1"
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
		ArgSpecs: core.ArgSpecs{},
		Run: func(ctx context.Context, args interface{}) (i interface{}, e error) {
			request := args.(*marketplace.ListImagesRequest)

			client := core.ExtractClient(ctx)
			api := marketplace.NewAPI(client)
			resp, err := api.ListImages(request, scw.WithAllPages())
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
