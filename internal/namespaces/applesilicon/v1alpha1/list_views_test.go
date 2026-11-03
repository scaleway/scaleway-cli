package applesilicon_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
)

func Test_ListOs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon os list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPrivateNetwork(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon private-network list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRunner(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon runner list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: applesilicon.GetCommands(),
		Cmd:      "scw apple-silicon server list",
		Check:    core.TestCheckGolden(),
	}))
}
