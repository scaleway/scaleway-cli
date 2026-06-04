package block

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/block/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchSnapshots struct{}

func (s *FetchSnapshots) Resource() string {
	return "snapshot"
}

func (*FetchSnapshots) Product() string {
	return "block"
}

func (s *FetchSnapshots) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all block storage snapshots in a given zone.
func (*FetchSnapshots) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := block.NewAPI(client)

	req := &block.ListSnapshotsRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
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

type FetchVolumes struct{}

func (f FetchVolumes) Product() string {
	return "block"
}

func (f FetchVolumes) Resource() string {
	return "volume"
}

func (f FetchVolumes) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all block volumes in a given zone.
func (f FetchVolumes) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := block.NewAPI(client)

	req := &block.ListVolumesRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
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
