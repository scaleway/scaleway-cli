package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_ServerSSH(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true ip=new")),
			startServer("Server"),
		),
		Cmd: "scw instance server ssh {{ .Server.ID }}",
		OverrideExec: core.OverrideExecSimple(
			"ssh {{ .Server.PublicIP.Address }} -p 22 -l root -t",
			0,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))

	t.Run("With-Exit-Code", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", testServerCommand("stopped=true ip=new")),
			startServer("Server"),
		),
		Cmd: "scw instance server ssh {{ .Server.ID }}",
		OverrideExec: core.OverrideExecSimple(
			"ssh {{ .Server.PublicIP.Address }} -p 22 -l root -t",
			130,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(130),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))

	t.Run("Stopped server", core.Test(&core.TestConfig{
		Commands:   instance.GetCommands(),
		BeforeFunc: createServerBionic("Server"),
		Cmd:        "scw instance server ssh {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))
}
