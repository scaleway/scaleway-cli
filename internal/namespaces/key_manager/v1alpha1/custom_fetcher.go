package key_manager

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	key_manager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchKey struct{}

func (f FetchKey) Product() string {
	return "keymanager"
}

func (f FetchKey) Resource() string {
	return "key"
}

func (f FetchKey) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all key manager keys in a given region.
func (f FetchKey) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := key_manager.NewAPI(client)

	req := &key_manager.ListKeysRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListKeys(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Keys))
	for _, key := range resp.Keys {
		results = append(results, fetch.ResourceResult{
			Locality: key.Region.String(),
			ID:       key.ID,
			Name:     key.Name,
		})
	}

	return results, nil
}
