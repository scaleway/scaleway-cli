package registry

import (
	"context"
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

// imageStatusMarshalerFunc marshals a registry.ImageStatus.
var (
	imageStatusMarshalSpecs = human.EnumMarshalSpecs{
		registry.ImageStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		registry.ImageStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.ImageStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.ImageStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

type CustomImage struct {
	registry.Image
	FullName           string
	ExplicitVisibility string `json:"-"`
}

func imageGetBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		getImageResp, err := runner(ctx, argsI)
		if err != nil {
			return nil, err
		}
		image := getImageResp.(*registry.Image)

		client := core.ExtractClient(ctx)
		api := registry.NewAPI(client)

		namespace, err := api.GetNamespace(&registry.GetNamespaceRequest{
			NamespaceID: image.NamespaceID,
		})
		if err != nil {
			return getImageResp, err
		}

		res := CustomImage{
			Image:    *image,
			FullName: fmt.Sprintf("%s/%s", namespace.Endpoint, image.Name),
		}

		return res, nil
	}

	return c
}

func imageListBuilder(c *core.Command) *core.Command {
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
				Label:     "Size",
				FieldName: "Size",
			},
			{
				Label:     "Visibility",
				FieldName: "ExplicitVisibility",
			},
			{
				Label:     "Status",
				FieldName: "Status",
			},
			{
				Label:     "Created At",
				FieldName: "CreatedAt",
			},
			{
				Label:     "Updated At",
				FieldName: "UpdatedAt",
			},
		},
	}

	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		listImageResp, err := runner(ctx, argsI)
		if err != nil {
			return listImageResp, err
		}
		listImage := listImageResp.([]*registry.Image)

		client := core.ExtractClient(ctx)
		api := registry.NewAPI(client)

		namespaces, err := api.ListNamespaces(&registry.ListNamespacesRequest{}, scw.WithAllPages())
		if err != nil {
			return listImageResp, err
		}

		namespaceEndpointByID := make(map[string]string)
		namespaceVisibilityByID := make(map[string]registry.ImageVisibility)
		for _, namespace := range namespaces.Namespaces {
			namespaceEndpointByID[namespace.ID] = namespace.Endpoint
			if namespace.IsPublic {
				namespaceVisibilityByID[namespace.ID] = registry.ImageVisibilityPublic
			} else {
				namespaceVisibilityByID[namespace.ID] = registry.ImageVisibilityPrivate
			}
		}

		var customRes []CustomImage
		for _, image := range listImage {
			img := CustomImage{
				Image: *image,
				FullName: fmt.Sprintf(
					"%s/%s",
					namespaceEndpointByID[image.NamespaceID],
					image.Name,
				),
			}

			if image.Visibility == registry.ImageVisibilityInherit {
				img.ExplicitVisibility = fmt.Sprintf(
					"%s (inherit)",
					namespaceVisibilityByID[image.NamespaceID],
				)
			} else {
				img.ExplicitVisibility = image.Visibility.String()
			}

			customRes = append(customRes, img)
		}

		return customRes, nil
	}

	return c
}
