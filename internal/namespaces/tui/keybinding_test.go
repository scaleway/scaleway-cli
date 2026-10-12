//go:build !wasm

package tui_test

import (
	"os"
	"path"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/tui"
	"github.com/stretchr/testify/assert"
)

func Test_TuiKeybindingList(t *testing.T) {
	t.Run("default", core.Test(&core.TestConfig{
		Commands: tui.GetCommands(),
		Cmd:      "scw tui keybinding list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				assert.Contains(t, lineFor(ctx.Stdout, "refresh"), "[r]")
				assert.Contains(t, lineFor(ctx.Stdout, "quit"), "[q]")
			},
		),
		TmpHomeDir: true,
	}))

	t.Run("with overrides", core.Test(&core.TestConfig{
		Commands:   tui.GetCommands(),
		BeforeFunc: writeCliConfig("tui:\n    keybindings:\n        refresh: R\n        menu: b\n"),
		Cmd:        "scw tui keybinding list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				// The cli.yaml overrides are applied (key column),
				// next to the default (default column).
				refresh := lineFor(ctx.Stdout, "refresh")
				assert.Contains(t, refresh, "[R]")
				assert.Contains(t, refresh, "[r]")
				menu := lineFor(ctx.Stdout, "menu")
				assert.Contains(t, menu, "[b]")
				assert.Contains(t, menu, "[m g]")
			},
		),
		TmpHomeDir: true,
	}))
}

// lineFor returns the table line of the given action.
func lineFor(stdout []byte, action string) string {
	for line := range strings.SplitSeq(string(stdout), "\n") {
		if strings.HasPrefix(line, action) {
			return line
		}
	}

	return ""
}

func writeCliConfig(content string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		dir := path.Join(ctx.OverrideEnv["HOME"], ".config", "scw")
		if err := os.MkdirAll(dir, 0o755); err != nil {
			return err
		}

		return os.WriteFile(path.Join(dir, "cli.yaml"), []byte(content), 0o644)
	}
}
