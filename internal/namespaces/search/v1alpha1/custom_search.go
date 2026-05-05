package search

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	search "github.com/scaleway/scaleway-sdk-go/api/search/v1alpha1"
)

func searchBuilder(c *core.Command) *core.Command {
	c.AddInterceptors(
		func(ctx context.Context, argsI any, runner core.CommandRunner) (i any, err error) {
			args := argsI.(*search.SearchResourcesRequest)

			request := &search.SearchResourcesRequest{
				Query: args.Query,
			}
			results, err := runner(ctx, request)
			if err != nil {
				return nil, err
			}

			response := results.(*search.SearchResourcesResponse)

			return response.Resources, nil
		},
	)

	return c
}

//
//Resources.0.ID              4f59cb01-97ec-400a-88c7-92935d36c8fa
//Resources.0.Name            scw-nostalgic-golick-system
//Resources.0.ProjectID       951df375-e094-4d26-97c1-ba548eeb9c42
//Resources.0.OrganizationID  951df375-e094-4d26-97c1-ba548eeb9c42
//Resources.0.Type            sbs_volume
//Resources.0.Zone            fr-par-2
//Resources.1.ID              a31ba2c0-878e-43aa-97c6-9cbf1537c209
//Resources.1.Name            scw-nostalgic-golick
//Resources.1.ProjectID       951df375-e094-4d26-97c1-ba548eeb9c42
//Resources.1.OrganizationID  951df375-e094-4d26-97c1-ba548eeb9c42
//Resources.1.Type            instance_server
//Resources.1.Zone            fr-par-2
