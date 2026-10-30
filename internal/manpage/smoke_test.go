package manpage_test

import (
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_SmokeFullCommandSet builds pages for the real full command set and
// checks the coverage invariant: exactly one page per visible command plus
// the top-level `scw` page, no collisions, and a successful install (spec
// SC-001).
func Test_SmokeFullCommandSet(t *testing.T) {
	// GetDocGenContext provides a client in the meta, which argument default
	// functions (e.g. zone) rely on.
	ctx := core.GetDocGenContext()
	cmds := commands.GetCommands(ctx)

	// The command set may contain duplicate registrations of the same path
	// (e.g. v1 and v1beta1 namespace roots); dedupe the same way the builder
	// does, so this mirrors the effective command surface.
	uniqueVisible := map[string]bool{}
	for _, cmd := range cmds.GetAll() {
		if cmd.Hidden {
			continue
		}
		uniqueVisible[cmd.Namespace+"|"+cmd.Resource+"|"+cmd.Verb] = true
	}

	pages, root, err := manpage.BuildPages(ctx, cmds)
	require.NoError(t, err)
	require.NotNil(t, root)

	// Every visible command path has exactly one page.
	assert.Len(t, pages, len(uniqueVisible), "one page per visible command")

	pageNames := map[string]bool{}
	for _, page := range pages {
		assert.False(t, pageNames[page.FileName], "duplicate page %s", page.FileName)
		pageNames[page.FileName] = true

		// Basic structural invariants on every page.
		content, err := manpage.Render(page)
		require.NoError(t, err)
		assert.Contains(t, content, ".TH "+page.PageName+" 1")
		assert.Contains(t, content, ".SH NAME")
		assert.Contains(t, content, ".SH SYNOPSIS")
		assert.Contains(t, content, ".SH ENVIRONMENT")
	}

	// The top-level page is present and distinct from command pages.
	assert.False(t, pageNames["scw.1"], "top-level page must not clash with a command page")

	// Total installed = visible commands + top-level page, in less than
	// 2 minutes (contract I5 / SC-004).
	target := t.TempDir()
	start := time.Now()
	count, err := manpage.Install(ctx, cmds, target)
	duration := time.Since(start)
	require.NoError(t, err)
	assert.Equal(
		t,
		len(uniqueVisible)+1,
		count,
		"installed pages = visible commands + top-level page",
	)
	assert.Less(t, duration, 2*time.Minute, "install should complete in less than 2 minutes")
}
