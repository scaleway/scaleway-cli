package docs_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/docs"
	"github.com/stretchr/testify/assert"
)

func TestSearchCommands(t *testing.T) {
	cmds := commands.GetCommands()
	results := docs.SearchCommands(cmds, "server")

	assert.NotEmpty(t, results, "should find commands matching 'server'")
}

func TestSearchCommands_Fuzzy(t *testing.T) {
	cmds := commands.GetCommands()

	// "isl" should fuzzy-match "instance server list" (subsequence)
	results := docs.SearchCommands(cmds, "isl")
	assert.NotEmpty(t, results, "should find commands via fuzzy subsequence match 'isl'")

	found := false
	for _, r := range results {
		if r.Path == "scw instance server list" {
			found = true
		}
	}
	assert.True(t, found, "fuzzy search 'isl' should match 'scw instance server list'")

	// "srl" should also fuzzy-match "server list" type commands
	results = docs.SearchCommands(cmds, "srl")
	assert.NotEmpty(t, results, "should find commands via fuzzy subsequence match 'srl'")
}

func TestSearchCommands_FuzzyRanking(t *testing.T) {
	cmds := commands.GetCommands()

	// "server" should rank exact substring matches higher than fuzzy-only matches
	results := docs.SearchCommands(cmds, "server")
	assert.NotEmpty(t, results)

	// The first result should contain "server" as a substring (exact match)
	// since exact matches score higher than pure subsequence matches
	assert.Contains(t, results[0].Path, "server")
}

func TestFuzzyMatch(t *testing.T) {
	tests := []struct {
		target string
		query  string
		match  bool
	}{
		{"instance server list", "isl", true},
		{"instance server list", "server", true},
		{"instance server list", "srl", true},
		{"instance server list", "xyz", false},
		{"instance server list", "", true},
		{"Instance Server List", "isl", true}, // case insensitive
		{"k8s cluster create", "kcc", true},
		{"k8s cluster create", "k8s", true},
		{"", "test", false},
	}

	for _, tt := range tests {
		t.Run(tt.target+"_"+tt.query, func(t *testing.T) {
			_, ok := docs.FuzzyMatch(tt.target, tt.query)
			assert.Equal(t, tt.match, ok, "fuzzyMatch(%q, %q)", tt.target, tt.query)
		})
	}
}

func TestSearchCommands_NoResults(t *testing.T) {
	cmds := commands.GetCommands()
	results := docs.SearchCommands(cmds, "xyznonexistentcommand12345")

	assert.Empty(t, results)
}

func TestBuildHierarchy(t *testing.T) {
	cmds := commands.GetCommands()
	hierarchy := docs.BuildHierarchy(cmds)

	assert.NotNil(t, hierarchy)
	assert.NotEmpty(t, hierarchy.Namespaces)

	_, hasInstance := hierarchy.Namespaces["instance"]
	assert.True(t, hasInstance, "should have instance namespace")

	_, hasHelp := hierarchy.Namespaces["help"]
	assert.True(t, hasHelp, "should have help namespace")
}

func TestRenderAllNamespaces(t *testing.T) {
	cmds := commands.GetCommands()
	output := docs.RenderAllNamespaces(cmds)

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Available namespaces:")
}

func TestRenderNamespacePage(t *testing.T) {
	cmds := commands.GetCommands()
	output := docs.RenderNamespacePage(cmds, "instance")

	assert.NotEmpty(t, output)
}

func TestRenderNamespacePage_NotFound(t *testing.T) {
	cmds := commands.GetCommands()
	output := docs.RenderNamespacePage(cmds, "nonexistent")

	assert.Contains(t, output, "not found")
}

func TestRenderHelpTopicPage(t *testing.T) {
	cmds := commands.GetCommands()
	helpCmd := cmds.Find("help", "output")
	if helpCmd == nil {
		t.Skip("help output command not found")
	}

	output := docs.RenderHelpTopicPage(helpCmd)
	assert.NotEmpty(t, output)
}

func TestRenderCommandPage(t *testing.T) {
	cmds := commands.GetCommands()
	ctx := core.GetDocGenContext()

	serverListCmd := cmds.Find("instance", "server", "list")
	if serverListCmd == nil {
		t.Skip("instance server list command not found")
	}

	output := docs.RenderCommandPage(ctx, serverListCmd, cmds)

	assert.NotEmpty(t, output)
	assert.Contains(t, output, "Usage:")
}

func Test_DocsList(t *testing.T) {
	t.Run("non-interactive", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw docs instance",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_DocsSearch(t *testing.T) {
	t.Run("search-server", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw docs search=server",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_DocsSearchNoResults(t *testing.T) {
	t.Run("no-results", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw docs search=xyznonexistent",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_DocsCommandPage(t *testing.T) {
	t.Run("instance-server-list", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw docs instance server list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}

func Test_DocsNamespaceNotFound(t *testing.T) {
	t.Run("not-found", core.Test(&core.TestConfig{
		Commands: commands.GetCommands(),
		Cmd:      "scw docs nonexistentns",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(1),
		),
	}))
}
