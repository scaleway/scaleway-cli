package inference_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	inference "github.com/scaleway/scaleway-cli/v2/internal/namespaces/inference/v1beta1"
)

func Test_ListModel(t *testing.T) {
	cmds := inference.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw inference model list",
		Check:    core.TestCheckGolden(),
	}))
}
