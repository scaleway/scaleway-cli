package baremetal

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Server
//
func Test_ListServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw baremetal server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))

	t.Run("List with tags", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw baremetal server list tags.0=a",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DefaultZone: scw.ZoneFrPar2,
	}))
}
