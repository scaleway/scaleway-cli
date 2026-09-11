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
