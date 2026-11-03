package searchdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/searchdb/v1alpha1"
)

func Test_ListDeployment(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: searchdb.GetCommands(),
		Cmd:      "scw searchdb deployment list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListNodeTypes(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: searchdb.GetCommands(),
		Cmd:      "scw searchdb node-types list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListUser(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: searchdb.GetCommands(),
		Cmd:      "scw searchdb user list deployment-id=11111111-1111-1111-1111-111111111111",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVersions(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: searchdb.GetCommands(),
		Cmd:      "scw searchdb versions list",
		Check:    core.TestCheckGolden(),
	}))
}
