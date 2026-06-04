package ipam

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/ipam/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchIPs struct{}

func (f FetchIPs) Product() string {
	return "ipam"
}

func (f FetchIPs) Resource() string {
	return "ip"
}

func (f FetchIPs) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all IPAM IPs in a given region.
func (f FetchIPs) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := ipam.NewAPI(client)

	req := &ipam.ListIPsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListIPs(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.IPs))
	for _, ip := range resp.IPs {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       ip.ID,
			Name:     ip.Address.String(),
		})
	}

	return results, nil
}
