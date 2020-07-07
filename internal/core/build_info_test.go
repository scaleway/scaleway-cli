package core

import (
	"context"
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/alecthomas/assert"
	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

var fakeCommand = &Command{
	Namespace:            "plop",
	ArgsType:             reflect.TypeOf(args.RawArgs{}),
	AllowAnonymousClient: true,
	Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		return &SuccessResult{}, nil
	},
}

// These tests needs to run in sequence since they are modifying a file on the filesystem
func Test_CheckVersion(t *testing.T) {
	t.Run("Outdated version", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: &BuildInfo{
			Version: version.Must(version.NewSemver("v1.20")),
		},
		Cmd: "scw plop",
		Check: TestCheckCombine(
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.Equal(t, "a new version of scw is available (2.0.0-beta.4), beware that you are currently running 1.20.0\n", ctx.LogBuffer)
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Up to date version", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: &BuildInfo{
			Version: version.Must(version.NewSemver("v99.99")),
		},
		Cmd: "scw plop -D",
		Check: TestCheckCombine(
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.Contains(t, ctx.LogBuffer, "version is up to date (99.99.0)\n")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Already checked", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: &BuildInfo{
			Version: version.Must(version.NewSemver("v1.0")),
		},
		BeforeFunc: func(ctx *BeforeFuncCtx) error {
			return createAndCloseFile(getLatestVersionUpdateFilePath(ctx.OverrideEnv[scw.ScwCacheDirEnv]))
		},
		Cmd: "scw plop -D",
		Check: TestCheckCombine(
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.Contains(t, ctx.LogBuffer, "version was already checked during past 24 hours\n")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("Check more than 24h ago", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: &BuildInfo{
			Version: version.Must(version.NewSemver("v1.0")),
		},
		BeforeFunc: func(ctx *BeforeFuncCtx) error {
			filePath := getLatestVersionUpdateFilePath(ctx.OverrideEnv[scw.ScwCacheDirEnv])
			err := createAndCloseFile(filePath)
			require.NoError(t, err)
			twoDaysAgo := time.Now().Local().Add(-2 * time.Hour * 24)
			return os.Chtimes(filePath, twoDaysAgo, twoDaysAgo)
		},
		Cmd: "scw plop",
		Check: TestCheckCombine(
			func(t *testing.T, ctx *CheckFuncCtx) {
				assert.Contains(t, ctx.LogBuffer, "a new version of scw is available (2.0.0-beta.4), beware that you are currently running 1.0.0\n")
			},
		),
		TmpHomeDir: true,
	}))
}
