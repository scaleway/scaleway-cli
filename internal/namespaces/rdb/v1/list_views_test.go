package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
)

func Test_ListAcl(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb acl list instance-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListDatabase(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb database list instance-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListEndpoint(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb endpoint list 11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb node-type list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListPrivilege(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb privilege list instance-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListSnapshot(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb snapshot list",
		Check:    core.TestCheckGolden(),
	}))
}
