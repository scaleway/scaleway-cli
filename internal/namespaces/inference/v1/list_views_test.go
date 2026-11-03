package inference_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/inference/v1"
)

func Test_ListDeployment(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: inference.GetCommands(),
		Cmd:      "scw inference deployment list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListModel(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: inference.GetCommands(),
		Cmd:      "scw inference model list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: inference.GetCommands(),
		Cmd:      "scw inference node-type list",
		Check:    core.TestCheckGolden(),
	}))
}
