package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Server
func Test_ListServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))

	t.Run("List with tags", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal server list tags.0=a",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))
}
