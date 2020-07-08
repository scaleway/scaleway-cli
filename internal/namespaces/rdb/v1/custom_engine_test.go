package rdb

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_EngineList(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw rdb engine list",
		Check:    core.TestCheckGolden(),
	}))
}
