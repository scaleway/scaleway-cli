package rdb_test

import (
	"os"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_CreateBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
		),
		Cmd:   "scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.ExecAfterCmd("scw rdb backup delete {{ .CmdResult.ID }}"),
			),
			deleteInstance(),
		),
		DefaultRegion: scw.RegionNlAms,
	}))
}

func Test_RestoreBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
		),
		Cmd: "scw rdb backup restore {{ .Backup.ID }} instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.ExecAfterCmd("scw rdb backup delete {{ .Backup.ID }}"),
			),
			deleteInstance(),
		),
	}))
}

func Test_ExportBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			core.BeforeFuncWhenUpdatingCassette(
				func(_ *core.BeforeFuncCtx) error {
					time.Sleep(1 * time.Minute)

					return nil
				},
			),
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
		),
		Cmd: "scw rdb backup export {{ .Backup.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.ExecAfterCmd("scw rdb backup delete {{ .Backup.ID }}"),
			),
			deleteInstance(),
		),
	}))
}

func Test_DownloadBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
			core.ExecStoreBeforeCmd(
				"BackupExport",
				"scw rdb backup export {{ .Backup.ID }} --wait",
			),
		),
		Cmd: "scw rdb backup download {{ .Backup.ID }} output=simple_dump",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.ExecAfterCmd("scw rdb backup delete {{ .Backup.ID }}"),
			),
			deleteInstance(),
			func(_ *core.AfterFuncCtx) error {
				err := os.Remove("simple_dump")

				return err
			},
		),
		DefaultRegion: scw.RegionNlAms,
		TmpHomeDir:    true,
	}))

	t.Run("With no previous export backup", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
		),
		Cmd: "scw rdb backup download {{ .Backup.ID }} output=no_previous_export_dump",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.ExecAfterCmd("scw rdb backup delete {{ .Backup.ID }}"),
			),
			deleteInstance(),
			func(_ *core.AfterFuncCtx) error {
				err := os.Remove("no_previous_export_dump")

				return err
			},
		),
		DefaultRegion: scw.RegionNlAms,
		TmpHomeDir:    true,
	}))
}

// If ran please update the cassette by changing 'download_url_expires_at'
// when it's not null to a much later date
// E.g. from "2022-09-05T13:14:54.437192Z" to "2032-09-05T13:14:54.437192Z"
func Test_ListBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createInstance(engine),
			core.ExecStoreBeforeCmd(
				"BackupA",
				"scw rdb backup create name=will_be_exported expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
			core.ExecStoreBeforeCmd(
				"BackupB",
				"scw rdb backup create name=will_not_be_exported expires-at=2032-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
			core.ExecStoreBeforeCmd(
				"BackupExport",
				"scw rdb backup export {{ .BackupA.ID }} --wait",
			),
		),
		Cmd: "scw rdb backup list instance-id={{ .Instance.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.AfterFuncWhenUpdatingCassette(
				core.AfterFuncCombine(
					core.ExecAfterCmd("scw rdb backup delete {{ .BackupA.ID }}"),
					core.ExecAfterCmd("scw rdb backup delete {{ .BackupB.ID }}"),
				),
			),
			deleteInstance(),
		),
	}))
}
