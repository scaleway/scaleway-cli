package rdb

import (
	"fmt"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_CreateBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Instance",
				fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
			),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			func(ctx *core.BeforeFuncCtx) error {
				time.Sleep(1 * time.Minute)
				return nil
			}),
		Cmd:   "scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}"),
			core.ExecAfterCmd("scw rdb backup delete {{ .CmdResult.ID }}"),
		),
		DefaultRegion: scw.RegionNlAms,
	}))
}

func Test_RestoreBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Instance",
				fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
			),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			func(ctx *core.BeforeFuncCtx) error {
				time.Sleep(1 * time.Minute)
				return nil
			},
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
		),
		Cmd: "scw rdb backup restore {{ .Backup.ID }} instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}"),
			core.ExecAfterCmd("scw rdb backup delete {{ .CmdResult.ID }}"),
		),
		DefaultRegion: scw.RegionNlAms,
	}))
}

func Test_ExportBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Instance",
				fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
			),
			// We opened an internal issue about the fact that the instance is considered ready even if rdb is not yet available.
			func(ctx *core.BeforeFuncCtx) error {
				time.Sleep(1 * time.Minute)
				return nil
			},
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-07:00 instance-id={{ .Instance.ID }} database-name=rdb --wait",
			),
		),
		Cmd: "scw rdb backup export {{ .Backup.ID }} --wait",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}"),
			core.ExecAfterCmd("scw rdb backup delete {{ .CmdResult.ID }}"),
		),
		DefaultRegion: scw.RegionNlAms,
	}))
}
