package key_manager_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/key_manager/v1alpha1"
)

func Test_ListKey(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: key_manager.GetCommands(),
		Cmd:      "scw keymanager key list",
		Check:    core.TestCheckGolden(),
	}))
}
