package object

import (
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_ConfigInstall(t *testing.T) {
	tmpDir := os.TempDir()

	t.Run("NoExistingConfig", func(t *testing.T) {
		t.Run("rclone", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=rclone",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					os.Stat(path.Join(tmpDir, ".config", "rclone", "rclone.conf"))
				},
				core.TestCheckExitCode(0),
			),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))

		t.Run("mc", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=mc",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					os.Stat(path.Join(tmpDir, ".mc", "config.json"))
				},
				core.TestCheckExitCode(0),
			),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))

		t.Run("s3cmd", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=s3cmd",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					os.Stat(path.Join(tmpDir, ".s3cfg"))
				},
				core.TestCheckExitCode(0),
			),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))
	})
}
