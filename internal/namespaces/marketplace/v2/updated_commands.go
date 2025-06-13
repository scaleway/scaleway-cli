package marketplace

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/marketplace/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

func updateCommands(commands *core.Commands) {
	updateMarketplaceGetImage(commands.MustFind("marketplace", "image", "get"))
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
	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		args := argsI.(*getImagesArgs)
		req := &marketplace.ListImagesRequest{}
		req.PageSize = scw.Uint32Ptr(100)
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
