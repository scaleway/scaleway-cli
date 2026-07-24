package vpc

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchVPCs struct{}

func (f FetchVPCs) Product() string {
	return "vpc"
}

func (f FetchVPCs) Resource() string {
	return "vpc"
}

func (f FetchVPCs) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all VPCs in a given region.
func (f FetchVPCs) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := vpc.NewAPI(client)

	req := &vpc.ListVPCsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListVPCs(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Vpcs))
	for _, v := range resp.Vpcs {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       v.ID,
			Name:     v.Name,
		})
	}

	return results, nil
}
