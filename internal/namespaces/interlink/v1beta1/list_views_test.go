package interlink_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/interlink/v1beta1"
)

func Test_ListConnection(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: interlink.GetCommands(),
		Cmd:      "scw interlink connection list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListLink(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: interlink.GetCommands(),
		Cmd:      "scw interlink link list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPartner(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: interlink.GetCommands(),
		Cmd:      "scw interlink partner list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPop(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: interlink.GetCommands(),
		Cmd:      "scw interlink pop list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRoutingPolicy(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: interlink.GetCommands(),
		Cmd:      "scw interlink routing-policy list",
		Check:    core.TestCheckGolden(),
	}))
}
