package webhosting

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchHostings struct{}

func (f FetchHostings) Product() string {
	return "webhosting"
}

func (f FetchHostings) Resource() string {
	return "hosting"
}

func (f FetchHostings) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all Web Hosting plans in a given region.
func (f FetchHostings) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := webhosting.NewHostingAPI(client)

	req := &webhosting.HostingAPIListHostingsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListHostings(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Hostings))
	for _, hosting := range resp.Hostings {
		result := fetch.ResourceResult{
			Locality: region.String(),
			Product:  "webhosting",
			Resource: "hosting",
			ID:       hosting.ID,
		}
		if hosting.Domain != nil {
			result.Name = *hosting.Domain
		}
		results = append(results, result)
	}

	return results, nil
}
