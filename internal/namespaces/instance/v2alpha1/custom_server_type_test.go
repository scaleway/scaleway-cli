package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
)

func Test_ServerTypeList(t *testing.T) {
	t.Run("on fr-par-1", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list zone=fr-par-2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("on nl-ams-2", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list zone=nl-ams-2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("on pl-waw-1", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list zone=pl-waw-1",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("on it-mil-1", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list zone=it-mil-1",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_Server_GetCompatibleTypes(t *testing.T) {
	cmds := instance.GetCommands()
	cmds.Merge(instanceV1.GetCommands())

	t.Run("DEV1-S", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			"scw instance server create stopped=true image=ubuntu-resolute ip=none type=DEV1-S",
		),
		Cmd: "scw instance server get-compatible-types {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("L4-2-24G", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.ExecStoreBeforeCmd(
			"Server",
			"scw instance server create stopped=true image=ubuntu_noble_gpu_os_13_nvidia ip=none type=L4-2-24G",
		),
		Cmd: "scw instance server get-compatible-types {{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
