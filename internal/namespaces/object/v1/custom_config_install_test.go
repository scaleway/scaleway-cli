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
					filePath := path.Join(tmpDir, ".config", "rclone", "rclone.conf")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
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
					filePath := path.Join(tmpDir, ".mc", "config.json")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
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
					filePath := path.Join(tmpDir, ".s3cfg")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
				},
				core.TestCheckExitCode(0),
			),
			OverrideEnv: map[string]string{
				"HOME": tmpDir,
			},
		}))
	})
}
