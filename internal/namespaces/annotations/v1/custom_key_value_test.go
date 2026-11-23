package annotations_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	annotations "github.com/scaleway/scaleway-cli/v2/internal/namespaces/annotations/v1"
)

func Test_ListKeyValue(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: annotations.GetCommands(),
		Cmd:      "scw annotations key-value list",
		Check:    core.TestCheckGolden(),
	}))
}
