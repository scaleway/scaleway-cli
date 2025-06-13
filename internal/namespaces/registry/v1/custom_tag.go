package registry

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

//
// Marshalers
//

// tagStatusMarshalerFunc marshals a registry.TagStatus.
var (
	tagStatusMarshalSpecs = human.EnumMarshalSpecs{
		registry.TagStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		registry.TagStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.TagStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.TagStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

type customTag struct {
	registry.Tag
	FullName string
}

func tagGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		getTagResp, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		tag := getTagResp.(*registry.Tag)

		client := core.ExtractClient(ctx)
		api := registry.NewAPI(client)

		image, err := api.GetImage(&registry.GetImageRequest{
			ImageID: tag.ImageID,
		})
		if err != nil {
			logger.Warningf("cannot get image %s %s", tag.ImageID, err)

			return getTagResp, nil
		}

		namespace, err := api.GetNamespace(&registry.GetNamespaceRequest{
			NamespaceID: image.NamespaceID,
		})
		if err != nil {
			logger.Warningf("cannot get namespace %s %s", image.NamespaceID, err)

			return getTagResp, nil
		}

		res := customTag{
			Tag:      *tag,
			FullName: fmt.Sprintf("%s/%s:%s", namespace.Endpoint, image.Name, tag.Name),
		}

		return res, nil
	}

	return c
}

func tagListBuilder(c *core.Command) *core.Command {
	c.View = &core.View{
		Fields: []*core.ViewField{
			{
				Label:     "ID",
				FieldName: "ID",
			},
			{
				Label:     "Full Name",
				FieldName: "FullName",
			},
			{
				Label:     "Status",
				FieldName: "Status",
			},
		},
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		listTagResp, err := runner(ctx, argsI)
		if err != nil {
			return listTagResp, err
		}

		client := core.ExtractClient(ctx)
		api := registry.NewAPI(client)

		request := argsI.(*registry.ListTagsRequest)
		image, err := api.GetImage(&registry.GetImageRequest{
			ImageID: request.ImageID,
		})
		if err != nil {
			return listTagResp, err
		}

		namespace, err := api.GetNamespace(&registry.GetNamespaceRequest{
			NamespaceID: image.NamespaceID,
		})
		if err != nil {
			return listTagResp, err
		}

		var customRes []customTag
		for _, tag := range listTagResp.([]*registry.Tag) {
			customRes = append(customRes, customTag{
				Tag: *tag,
				FullName: fmt.Sprintf("%s/%s:%s",
					namespace.Endpoint,
					image.Name,
					tag.Name,
				),
			})
		}

		return customRes, nil
	}

	return c
}
