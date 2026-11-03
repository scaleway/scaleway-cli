package messageq_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/messageq/v1alpha1"
)

func Test_ListDeployment(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: messageq.GetCommands(),
		Cmd:      "scw messageq deployment list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: messageq.GetCommands(),
		Cmd:      "scw messageq node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: messageq.GetCommands(),
		Cmd:      "scw messageq user list deployment-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: messageq.GetCommands(),
		Cmd:      "scw messageq version list",
		Check:    core.TestCheckGolden(),
	}))
}
