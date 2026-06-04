package searchdb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	searchdb "github.com/scaleway/scaleway-sdk-go/api/searchdb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchDeployments struct{}

func (f FetchDeployments) Product() string {
	return "searchdb"
}

func (f FetchDeployments) Resource() string {
	return "deployment"
}

func (f FetchDeployments) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all SearchDB deployments in a given region.
func (f FetchDeployments) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := searchdb.NewAPI(client)

	req := &searchdb.ListDeploymentsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListDeployments(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Deployments))
	for _, deployment := range resp.Deployments {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       deployment.ID,
			Name:     deployment.Name,
		})
	}

	return results, nil
}
