package rdb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Test_EngineList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: rdb.GetCommands(),
		Cmd:      "scw rdb engine list",
		Check:    core.TestCheckGolden(),
	}))
}
