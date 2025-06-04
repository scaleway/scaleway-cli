package object_test

import (
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConfigInstall(t *testing.T) {
	client, err := scw.NewClient(
		scw.WithAuth(
			"SCWXXXXXXXXXXXXXXXXX",
			"11111111-1111-1111-1111-111111111111",
		),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	require.NoError(t, err)

	t.Run("NoExistingConfig", func(t *testing.T) {
		t.Run("rclone", core.Test(&core.TestConfig{
			Commands: object.GetCommands(),
			Cmd:      "scw object config install type=rclone",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					filePath := path.Join(
						ctx.OverrideEnv["HOME"],
						".config",
						"rclone",
						"rclone.conf",
					)
					assert.FileExists(t, filePath)
				},
				core.TestCheckExitCode(0),
			),
			TmpHomeDir: true,
			Client:     client,
		}))

		t.Run("mc", core.Test(&core.TestConfig{
			Commands: object.GetCommands(),
			Cmd:      "scw object config install type=mc",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".mc", "config.json")
					assert.FileExists(t, filePath)
				},
				core.TestCheckExitCode(0),
			),
			TmpHomeDir: true,
			Client:     client,
		}))

		t.Run("s3cmd", core.Test(&core.TestConfig{
			Commands: object.GetCommands(),
			Cmd:      "scw object config install type=s3cmd",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					t.Helper()
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
