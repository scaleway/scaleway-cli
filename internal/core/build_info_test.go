package core

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/internal/args"
)

var fakeCommand = &Command{
	Namespace:            "plop",
	ArgsType:             reflect.TypeOf(args.RawArgs{}),
	AllowAnonymousClient: true,
	Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
		return &SuccessResult{}, nil
	},
}

func deleteLatestVersionUpdateFile(*BeforeFuncCtx) error {
	os.Remove(getLatestVersionUpdateFilePath())
	return nil
}

// These tests needs to run in sequence since they are modifying a file on the filesystem
func Test_CheckVersion(t *testing.T) {
	t.Run("Outdated version", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: BuildInfo{
			Version: version.Must(version.NewSemver("v1.20")),
		},
		BeforeFunc: deleteLatestVersionUpdateFile,
		Cmd:        "scw plop",
		Check: TestCheckCombine(
			TestCheckStderrGolden(),
			TestCheckExitCode(0),
		),
		DisableParallel: true,
	}))

	t.Run("Up to date version", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: BuildInfo{
			Version: version.Must(version.NewSemver("v99.99")),
		},
		BeforeFunc: deleteLatestVersionUpdateFile,
		Cmd:        "scw plop -D",
		Check: TestCheckCombine(
			TestCheckStderrGolden(),
			TestCheckExitCode(0),
		),
		DisableParallel: true,
	}))

	t.Run("Already checked", Test(&TestConfig{
		Commands: NewCommands(fakeCommand),
		BuildInfo: BuildInfo{
			Version: version.Must(version.NewSemver("v1.0")),
		},
		BeforeFunc: func(ctx *BeforeFuncCtx) error {
			if createAndCloseFile(getLatestVersionUpdateFilePath()) {
				return nil
			}
			return fmt.Errorf("failed to create latestVersionUpdateFile")
		},
		Cmd: "scw plop -D",
		Check: TestCheckCombine(
			TestCheckStderrGolden(),
			TestCheckExitCode(0),
		),
		DisableParallel: true,
	}))
}
