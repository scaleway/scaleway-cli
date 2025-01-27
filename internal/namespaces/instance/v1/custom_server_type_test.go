package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_ServerTypeList(t *testing.T) {
	t.Run("server-type list", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
