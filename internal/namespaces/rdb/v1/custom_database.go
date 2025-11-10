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
		if req.Name == "" {
			return runner(ctx, argsI)
		}

		api := rdb.NewAPI(core.ExtractClient(ctx))
		name := req.Name
		list, err := api.ListDatabases(&rdb.ListDatabasesRequest{
			Region:     req.Region,
			InstanceID: req.InstanceID,
			Name:       &name,
		}, scw.WithAllPages())

		if err == nil && list.TotalCount > 0 {
			return list.Databases[0], nil
		}

		return runner(ctx, argsI)
	}

	return c
}
