package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
)

func Test_StartServerErrors(t *testing.T) {
	t.Run("Error: cannot be started while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable, offerFilterHDD),
			createServer(),
		),
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal server start {{ .Server.ID }} zone={{ .Server.Zone }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitForServer(),
			deleteServer(),
		),
	}))
}

func Test_StopServerErrors(t *testing.T) {
	t.Run("Error: cannot be stopped while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable, offerFilterHDD),
			createServer(),
		),
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal server stop {{ .Server.ID }} zone={{ .Server.Zone }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitForServer(),
			deleteServer(),
		),
	}))
}

func Test_RebootServerErrors(t *testing.T) {
	t.Run("Error: cannot be rebooted while not delivered", core.Test(&core.TestConfig{
		BeforeFunc: core.BeforeFuncCombine(
			selectOffer(zone, offerFilterAvailable, offerFilterHDD),
			createServer(),
		),
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal server reboot {{ .Server.ID }} zone={{ .Server.Zone }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
		AfterFunc: core.AfterFuncCombine(
			waitForServer(),
			deleteServer(),
		),
	}))
}
