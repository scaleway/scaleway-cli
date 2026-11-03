package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
)

func Test_ListServerType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: instance.GetCommands(),
		Cmd:      "scw instance server-type list",
		Check:    core.TestCheckGolden(),
	}))
}
