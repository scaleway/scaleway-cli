package k8s

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

////
// Simple use cases
////
func Test_GetVersion(t *testing.T) {
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw k8s version get 1.17.4",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("error", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw k8s version get test",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
