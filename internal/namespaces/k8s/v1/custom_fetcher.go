package k8s

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	k8s2 "github.com/scaleway/scaleway-sdk-go/api/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchClusters struct{}

func (f FetchClusters) Product() string {
	return "k8s"
}

func (f FetchClusters) Resource() string {
	return "clusters"
}

func (f FetchClusters) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all Kubernetes clusters in a given region.
func (f FetchClusters) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := k8s2.NewAPI(client)

	req := &k8s2.ListClustersRequest{
		Region: region,
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
			Locality: region.String(),
			ID:       cluster.ID,
			Name:     cluster.Name,
		})
	}

	return results, nil
}
