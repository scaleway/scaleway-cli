package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_ServerTypeList(t *testing.T) {
	t.Run("on fr-par-1", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list",
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
