package baremetal_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
)

func Test_ListOptions(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal options list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListOs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal os list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPrivateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal private-network list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSettings(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: baremetal.GetCommands(),
		Cmd:      "scw baremetal settings list project-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}
