package container

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchNamespaces struct{}

func (f FetchNamespaces) Product() string {
	return "container"
}

func (f FetchNamespaces) Resource() string {
	return "namespace"
}

func (f FetchNamespaces) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all container namespaces in a given region.
func (f FetchNamespaces) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	req := &container.ListNamespacesRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListNamespaces(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Namespaces))
	for _, ns := range resp.Namespaces {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       ns.ID,
			Name:     ns.Name,
		})
	}

	return results, nil
}
