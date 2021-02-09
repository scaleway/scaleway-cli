package k8s

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_GetVersion(t *testing.T) {
	////
	// Simple use cases
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		Cmd:        "scw k8s version get 1.20.2",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
