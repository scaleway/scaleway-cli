package flexibleip

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	flexibleip2 "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchFlexibleIPs struct{}

func (f FetchFlexibleIPs) Product() string {
	return "fip"
}

func (f FetchFlexibleIPs) Resource() string {
	return "ip"
}

func (f FetchFlexibleIPs) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all flexible IPs in a given zone.
func (f FetchFlexibleIPs) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := flexibleip2.NewAPI(client)

	req := &flexibleip2.ListFlexibleIPsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListFlexibleIPs(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.FlexibleIPs))
	for _, ip := range resp.FlexibleIPs {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       ip.ID,
			Name:     ip.Description,
		})
	}

	return results, nil
}
