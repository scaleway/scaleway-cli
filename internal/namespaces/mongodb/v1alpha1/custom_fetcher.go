package mongodb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchInstances struct{}

func (f FetchInstances) Product() string {
	return "mongodb"
}

func (f FetchInstances) Resource() string {
	return "instance"
}

func (f FetchInstances) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all MongoDB instances in a given region.
func (f FetchInstances) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := mongodb.NewAPI(client)

	req := &mongodb.ListInstancesRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListInstances(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Instances))
	for _, inst := range resp.Instances {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       inst.ID,
			Name:     inst.Name,
		})
	}

	return results, nil
}

type FetchSnapshot struct{}

func (f FetchSnapshot) Product() string {
	return "mongodb"
}

func (f FetchSnapshot) Resource() string {
	return "snapshot"
}

func (f FetchSnapshot) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all MongoDB snapshots in a given region.
func (f FetchSnapshot) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := mongodb.NewAPI(client)

	req := &mongodb.ListSnapshotsRequest{
		Region: region,
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
	for _, snapshot := range resp.Snapshots {
		results = append(results, fetch.ResourceResult{
			Locality: region.String(),
			ID:       snapshot.ID,
			Name:     snapshot.Name,
		})
	}

	return results, nil
}
