package object

import (
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_ConfigInstall(t *testing.T) {
	tmpDir := os.TempDir()
	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithEnv(),
		scw.WithUserAgent("cli-e2e-test"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
	}

	client, err := scw.NewClient(clientOpts...)
	if err != nil {
		t.Fail()
	}

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
			Client: client,
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
			Client: client,
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
			Client: client,
		}))
	})
}
