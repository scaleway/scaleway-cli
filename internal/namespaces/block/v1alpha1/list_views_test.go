package block_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
)

func Test_ListSnapshot(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		Cmd:      "scw block snapshot list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVolume(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		Cmd:      "scw block volume list",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_ListVolumeType(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: block.GetCommands(),
		Cmd:      "scw block volume-type list",
		Check:    core.TestCheckGolden(),
	}))
}
