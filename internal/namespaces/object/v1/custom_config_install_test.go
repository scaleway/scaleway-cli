package object

import (
	"os"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ConfigInstall(t *testing.T) {
	tmpDir := os.TempDir()

	t.Run("NoExistingConfig", func(t *testing.T) {
		t.Run("rclone", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=rclone",
			Check:    core.TestCheckExitCode(0),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))

		t.Run("mc", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=mc",
			Check:    core.TestCheckExitCode(0),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))

		t.Run("s3cmd", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=s3cmd",
			Check:    core.TestCheckExitCode(0),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))
	})
}
