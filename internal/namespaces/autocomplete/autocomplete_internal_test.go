package autocomplete

import (
	"context"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/assert"
)

func TestDefaultBasename(t *testing.T) {
	run := func(binaryName string, expected string) func(*testing.T) {
		return func(t *testing.T) {
			t.Helper()
			ctx := core.InjectMeta(context.Background(), &core.Meta{BinaryName: binaryName})
			assert.Equal(t, expected, defaultBasename(ctx))
		}
	}

	t.Run("bare name", run("scw", "scw"))
	t.Run("absolute path", run("/opt/homebrew/Cellar/scw/2.65.2/bin/scw", "scw"))
	t.Run("relative path", run("./bin/scw", "scw"))
	t.Run("empty", run("", "scw"))
}

func TestZshScriptRegistersBasename(t *testing.T) {
	script := TrimText(autocompleteScripts(context.Background(), "scw")["zsh"].CompleteFunc)

	// compinit only registers an autoloaded _scw file when the first line is a #compdef tag.
	assert.True(t, strings.HasPrefix(script, "#compdef scw\n"), script)
	assert.Contains(t, script, "compdef _scw scw")
}
