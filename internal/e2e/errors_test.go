package e2e_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/test/v1"
)

func TestStandardErrors(t *testing.T) {
	t.Skip("Test API not available")

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
