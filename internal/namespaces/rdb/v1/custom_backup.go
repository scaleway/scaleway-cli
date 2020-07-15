package rdb

import (
	"context"
	"reflect"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	backupActionTimeout = 20 * time.Minute
)

var (
	backupStatusMarshalSpecs = human.EnumMarshalSpecs{
		rdb.DatabaseBackupStatusUnknown:   &human.EnumMarshalSpec{Attribute: color.Faint, Value: "unknown"},
		rdb.DatabaseBackupStatusCreating:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "creating"},
		rdb.DatabaseBackupStatusReady:     &human.EnumMarshalSpec{Attribute: color.FgGreen, Value: "ready"},
		rdb.DatabaseBackupStatusRestoring: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "restoring"},
		rdb.DatabaseBackupStatusDeleting:  &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "deleting"},
		rdb.DatabaseBackupStatusError:     &human.EnumMarshalSpec{Attribute: color.FgRed, Value: "error"},
		rdb.DatabaseBackupStatusExporting: &human.EnumMarshalSpec{Attribute: color.FgBlue, Value: "exporting"},
	}
)

type backupWaitRequest struct {
	DatabaseBackupID string
	Region           scw.Region
}

func backupWaitCommand() *core.Command {
	return &core.Command{
		Short:     `Wait for a backup to reach a stable state`,
		Long:      `Wait for a backup to reach a stable state. This is similar to using --wait flag.`,
		Namespace: "rdb",
		Resource:  "backup",
		Verb:      "wait",
		ArgsType:  reflect.TypeOf(backupWaitRequest{}),
		Run: func(ctx context.Context, argsI interface{}) (i interface{}, err error) {
			api := rdb.NewAPI(core.ExtractClient(ctx))
			return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
				DatabaseBackupID: argsI.(*backupWaitRequest).DatabaseBackupID,
				Region:           argsI.(*backupWaitRequest).Region,
				Timeout:          scw.TimeDurationPtr(backupActionTimeout),
				RetryInterval:    core.DefaultRetryInterval,
			})
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "backup-id",
				Short:      `ID of the backup you want to wait for.`,
				Required:   true,
				Positional: true,
			},
			core.RegionArgSpec(scw.RegionFrPar, scw.RegionNlAms),
		},
		Examples: []*core.Example{
			{
				Short:    "Wait for a backup to reach a stable state",
				ArgsJSON: `{"backup_id": "11111111-1111-1111-1111-111111111111"}`,
			},
		},
	}
}

func backupCreateBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func backupExportBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}

func backupRestoreBuilder(c *core.Command) *core.Command {
	timeout := backupActionTimeout
	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := rdb.NewAPI(core.ExtractClient(ctx))
		return api.WaitForDatabaseBackup(&rdb.WaitForDatabaseBackupRequest{
			DatabaseBackupID: respI.(*rdb.DatabaseBackup).ID,
			Region:           respI.(*rdb.DatabaseBackup).Region,
			Timeout:          &timeout,
			RetryInterval:    core.DefaultRetryInterval,
		})
	}

	return c
}
