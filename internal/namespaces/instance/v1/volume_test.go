package instance

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_VolumeCreate(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create name=test size=20G",
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			ctx.ExecuteCmd("scw instance volume delete volume-id={{ .Volume.ID }}")
			return nil
		},
		Check: core.TestCheckGolden(),
	}))

	t.Run("Bad size unit", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume create name=test size=20",
		Check:    core.TestCheckGolden(),
	}))
}

func Test_VolumeDelete(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume delete -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume delete",
		Check:    core.TestCheckGolden(),
	}))

}

func Test_VolumeGet(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume get -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume get",
		Check:    core.TestCheckGolden(),
	}))

}
func Test_VolumeList(t *testing.T) {

	t.Run("Usage", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume list -h",
		Check:    core.TestCheckGolden(),
	}))

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw instance volume list",
		Check:    core.TestCheckGolden(),
	}))

}
