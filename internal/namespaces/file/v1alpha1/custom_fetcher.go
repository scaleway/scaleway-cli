package file

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	file "github.com/scaleway/scaleway-sdk-go/api/file/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type FetchFileSystem struct{}

func (f FetchFileSystem) Product() string {
	return "file"
}

func (f FetchFileSystem) Resource() string {
	return "filesystem"
}

func (f FetchFileSystem) LocalityType() fetch.LocalityType {
	return fetch.LocalityTypeRegion
}

// Fetch fetches all file systems in a given region.
func (f FetchFileSystem) Fetch(
	ctx context.Context,
	region scw.Region,
	projectID string,
) ([]fetch.ResourceResult, error) {
	client := core.ExtractClient(ctx)
	api := file.NewAPI(client)

	req := &file.ListFileSystemsRequest{
		Region: region,
	}
	if projectID != "" {
		req.ProjectID = &projectID
	}

	resp, err := api.ListFileSystems(req, scw.WithAllPages(), scw.WithContext(ctx))
	if err != nil {
		if fetch.ShouldIgnoreError(err) {
			return nil, nil
		}

		return nil, err
	}

	results := make([]fetch.ResourceResult, 0, len(resp.Filesystems))
	for _, fs := range resp.Filesystems {
		results = append(results, fetch.ResourceResult{
			Locality: fs.Region.String(),
			ID:       fs.ID,
			Name:     fs.Name,
		})
	}

	return results, nil
}
