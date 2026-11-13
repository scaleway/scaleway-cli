package man_test

import (
	"bytes"
	"context"
	"os"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Man(t *testing.T) {
	t.Run("Install", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(context.Background()),
		Cmd:      "scw man install man-test-target",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				// The man section directory must have been created with the
				// top-level page and per-command pages.
				entries, err := os.ReadDir("man-test-target/man1")
				require.NoError(t, err)
				require.NotEmpty(t, entries)

				foundRootPage := false
				for _, entry := range entries {
					if entry.Name() == "scw.1" {
						foundRootPage = true
					}
				}
				assert.True(t, foundRootPage, "scw.1 not found in man-test-target/man1")

				// Clean up the install target created in the test directory.
				require.NoError(t, os.RemoveAll("man-test-target"))
			},
		),
	}))

	t.Run("Unattended", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(context.Background()),
		Stdin:    &bytes.Buffer{},
		Cmd:      "scw man install man-unattended-target",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stdout), "Installed")
				require.NoError(t, os.RemoveAll("man-unattended-target"))
			},
		),
	}))

	t.Run("TargetIsFile", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(context.Background()),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			file, err := os.Create("man-target-file")
			require.NoError(ctx.T, err)

			return file.Close()
		},
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			return os.Remove("man-target-file")
		},
		Cmd: "scw man install man-target-file",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, string(ctx.Stderr), "is not a directory")
			},
		),
	}))

	t.Run("NonWritableTarget", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(context.Background()),
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			if os.Geteuid() == 0 {
				ctx.T.Skip("cannot test non-writable target as root")
			}

			require.NoError(ctx.T, os.MkdirAll("man-readonly-target", 0o555))

			return nil
		},
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			if err := os.Chmod("man-readonly-target", 0o755); err != nil {
				return err
			}

			return os.Remove("man-readonly-target")
		},
		Cmd: "scw man install man-readonly-target",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.True(
					t,
					strings.Contains(string(ctx.Stderr), "man-readonly-target") &&
						strings.Contains(string(ctx.Stderr), "permission denied"),
					"expected a permission denied error naming the target, got: %s",
					string(ctx.Stderr),
				)
			},
		),
	}))
}
