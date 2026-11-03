package annotations_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/annotations/v1"
)

func Test_ListBinding(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: annotations.GetCommands(),
		Cmd:      "scw annotations binding list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListKey(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: annotations.GetCommands(),
		Cmd:      "scw annotations key list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListValue(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: annotations.GetCommands(),
		Cmd:      "scw annotations value list",
		Check:    core.TestCheckGolden(),
	}))
}
