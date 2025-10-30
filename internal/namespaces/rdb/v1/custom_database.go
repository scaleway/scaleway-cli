package rdb

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func databaseCreateBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		req := argsI.(*rdb.CreateDatabaseRequest)
		api := rdb.NewAPI(core.ExtractClient(ctx))
		if req != nil && req.Name != "" {
			name := req.Name
			list, err := api.ListDatabases(&rdb.ListDatabasesRequest{
				Region:     req.Region,
				InstanceID: req.InstanceID,
				Name:       &name,
			}, scw.WithAllPages())
			if err == nil && list.TotalCount > 0 {
				// Return the existing database to maintain compatibility with workflows
				// that expect database object (not just SuccessResult)
				return list.Databases[0], nil
			}
		}

		return runner(ctx, argsI)
	}

	return c
}
