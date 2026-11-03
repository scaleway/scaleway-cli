package ipam_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipam/v1"
)

func Test_ListIp(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: ipam.GetCommands(),
		Cmd:      "scw ipam ip list",
		Check:    core.TestCheckGolden(),
	}))
}
