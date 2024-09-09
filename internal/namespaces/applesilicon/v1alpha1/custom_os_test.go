package applesilicon_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
)

func Test_OSDisplay(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon os get e08d1e5d-b4b9-402a-9f9a-97732d17e374",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		DisableParallel: false,
	}))
}
