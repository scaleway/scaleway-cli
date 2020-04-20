package object

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ConfigGet(t *testing.T) {
	t.Run("rclone", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw object config get type=rclone",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("mc", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw object config get type=mc",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))

	t.Run("s3cmd", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw object config get type=s3cmd",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_ConfigInstall(t *testing.T) {
	t.Run("rclone", func(t *testing.T) {

	})

	t.Run("mc", func(t *testing.T) {

	})

	t.Run("s3cmd", func(t *testing.T) {

	})
}
