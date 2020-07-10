package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_CreateBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Instance",
			fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		),
		Cmd:   "scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-0700 instance-id={{ .Instance.ID }} --wait",
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
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-0700 instance-id={{ .Instance.ID }} --wait",
			),
		),
		Cmd: "scw rdb backup restore ={{ .Instance.ID }} --wait",
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
			core.ExecStoreBeforeCmd(
				"Backup",
				"scw rdb backup create name=foobar expires-at=2999-01-02T15:04:05-0700 instance-id={{ .Instance.ID }} --wait",
			),
		),
		Cmd: "scw rdb backup export ={{ .Instance.ID }} --wait",
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
