package mnq

import (
	"context"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

var mnqSqsInfoStatusMarshalSpecs = human.EnumMarshalSpecs{
	mnq.SqsInfoStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
	mnq.SqsInfoStatusEnabled:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
	mnq.SqsInfoStatusDisabled:      &human.EnumMarshalSpec{Attribute: color.FgRed},
}

func mnqSqsListCredentialsBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("project-id").AutoCompleteFunc = autocompleteSqsProjectID

	return c
}

func mnqSqsGetCredentialsBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("sqs-credentials-id").AutoCompleteFunc = autocompleteSqsCredentialsID

	return c
}

var completeSqsInfoCache *mnq.SqsInfo

func autocompleteSqsProjectID(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	_ = prefix
	req := request.(*mnq.SqsAPIListSqsCredentialsRequest)
	suggestions := core.AutocompleteSuggestions(nil)
	client := core.ExtractClient(ctx)
	api := mnq.NewSqsAPI(client)
	if completeSqsInfoCache == nil {
		res, err := api.GetSqsInfo(&mnq.SqsAPIGetSqsInfoRequest{
			Region: req.Region,
		})
		if err != nil {
			return nil
		}
		completeSqsInfoCache = res
	}

	suggestions = append(suggestions, completeSqsInfoCache.ProjectID)

	return suggestions
}

var completeSqsListCredentials *mnq.ListSqsCredentialsResponse

func autocompleteSqsCredentialsID(
	ctx context.Context,
	prefix string,
	request any,
) core.AutocompleteSuggestions {
	req := request.(*mnq.SqsAPIGetSqsCredentialsRequest)

	suggestions := core.AutocompleteSuggestions(nil)

	if req.Region != "" {
		client := core.ExtractClient(ctx)
		api := mnq.NewSqsAPI(client)
		if completeSqsInfoCache == nil {
			res, err := api.GetSqsInfo(&mnq.SqsAPIGetSqsInfoRequest{
				Region: req.Region,
			})
			if err != nil {
				return nil
			}
			completeSqsInfoCache = res
		}
		if completeSqsListCredentials != nil {
			res, err := api.ListSqsCredentials(&mnq.SqsAPIListSqsCredentialsRequest{
				Region:    req.Region,
				ProjectID: &completeSqsInfoCache.ProjectID,
			})
			if err != nil {
				return nil
			}
			completeSqsListCredentials = res
		}

		for _, sqsCredentials := range completeSqsListCredentials.SqsCredentials {
			if strings.HasPrefix(sqsCredentials.ID, prefix) {
				suggestions = append(suggestions, sqsCredentials.ID)
			}
		}
	}

	return suggestions
}
