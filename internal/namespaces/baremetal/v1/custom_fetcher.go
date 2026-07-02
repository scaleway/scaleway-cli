package baremetal

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchServers struct{}

func (s *FetchServers) Resource() string {
	return "server"
}

func (*FetchServers) Product() string {
	return "baremetal"
}

func (s *FetchServers) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeZone
}

// Fetch fetches all baremetal servers in a given zone.
func (*FetchServers) Fetch(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := baremetal.NewAPI(client)

	req := &baremetal.ListServersRequest{
		Zone: zone,
	}
	if projectID != "" {
		req.ProjectID = &projectID
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
