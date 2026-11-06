package cockpit

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/cockpit/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchTokens struct{}

func (f *FetchTokens) Namespace() string {
	return cockpitToken().Namespace
}

func (f *FetchTokens) Resource() string {
	return cockpitToken().Resource
}

func (f *FetchTokens) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

func (f *FetchTokens) Fetch(
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

type FetchDataSources struct{}

func (f *FetchDataSources) Namespace() string {
	return cockpitDataSource().Namespace
}

func (f *FetchDataSources) Resource() string {
	return cockpitDataSource().Resource
}

func (f *FetchDataSources) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

func (f *FetchDataSources) Fetch(
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
