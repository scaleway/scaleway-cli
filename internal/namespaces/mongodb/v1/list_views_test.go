package mongodb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
)

func Test_ListDatabase(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb database list instance-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListInstance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb instance list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListMaintenance(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb maintenance list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSnapshot(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb snapshot list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb user list instance-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersion(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: mongodb.GetCommands(),
		Cmd:      "scw mongodb version list",
		Check:    core.TestCheckGolden(),
	}))
}
