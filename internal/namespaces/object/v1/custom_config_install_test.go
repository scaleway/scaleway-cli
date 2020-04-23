package object

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_ConfigInstall(t *testing.T) {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
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
					fmt.Println(ctx.Err)
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".config", "rclone", "rclone.conf")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
				},
				core.TestCheckExitCode(0),
			),
			MockHomeDir: true,
			Client:      client,
		}))

		t.Run("mc", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=mc",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".mc", "config.json")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
				},
				core.TestCheckExitCode(0),
			),
			MockHomeDir: true,
			Client:      client,
		}))

		t.Run("s3cmd", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw object config install type=s3cmd",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					filePath := path.Join(ctx.OverrideEnv["HOME"], ".s3cfg")
					_, err := os.Stat(filePath)
					if err != nil {
						t.Logf("No file at %s", filePath)
						t.Fail()
					}
				},
				core.TestCheckExitCode(0),
			),
			MockHomeDir: true,
			Client:      client,
		}))
	})
}
