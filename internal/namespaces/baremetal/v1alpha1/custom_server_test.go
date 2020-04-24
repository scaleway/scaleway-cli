package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_StartServerErrors(t *testing.T) {
	t.Run("Error: cannot be started while not delivered", core.Test(&core.TestConfig{
		BeforeFunc:  createServerAndWait("Server"),
		Commands:    GetCommands(),
		Cmd:         "scw baremetal server start {{ .Server.ID }} -w",
		Check:       core.TestCheckExitCode(1),
		AfterFunc:   deleteServer("Server"),
		DefaultZone: scw.ZoneFrPar2,
	}))
}

func Test_StopServerErrors(t *testing.T) {
	t.Run("Error: cannot be stopped while not delivered", core.Test(&core.TestConfig{
		BeforeFunc:  createServerAndWait("Server"),
		Commands:    GetCommands(),
		Cmd:         "scw baremetal server stop {{ .Server.ID }} -w",
		Check:       core.TestCheckExitCode(1),
		AfterFunc:   deleteServer("Server"),
		DefaultZone: scw.ZoneFrPar2,
	}))
}

func Test_RebootServerErrors(t *testing.T) {
	t.Run("Error: cannot be rebooted while not delivered", core.Test(&core.TestConfig{
		BeforeFunc:  createServerAndWait("Server"),
		Commands:    GetCommands(),
		Cmd:         "scw baremetal server reboot {{ .Server.ID }} -w",
		Check:       core.TestCheckExitCode(1),
		AfterFunc:   deleteServer("Server"),
		DefaultZone: scw.ZoneFrPar2,
	}))
}
