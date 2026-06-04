package registry

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchNamespaces struct{}

func (f FetchNamespaces) Product() string {
	return "registry"
}

func (f FetchNamespaces) Resource() string {
	return "namespace"
}

func (f FetchNamespaces) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all registry namespaces in a given region.
func (f FetchNamespaces) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := registry.NewAPI(client)

	req := &registry.ListNamespacesRequest{
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
	for _, namespace := range resp.Namespaces {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       namespace.ID,
			Name:     namespace.Name,
		})
	}

	return results, nil
}
