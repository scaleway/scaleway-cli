package registry

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_Logout(t *testing.T) {
	t.Run("docker", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry logout program=docker",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple("docker logout rg.fr-par.scw.cloud", 0),
	}))
	t.Run("podman", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry logout program=podman",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExec: core.OverrideExecSimple("podman logout rg.fr-par.scw.cloud", 0),
	}))
}
