package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
)

func Test_ListPrivateNic(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance private-nic list server-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSnapshot(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance snapshot list",
		Check:    core.TestCheckGolden(),
	}))
}
