package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v3"
)

func Test_ListPrivateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal private-network list",
		Check:    core.TestCheckGolden(),
	}))
}
