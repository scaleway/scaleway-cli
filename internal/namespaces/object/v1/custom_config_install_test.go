package object

import (
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/stretchr/testify/assert"
)

func Test_ConfigInstall(t *testing.T) {
	t.Run("NoExistingConfig", func(t *testing.T) {
		t.Run("rclone", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=rclone",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".config", "rclone", "rclone.conf")
					assert.FileExists(t, filePath)
				},
				core.TestCheckExitCode(0),
			),
			TmpHomeDir: true,
		}))

		t.Run("mc", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=mc",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".mc", "config.json")
					assert.FileExists(t, filePath)
				},
				core.TestCheckExitCode(0),
			),
			TmpHomeDir: true,
		}))

		t.Run("s3cmd", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=s3cmd",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".s3cfg")
					assert.FileExists(t, filePath)
				},
				core.TestCheckExitCode(0),
			),
			TmpHomeDir: true,
		}))
	})
}
