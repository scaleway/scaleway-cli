package manpage_test

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// syntheticCommands returns a small deterministic command set used to test
// rendering without depending on the full CLI command tree.
func syntheticCommands() *core.Commands {
	return core.NewCommands(
		&core.Command{
			Namespace: "demo",
			Short:     "Demo namespace",
			ArgsType:  reflect.TypeOf(struct{}{}),
			ArgSpecs:  core.ArgSpecs{},
		},
		&core.Command{
			Namespace: "demo",
			Resource:  "widget",
			Short:     "Widget management commands",
			Long:      "Manage the demo widgets.",
			ArgsType:  reflect.TypeOf(struct{}{}),
			ArgSpecs:  core.ArgSpecs{},
		},
		&core.Command{
			Namespace: "demo",
			Resource:  "widget",
			Verb:      "create",
			Short:     "Create a widget",
			Long:      "Create a new widget\nin the demo namespace.",
			ArgsType:  reflect.TypeOf(struct{}{}),
			ArgSpecs: core.ArgSpecs{
				{Name: "name", Short: "The name of the widget", Required: true},
				{
					Name:       "zone",
					Short:      "Zone to target",
					EnumValues: []string{"fr-par-1", "nl-ams-1"},
					Default: func(_ context.Context) (string, string) {
						return "fr-par-1", "fr-par-1"
					},
				},
				{Name: "legacy-flag", Short: "An old flag", Deprecated: true},
				{Name: "server-id", Short: "The ID of the server", Positional: true},
			},
			Examples: []*core.Example{
				{
					Short: "Create a widget",
					Raw:   "scw demo widget create name=demo server-id=11111111-1111-1111-1111-111111111111",
				},
			},
			SeeAlsos: []*core.SeeAlso{
				{Short: "List widgets", Command: "scw demo widget list"},
				{Short: "Widget documentation", Command: "https://docs.scaleway.com/widget"},
			},
		},
	)
}

// writeGolden writes the golden file when UPDATE_MAN_GOLDENS is set, otherwise
// it compares the output against the existing golden file.
func writeGolden(t *testing.T, name string, content string) {
	t.Helper()
	path := filepath.Join("testdata", name)

	if os.Getenv("UPDATE_MAN_GOLDENS") != "" {
		require.NoError(t, os.MkdirAll("testdata", 0o755))
		require.NoError(t, os.WriteFile(path, []byte(content), 0o644))
		t.Logf("updated golden file %s", path)

		return
	}

	expected, err := os.ReadFile(path)
	require.NoErrorf(
		t,
		err,
		"golden file %s not found: run with UPDATE_MAN_GOLDENS=1 to create it",
		path,
	)
	assert.Equal(t, string(expected), content, "golden file %s mismatch", path)
}

func Test_RenderCommandPage(t *testing.T) {
	cmds := syntheticCommands()

	pages, _, err := manpage.BuildPages(core.GetDocGenContext(), cmds)
	require.NoError(t, err)

	var page *manpage.ManPage
	for _, p := range pages {
		if p.PageName == "scw-demo-widget-create" {
			page = p
		}
	}
	require.NotNil(t, page, "page scw-demo-widget-create not found")

	rendered, err := manpage.Render(page)
	require.NoError(t, err)

	writeGolden(t, "render-command-page.golden", rendered)

	// Structural assertions on the roff output.
	assert.Contains(t, rendered, `.TH scw-demo-widget-create 1`)
	assert.Contains(t, rendered, ".SH NAME")
	assert.Contains(
		t,
		rendered,
		"scw-demo-widget-create, scw demo widget create \\- Create a widget",
	)
	assert.Contains(t, rendered, ".SH SYNOPSIS")
	assert.Contains(t, rendered, "\\fBscw demo widget create\\fP")
	assert.Contains(t, rendered, ".SH DESCRIPTION")
	assert.Contains(t, rendered, "Create a new widget")
	assert.Contains(t, rendered, ".SH OPTIONS")
	assert.Contains(t, rendered, "\\fBname=value\\fP")
	assert.Contains(t, rendered, "The name of the widget")
	assert.Contains(t, rendered, "[Required]")
	assert.Contains(t, rendered, "[Default: fr-par-1]")
	assert.Contains(t, rendered, "[One of: fr-par-1, nl-ams-1]")
	assert.Contains(t, rendered, "[Deprecated]")
	assert.Contains(t, rendered, "\\fB<server-id>\\fP")
	assert.Contains(t, rendered, ".SH EXAMPLES")
	assert.Contains(t, rendered, "Create a widget")
	assert.Contains(
		t,
		rendered,
		"scw demo widget create name=demo server-id=11111111-1111-1111-1111-111111111111",
	)
	assert.Contains(t, rendered, ".SH SEE ALSO")
	assert.Contains(t, rendered, "scw-demo-widget-list")
	assert.Contains(t, rendered, "https://docs.scaleway.com/widget")
	assert.Contains(t, rendered, ".SH ENVIRONMENT")
	// A leaf command has no SUBCOMMANDS section.
	assert.NotContains(t, rendered, ".SH SUBCOMMANDS")
}

func Test_RenderNamespacePage(t *testing.T) {
	cmds := syntheticCommands()

	pages, _, err := manpage.BuildPages(core.GetDocGenContext(), cmds)
	require.NoError(t, err)

	var page *manpage.ManPage
	for _, p := range pages {
		if p.PageName == "scw-demo" {
			page = p
		}
	}
	require.NotNil(t, page, "page scw-demo not found")

	rendered, err := manpage.Render(page)
	require.NoError(t, err)

	writeGolden(t, "render-namespace-page.golden", rendered)

	// The namespace page lists its direct subcommands as command lines.
	assert.Contains(t, rendered, ".SH SUBCOMMANDS")
	assert.Contains(t, rendered, "\\fBscw demo widget\\fP")
	assert.Contains(t, rendered, "Widget management commands")
	// Verbs two levels down are not listed on the namespace page.
	assert.NotContains(t, rendered, "scw demo widget create")
}

func Test_RenderRootPage(t *testing.T) {
	cmds := syntheticCommands()

	_, root, err := manpage.BuildPages(core.GetDocGenContext(), cmds)
	require.NoError(t, err)
	require.NotNil(t, root)

	rendered, err := manpage.RenderRoot(root)
	require.NoError(t, err)

	writeGolden(t, "render-root-page.golden", rendered)

	assert.Contains(t, rendered, `.TH scw 1`)
	assert.Contains(t, rendered, "scw \\- Scaleway command line tool")
	assert.Contains(t, rendered, ".SH OPTIONS")
	assert.Contains(t, rendered, "--profile")
	assert.Contains(t, rendered, ".SH NAMESPACES")
	assert.Contains(t, rendered, "\\fBdemo\\fP")
	assert.Contains(t, rendered, "Demo namespace")
	assert.Contains(t, rendered, ".SH ENVIRONMENT")
}

func Test_RenderEscaping(t *testing.T) {
	cmds := core.NewCommands(
		&core.Command{
			Namespace: "demo",
			Resource:  "tricky",
			Verb:      "run",
			Short:     "Tricky content",
			Long:      ".SH FAKE\n\\backslash and \x1b[1mansi bold\x1b[0m text",
			ArgsType:  reflect.TypeOf(struct{}{}),
			ArgSpecs:  core.ArgSpecs{},
		},
	)

	pages, _, err := manpage.BuildPages(core.GetDocGenContext(), cmds)
	require.NoError(t, err)

	var page *manpage.ManPage
	for _, p := range pages {
		if p.PageName == "scw-demo-tricky-run" {
			page = p
		}
	}
	require.NotNil(t, page)

	rendered, err := manpage.Render(page)
	require.NoError(t, err)

	// A line starting with a dot must be neutralized so roff does not
	// interpret it as a macro.
	assert.Contains(t, rendered, `\&.SH FAKE`)
	// Backslashes must be escaped.
	assert.Contains(t, rendered, `\Ebackslash`)
	// ANSI escape sequences must be stripped.
	assert.NotContains(t, rendered, "\x1b[1m")
	assert.Contains(t, rendered, "ansi bold")
}
