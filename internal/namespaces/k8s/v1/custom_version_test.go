package k8s

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_GetVersion(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw k8s version get " + kapsuleVersion,
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_ListVersion_Basic(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw k8s version list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
	}))
}
