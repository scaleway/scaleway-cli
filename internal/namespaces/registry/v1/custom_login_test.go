package registry

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_Login(t *testing.T) {
	t.Run("docker", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry login program=docker",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExecCommand: map[string]core.ExecCmd{"docker": dockerFakeCommand},
	}))
	t.Run("podman", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw registry login program=podman",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideExecCommand: map[string]core.ExecCmd{"podman": podmanFakeCommand},
	}))
}
