package rdb

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func backupCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: argsI.(*rdb.WaitForDatabaseBackupRequest).DatabaseBackupID,
			Region:           respI.(*rdb.WaitForDatabaseBackupRequest).Region,
			Timeout:          scw.TimeDurationPtr(backupActionTimeout),
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}
