package s2s_vpn

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	s2s_vpn "github.com/scaleway/scaleway-sdk-go/api/s2s_vpn/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchVpnGateways struct{}

func (f FetchVpnGateways) Product() string {
	return "s2s-vpn"
}

func (f FetchVpnGateways) Resource() string {
	return "vpn-gateway"
}

func (f FetchVpnGateways) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all VPN gateways in a given region.
func (f FetchVpnGateways) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := s2s_vpn.NewAPI(client)

	req := &s2s_vpn.ListVpnGatewaysRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListVpnGateways(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Gateways))
	for _, gateway := range resp.Gateways {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       gateway.ID,
			Name:     gateway.Name,
		})
	}

	return results, nil
}
