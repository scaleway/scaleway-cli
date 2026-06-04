package inference

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/inference/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchDeployment struct{}

func (f FetchDeployment) Product() string {
	return "inference"
}

func (f FetchDeployment) Resource() string {
	return "deployment"
}

func (f FetchDeployment) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all inference deployments in a given region.
func (f FetchDeployment) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := inference.NewAPI(client)

	req := &inference.ListDeploymentsRequest{
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
			Locality: deployment.Region.String(),
			ID:       deployment.ID,
			Name:     deployment.Name,
		})
	}

	return results, nil
}
