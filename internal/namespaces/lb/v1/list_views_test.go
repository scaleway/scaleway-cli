package lb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
)

func Test_ListAcl(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb acl list frontend-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListBackend(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb backend list lb-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListFrontend(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb frontend list lb-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListIp(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb ip list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListLb(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb lb list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListLbTypes(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb lb-types list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListRoute(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb route list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSubscriber(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: lb.GetCommands(),
		Cmd:      "scw lb subscriber list name=x",
		Check:    core.TestCheckGolden(),
	}))
}
