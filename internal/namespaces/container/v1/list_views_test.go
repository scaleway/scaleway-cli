package container_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
)

func Test_ListContainer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		Cmd:      "scw container container list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDomain(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		Cmd:      "scw container domain list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNamespace(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		Cmd:      "scw container namespace list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListTrigger(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		Cmd:      "scw container trigger list",
		Check:    core.TestCheckGolden(),
	}))
}
