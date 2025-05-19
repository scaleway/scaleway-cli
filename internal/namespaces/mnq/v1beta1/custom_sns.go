package mnq

import (
	"context"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

func mnqSnsListCredentialsBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("project-id").AutoCompleteFunc = autocompleteSnsProjectID

	return c
}

func mnqSnsGetCredentialsBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("sns-credentials-id").AutoCompleteFunc = autocompleteSnsCredentialsID

	return c
}

var completeSnsInfoCache *mnq.SnsInfo

func autocompleteSnsProjectID(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	_ = prefix
	req := request.(*mnq.SnsAPIListSnsCredentialsRequest)
	suggestions := core.AutocompleteSuggestions(nil)

	client := core.ExtractClient(ctx)
	api := mnq.NewSnsAPI(client)

	if completeSnsInfoCache == nil {
		res, err := api.GetSnsInfo(&mnq.SnsAPIGetSnsInfoRequest{
			Region: req.Region,
		})
		if err != nil {
			return nil
		}
		completeSnsInfoCache = res
	}

	suggestions = append(suggestions, completeSnsInfoCache.ProjectID)

	return suggestions
}

var completeSnsListCredentials *mnq.ListSnsCredentialsResponse

func autocompleteSnsCredentialsID(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	req := request.(*mnq.SnsAPIGetSnsCredentialsRequest)

	suggestions := core.AutocompleteSuggestions(nil)

	if req.Region != "" {
		client := core.ExtractClient(ctx)
		api := mnq.NewSnsAPI(client)
		if completeSnsInfoCache == nil {
			res, err := api.GetSnsInfo(&mnq.SnsAPIGetSnsInfoRequest{
				Region: req.Region,
			})
			if err != nil {
				return nil
			}
			completeSnsInfoCache = res
		}
		if completeSnsListCredentials == nil {
			res, err := api.ListSnsCredentials(&mnq.SnsAPIListSnsCredentialsRequest{
				Region:    req.Region,
				ProjectID: &completeSnsInfoCache.ProjectID,
			})
			if err != nil {
				return nil
			}
			completeSnsListCredentials = res
		}

		for _, snsCredentials := range completeSnsListCredentials.SnsCredentials {
			if strings.HasPrefix(snsCredentials.ID, prefix) {
				suggestions = append(suggestions, snsCredentials.ID)
			}
		}
	}

	return suggestions
}
