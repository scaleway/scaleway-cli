package dedibox_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/dedibox/v1"
)

func Test_ListFip(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox fip list project-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListOffer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox offer list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListOs(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox os list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRpnInfo(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox rpn-info list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRpnV1(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox rpn-v1 list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRpnV2(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox rpn-v2 list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSan(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox san list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListServer(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox server list project-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListService(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: dedibox.GetCommands(),
		Cmd:      "scw dedibox service list",
		Check:    core.TestCheckGolden(),
	}))
}
