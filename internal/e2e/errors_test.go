package e2e

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/namespaces/test/v1"
)

func TestStandardErrors(t *testing.T) {
	t.Run("unknown-command", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw bob",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))
}
