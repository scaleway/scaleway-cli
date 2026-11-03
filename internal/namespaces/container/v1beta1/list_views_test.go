package container_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1beta1"
)

func Test_ListToken(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: container.GetCommands(),
		Cmd:      "scw container token list",
		Check:    core.TestCheckGolden(),
	}))
}
