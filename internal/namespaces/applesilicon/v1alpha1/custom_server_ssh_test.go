package applesilicon_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
)

func Test_ServerSSH(t *testing.T) {
	t.Skip("Cannot run test as it's lock for 24h before deletion")
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Server",
				"scw apple-silicon server create server-type=M2-M --wait",
			),
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
		Commands: applesilicon.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd(
				"Server",
				"scw apple-silicon server create server-type=M2-M --wait",
			),
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
