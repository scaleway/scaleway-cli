package rdb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchInstances struct{}

func (f FetchInstances) Product() string {
	return "rdb"
}

func (f FetchInstances) Resource() string {
	return "instances"
}

func (f FetchInstances) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all RDB instances in a given region.
func (f FetchInstances) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := rdb.NewAPI(client)

	req := &rdb.ListInstancesRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListInstances(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Instances))
	for _, inst := range resp.Instances {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       inst.ID,
			Name:     inst.Name,
		})
	}

	return results, nil
}
