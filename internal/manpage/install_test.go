package manpage_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_InstallSyncOwnedFiles(t *testing.T) {
	ctx := core.GetDocGenContext()
	cmds := syntheticCommands()

	target := t.TempDir()
	sectionDir := filepath.Join(target, manpage.SectionDir)
	require.NoError(t, os.MkdirAll(sectionDir, 0o755))

	// First install.
	count, err := manpage.Install(ctx, cmds, target)
	require.NoError(t, err)
	require.Positive(t, count)

	// Seed stale owned and foreign content.
	staleOwned := filepath.Join(sectionDir, "scw-demo-removed-gone.1")
	require.NoError(t, os.WriteFile(staleOwned, []byte("stale"), 0o644))

	foreignFile := filepath.Join(sectionDir, "other-tool.1")
	require.NoError(t, os.WriteFile(foreignFile, []byte("foreign"), 0o644))

	foreignSubDir := filepath.Join(sectionDir, "man7")
	require.NoError(t, os.MkdirAll(foreignSubDir, 0o755))
	foreignPage := filepath.Join(foreignSubDir, "keep.7")
	require.NoError(t, os.WriteFile(foreignPage, []byte("keep"), 0o644))

	// Second install: the synthetic command set no longer contains any
	// previously installed command pages except the current ones.
	_, err = manpage.Install(ctx, cmds, target)
	require.NoError(t, err)

	// The stale owned page must be removed.
	_, err = os.Stat(staleOwned)
	assert.True(t, os.IsNotExist(err), "stale owned page was not removed")

	// Foreign files must be untouched.
	content, err := os.ReadFile(foreignFile)
	require.NoError(t, err)
	assert.Equal(t, "foreign", string(content))

	content, err = os.ReadFile(foreignPage)
	require.NoError(t, err)
	assert.Equal(t, "keep", string(content))
}

func Test_InstallCreatesTarget(t *testing.T) {
	ctx := core.GetDocGenContext()
	cmds := syntheticCommands()

	target := filepath.Join(t.TempDir(), "does", "not", "exist")

	count, err := manpage.Install(ctx, cmds, target)
	require.NoError(t, err)
	assert.Equal(
		t,
		4,
		count,
		"expected 4 pages: scw, scw-demo, scw-demo-widget, scw-demo-widget-create",
	)

	sectionDir := filepath.Join(target, manpage.SectionDir)
	entries, err := os.ReadDir(sectionDir)
	require.NoError(t, err)
	require.Len(t, entries, 4)

	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		names = append(names, entry.Name())
	}
	assert.ElementsMatch(
		t,
		[]string{"scw.1", "scw-demo.1", "scw-demo-widget.1", "scw-demo-widget-create.1"},
		names,
	)
}

func Test_InstallTargetIsFile(t *testing.T) {
	ctx := core.GetDocGenContext()
	cmds := syntheticCommands()

	target := filepath.Join(t.TempDir(), "not-a-dir")
	require.NoError(t, os.WriteFile(target, []byte("file"), 0o644))

	_, err := manpage.Install(ctx, cmds, target)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "is not a directory")
}
