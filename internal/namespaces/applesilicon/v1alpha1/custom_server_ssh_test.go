package applesilicon

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ServerSSH(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", "scw apple-silicon server create --wait"),
		),
		Cmd: "scw apple-silicon server ssh {{ .Server.ID }}",
		OverrideExec: core.OverrideExecSimple(
			"ssh {{ .Server.IP }} -p 22 -l m1 -t",
			0,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw apple-silicon server delete {{ .Server.ID }}"),
		),
		DisableParallel: true,
	}))

	t.Run("With-Exit-Code", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", "scw apple-silicon server create --wait"),
		),
		Cmd: "scw apple-silicon server ssh {{ .Server.ID }}",
		OverrideExec: core.OverrideExecSimple(
			"ssh {{ .Server.IP }} -p 22 -l m1 -t",
			130,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(130),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw apple-silicon server delete {{ .Server.ID }}"),
		),
		DisableParallel: true,
	}))
}
