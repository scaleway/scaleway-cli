package applesilicon_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
)

func Test_ServerTypeList(t *testing.T) {
	t.Run("base", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon server-type list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
