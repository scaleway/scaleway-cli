package core_test

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var fakeCommand = &core.Command{
	Namespace:            "plop",
	ArgsType:             reflect.TypeOf(args.RawArgs{}),
	AllowAnonymousClient: true,
	Run: func(_ context.Context, _ any) (i any, e error) {
		return &core.SuccessResult{}, nil
	},
}

// These tests needs to run in sequence since they are modifying a file on the filesystem
func Test_CheckVersion(t *testing.T) {
	t.Run("Outdated version", core.Test(&core.TestConfig{
		Commands: core.NewCommands(fakeCommand),
		BuildInfo: &core.BuildInfo{
			Version: version.Must(version.NewSemver("v1.20")),
		},
		Cmd: "scw plop",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Equal(
					t,
					"A new version of scw is available (2.5.4), beware that you are currently running 1.20.0\n",
					ctx.LogBuffer,
				)
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Up to date version", core.Test(&core.TestConfig{
		Commands: core.NewCommands(fakeCommand),
		BuildInfo: &core.BuildInfo{
			Version: version.Must(version.NewSemver("v99.99")),
		},
		Cmd: "scw plop -D",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, ctx.LogBuffer, "version is up to date (99.99.0)\n")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Already checked", core.Test(&core.TestConfig{
		Commands: core.NewCommands(fakeCommand),
		BuildInfo: &core.BuildInfo{
			Version: version.Must(version.NewSemver("v1.0")),
		},
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			return core.CreateAndCloseFile(
				core.GetLatestVersionUpdateFilePath(ctx.OverrideEnv[scw.ScwCacheDirEnv]),
			)
		},
		Cmd: "scw plop -D",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(
					t,
					ctx.LogBuffer,
					"version was already checked during past 24 hours\n",
				)
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Check more than 24h ago", core.Test(&core.TestConfig{
		Commands: core.NewCommands(fakeCommand),
		BuildInfo: &core.BuildInfo{
			Version: version.Must(version.NewSemver("v1.0")),
		},
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			filePath := core.GetLatestVersionUpdateFilePath(ctx.OverrideEnv[scw.ScwCacheDirEnv])
			err := core.CreateAndCloseFile(filePath)
			require.NoError(t, err)
			twoDaysAgo := time.Now().Local().Add(-2 * time.Hour * 24)

			return os.Chtimes(filePath, twoDaysAgo, twoDaysAgo)
		},
		Cmd: "scw plop",
		Check: core.TestCheckCombine(
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(
					t,
					ctx.LogBuffer,
					"A new version of scw is available (2.5.4), beware that you are currently running 1.0.0\n",
				)
			},
		),
		TmpHomeDir: true,
	}))
}
