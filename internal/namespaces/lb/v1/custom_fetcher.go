package lb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchLoadBalancers struct{}

func (f FetchLoadBalancers) Product() string {
	return "lb"
}

func (f FetchLoadBalancers) Resource() string {
	return "lb"
}

func (f FetchLoadBalancers) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all load balancers in a given region.
func (f FetchLoadBalancers) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := lb.NewAPI(client)

	req := &lb.ListLBsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListLBs(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.LBs))
	for _, lbItem := range resp.LBs {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       lbItem.ID,
			Name:     lbItem.Name,
		})
	}

	return results, nil
}
