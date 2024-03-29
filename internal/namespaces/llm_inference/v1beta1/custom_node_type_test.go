package llm_inference_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	llm_inference "github.com/scaleway/scaleway-cli/v2/internal/namespaces/llm_inference/v1beta1"
)

func Test_ListNodeType(t *testing.T) {
	cmds := llm_inference.GetCommands()

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw llm-inference node-type list",
		Check:    core.TestCheckGolden(),
	}))
}
