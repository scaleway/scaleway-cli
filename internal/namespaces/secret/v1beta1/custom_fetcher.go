package secret

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	secret "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchSecrets struct{}

func (f FetchSecrets) Product() string {
	return "secret"
}

func (f FetchSecrets) Resource() string {
	return "secrets"
}

func (f FetchSecrets) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all secrets in a given region.
func (f FetchSecrets) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := secret.NewAPI(client)

	req := &secret.ListSecretsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListSecrets(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Secrets))
	for _, secretItem := range resp.Secrets {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       secretItem.ID,
			Name:     secretItem.Name,
		})
	}

	return results, nil
}
