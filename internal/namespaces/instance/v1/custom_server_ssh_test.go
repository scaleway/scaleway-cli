package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ServerSSH(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServer("Server"),
			startServer("Server"),
		),
		Cmd: "scw instance server ssh {{ .Server.ID }}",
		OverrideExec: core.OverrideExecSimple(
			"/usr/bin/ssh {{ .Server.PublicIP.Address }} -p 22 -l root -t",
			0,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc:       deleteServer("Server"),
		DisableParallel: true,
	}))
}
