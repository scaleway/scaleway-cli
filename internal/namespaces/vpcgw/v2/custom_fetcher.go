package vpcgw

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/vpcgw/v2"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchGateways struct{}

func (f FetchGateways) Product() string {
	return "vpc-gw"
}

func (f FetchGateways) Resource() string {
	return "gateway"
}

func (f FetchGateways) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all VPC gateways in a given zone.
func (f FetchGateways) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := vpcgw.NewAPI(client)

	req := &vpcgw.ListGatewaysRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListGateways(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Gateways))
	for _, gateway := range resp.Gateways {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       gateway.ID,
			Name:     gateway.Name,
		})
	}

	return results, nil
}

type FetchIPs struct{}

func (f FetchIPs) Product() string {
	return "vpc-gw"
}

func (f FetchIPs) Resource() string {
	return "ip"
}

func (f FetchIPs) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all VPC gateway IPs in a given zone.
func (f FetchIPs) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := vpcgw.NewAPI(client)

	req := &vpcgw.ListIPsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListIPs(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.IPs))
	for _, ip := range resp.IPs {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       ip.ID,
			Name:     ip.Address.String(),
		})
	}

	return results, nil
}
