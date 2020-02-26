package marketplace

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

func updateCommands(commands *core.Commands) {
	updateMarketplaceListImages(commands.MustFind("marketplace", "image", "list"))
	updateMarketplaceGetImage(commands.MustFind("marketplace", "image", "get"))
}

// updateMarketplaceListImages is the custom Run for scw marketplace list images
// TODO: remove when [APIGW-1959] is done
func updateMarketplaceListImages(c *core.Command) {
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		client := core.ExtractClient(ctx)
		req := argsI.(*marketplace.ListImagesRequest)
		api := marketplace.NewAPI(client)
		req.PerPage = scw.Uint32Ptr(100)
		imagesResponse, err := api.ListImages(req, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		return imagesResponse.Images, nil
	}
}

// TODO : use generated command with label 'param'
// marketplaceGetImagesCustom is the custom command for scw marketplace get image
func updateMarketplaceGetImage(c *core.Command) {
	type getImagesArgs struct {
		Label string
	}

	c.ArgsType = reflect.TypeOf(getImagesArgs{})
	c.ArgSpecs = core.ArgSpecs{
		{
			Name:     "label",
			Short:    ``,
			Required: true,
		},
	}
	c.Run = func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		args := argsI.(*getImagesArgs)
		req := &marketplace.ListImagesRequest{}
		req.PerPage = scw.Uint32Ptr(100)
		client := core.ExtractClient(ctx)
		api := marketplace.NewAPI(client)
		resp, err := api.ListImages(req, scw.WithAllPages())
		if err != nil {
			return nil, err
		}
		return getImageByLabel(resp.Images, args.Label)
	}
}

// getImageByLabel returns a single Image from a slice based on the label of the Image
func getImageByLabel(images []*marketplace.Image, label string) (*marketplace.Image, error) {
	for _, image := range images {
		if strcase.ToBashArg(image.Label) == strcase.ToBashArg(label) {
			return image, nil
		}
	}
	return nil, fmt.Errorf("image not found: no image with label '%v'", label)
}
