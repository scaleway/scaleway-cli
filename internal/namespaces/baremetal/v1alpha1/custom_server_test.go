package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_StartServerErrors(t *testing.T) {
	t.Run("Error: cannot be started while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: createServer("Server"),
		Commands:   GetCommands(),
		Cmd:        "scw baremetal server start server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitServerAfter("Server"),
			deleteServer("Server"),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))
}

func Test_StopServerErrors(t *testing.T) {
	t.Run("Error: cannot be stopped while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: createServer("Server"),
		Commands:   GetCommands(),
		Cmd:        "scw baremetal server stop server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitServerAfter("Server"),
			deleteServer("Server"),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))
}

func Test_RebootServerErrors(t *testing.T) {
	t.Run("Error: cannot be rebooted while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: createServer("Server"),
		Commands:   GetCommands(),
		Cmd:        "scw baremetal server reboot server-id={{ .Server.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitServerAfter("Server"),
			deleteServer("Server"),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))
}
