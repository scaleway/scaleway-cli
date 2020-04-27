package object

import (
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConfigInstall(t *testing.T) {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
	}

	client, err := scw.NewClient(clientOpts...)
	require.NoError(t, err)

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
			Client:     client,
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
			Client:     client,
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
			Client:     client,
		}))
	})
}
