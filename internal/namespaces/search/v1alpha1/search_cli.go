// This file was automatically generated. DO NOT EDIT.
// If you have any remark or suggestion do not hesitate to open an issue.

package search

import (
	"context"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	search "github.com/scaleway/scaleway-sdk-go/api/search/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// always import dependencies
var (
	_ = scw.RegionFrPar
)

func GetGeneratedCommands() *core.Commands {
	return core.NewCommands(
		searchRoot(),
		searchResource(),
		searchResourceSearch(),
	)
}

func searchRoot() *core.Command {
	return &core.Command{
		Short:     `Search API`,
		Long:      ``,
		Namespace: "search",
	}
}

func searchResource() *core.Command {
	return &core.Command{
		Short:     `Resource search commands`,
		Long:      `Resource search commands.`,
		Namespace: "search",
		Resource:  "resource",
	}
}

func searchResourceSearch() *core.Command {
	return &core.Command{
		Short:     `Search API`,
		Long:      `Search API.`,
		Namespace: "search",
		Resource:  "resource",
		Verb:      "search",
		// Deprecated:    false,
		ArgsType: reflect.TypeOf(search.SearchResourcesRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "query",
				Short:      `Search query`,
				Required:   true,
				Deprecated: false,
				Positional: true,
			},
			core.OrganizationIDArgSpec(),
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			request := args.(*search.SearchResourcesRequest)

			client := core.ExtractClient(ctx)
			api := search.NewAPI(client)

			return api.SearchResources(request)
		},
	}
}
