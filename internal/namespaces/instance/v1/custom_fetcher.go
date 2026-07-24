package instance

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchServers struct{}

func (f FetchServers) Product() string {
	return "instance"
}

func (f FetchServers) Resource() string {
	return "server"
}

func (f FetchServers) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all instances in a given zone.
func (f FetchServers) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)

	req := &instance.ListServersRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.Project = &projectID
	}

	resp, err := api.ListServers(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Servers))
	for _, server := range resp.Servers {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       server.ID,
			Name:     server.Name,
		})
	}

	return results, nil
}

type FetchIPs struct{}

func (f FetchIPs) Product() string {
	return "instance"
}

func (f FetchIPs) Resource() string {
	return "ip"
}

func (f FetchIPs) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all instance IPs in a given zone.
func (f FetchIPs) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)

	req := &instance.ListIPsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.Project = &projectID
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

type FetchVolumes struct{}

func (f FetchVolumes) Product() string {
	return "instance"
}

func (f FetchVolumes) Resource() string {
	return "volume"
}

func (f FetchVolumes) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all instance volumes in a given zone.
func (f FetchVolumes) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)

	req := &instance.ListVolumesRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.Project = &projectID
	}

	resp, err := api.ListVolumes(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Volumes))
	for _, vol := range resp.Volumes {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       vol.ID,
			Name:     vol.Name,
		})
	}

	return results, nil
}

type FetchSnapshots struct{}

func (f FetchSnapshots) Product() string {
	return "instance"
}

func (f FetchSnapshots) Resource() string {
	return "snapshot"
}

func (f FetchSnapshots) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all instance snapshots in a given zone.
func (f FetchSnapshots) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)

	req := &instance.ListSnapshotsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.Project = &projectID
	}

	resp, err := api.ListSnapshots(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Snapshots))
	for _, snap := range resp.Snapshots {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       snap.ID,
			Name:     snap.Name,
		})
	}

	return results, nil
}

type FetchSecurityGroups struct{}

func (f FetchSecurityGroups) Product() string {
	return "instance"
}

func (f FetchSecurityGroups) Resource() string {
	return "security-group"
}

func (f FetchSecurityGroups) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all instance security groups in a given zone.
func (f FetchSecurityGroups) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := instance.NewAPI(client)

	req := &instance.ListSecurityGroupsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.Project = &projectID
	}

	resp, err := api.ListSecurityGroups(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.SecurityGroups))
	for _, sg := range resp.SecurityGroups {
		results = append(results, fetch.ResourceResult{
			Locality: zone.String(),
			ID:       sg.ID,
			Name:     sg.Name,
		})
	}

	return results, nil
}
