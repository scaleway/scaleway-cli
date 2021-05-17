package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_IPAttach(t *testing.T) {
	t.Run("With UUID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			createIP("Ip"),
		),
		Cmd: "scw instance ip attach {{ .Ip.Address }} server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("Ip"),
		),
		DisableParallel: true,
	}))

	t.Run("With IP", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			createIP("Ip"),
		),
		Cmd: "scw instance ip attach {{ .Ip.Address }} server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("Ip"),
		),
		DisableParallel: true,
	}))
}

func Test_IPDetach(t *testing.T) {
	t.Run("With UUID", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic"),
			core.ExecStoreBeforeCmd("Ip", "scw instance ip create"),
			core.ExecBeforeCmd("scw instance ip attach {{ .Ip.Address }} server-id={{ .Server.Id }}"),
		),
		Cmd: "scw instance ip detach {{ .Ip.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("Ip"),
		),
		DisableParallel: true,
	}))

	t.Run("With IP", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			core.ExecStoreBeforeCmd("Server", "scw instance server create stopped=true image=ubuntu-bionic"),
			core.ExecStoreBeforeCmd("Ip", "scw instance ip create"),
			core.ExecBeforeCmd("scw instance ip attach {{ .Ip.Address }} server-id={{ .Server.ID }}"),
		),
		Cmd: "scw instance ip detach {{ .Ip.Address }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: core.AfterFuncCombine(
			deleteServer("Server"),
			deleteIP("Ip"),
		),
		DisableParallel: true,
	}))

}
