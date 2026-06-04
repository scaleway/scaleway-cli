package cockpit

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchToken struct{}

func (f FetchToken) Product() string {
	return "cockpit"
}

func (f FetchToken) Resource() string {
	return "token"
}

func (f FetchToken) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

func (f FetchToken) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := cockpit.NewRegionalAPI(client)

	req := &cockpit.RegionalAPIListTokensRequest{
		Region:    region,
		ProjectID: projectID,
	}

	resp, err := api.ListTokens(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Tokens))
	for _, token := range resp.Tokens {
		results = append(results, fetch.ResourceResult{
			Locality: token.Region.String(),
			ID:       token.ID,
			Name:     token.Name,
		})
	}

	return results, nil
}

type FetchDataSource struct{}

func (f FetchDataSource) Product() string {
	return "cockpit"
}

func (f FetchDataSource) Resource() string {
	return "data-source"
}

func (f FetchDataSource) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

func (f FetchDataSource) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := cockpit.NewRegionalAPI(client)

	req := &cockpit.RegionalAPIListDataSourcesRequest{
		Region:    region,
		Origin:    cockpit.DataSourceOriginCustom,
		ProjectID: projectID,
	}

	resp, err := api.ListDataSources(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.DataSources))
	for _, dataSource := range resp.DataSources {
		results = append(results, fetch.ResourceResult{
			Locality: dataSource.Region.String(),
			ID:       dataSource.ID,
			Name:     dataSource.Name,
		})
	}

	return results, nil
}
