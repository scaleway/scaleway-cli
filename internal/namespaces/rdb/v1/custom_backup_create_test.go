package rdb

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_CreateBackup(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Instance",
			fmt.Sprintf("scw rdb instance create node-type=DB-DEV-S is-ha-cluster=false name=%s engine=%s user-name=%s password=%s --wait", name, engine, user, password),
		),
		Cmd:   "scw rdb backup create instance-id={{ .Instance.ID }} --wait",
		Check: core.TestCheckGolden(),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw rdb instance delete {{ .Instance.ID }}"),
			core.ExecAfterCmd("scw rdb backup delete {{ .CmdResult.ID }}"),
		),
	}))
}
