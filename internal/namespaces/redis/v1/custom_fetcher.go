package redis

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/redis/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchClusters struct{}

func (f FetchClusters) Product() string {
	return "redis"
}

func (f FetchClusters) Resource() string {
	return "clusters"
}

func (f FetchClusters) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all Redis clusters in a given zone.
func (f FetchClusters) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := redis.NewAPI(client)

	req := &redis.ListClustersRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListClusters(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Clusters))
	for _, cluster := range resp.Clusters {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       cluster.ID,
			Name:     cluster.Name,
		})
	}

	return results, nil
}
