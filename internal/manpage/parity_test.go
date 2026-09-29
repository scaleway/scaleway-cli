package manpage_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/manpage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_Parity checks that a sample of real commands across several namespaces
// have man pages that cover all of their documented arguments, examples and
// the global environment variables (spec SC-003).
func Test_Parity(t *testing.T) {
	samples := [][]string{
		{"instance", "server", "create"},
		{"instance", "server", "list"},
		{"k8s", "cluster", "get"},
		{"rdb", "instance", "create"},
		{"lb", "lb", "list"},
		{"config", "dump"},
		{"secret", "secret", "create"},
		{"object", "bucket", "create"},
		{"iam", "user", "create"},
		{"version"},
		{"block", "volume", "get"},
		{"vpc", "private-network", "list"},
	}

	ctx := core.GetDocGenContext()
	cmds := commands.GetCommands(ctx)

	pages, _, err := manpage.BuildPages(ctx, cmds)
	require.NoError(t, err)

	pagesByName := map[string]*manpage.ManPage{}
	for _, page := range pages {
		pagesByName[page.PageName] = page
	}

	rendered := map[string]string{}
	renderPage := func(page *manpage.ManPage) string {
		t.Helper()
		content, ok := rendered[page.PageName]
		if !ok {
			var err error
			content, err = manpage.Render(page)
			require.NoError(t, err)
			rendered[page.PageName] = content
		}

		return content
	}

	for _, path := range samples {
		t.Run(strings.Join(path, "-"), func(t *testing.T) {
			cmd := cmds.Find(path...)
			require.NotNil(t, cmd, "command %v not found", path)
			require.False(t, cmd.Hidden, "command %v is hidden", path)

			page := pagesByName[manpage.PageName(path...)]
			require.NotNil(t, page, "no page for %v", path)
			content := renderPage(page)

			// Every documented argument appears in OPTIONS with its name and description.
			for _, arg := range cmd.ArgSpecs {
				if arg.Positional {
					assert.Contains(t, content, "<"+arg.Name+">")
				} else {
					assert.Contains(t, content, arg.Name+"=value")
				}
				assert.Contains(t, content, arg.Short)
			}

			// Every documented example appears in EXAMPLES.
			if len(cmd.Examples) > 0 {
				assert.Contains(t, content, ".SH EXAMPLES")
				for _, example := range cmd.Examples {
					assert.Contains(t, content, example.GetCommandLine("scw", cmd))
				}
			}

			// The global environment variables are documented.
			assert.Contains(t, content, ".SH ENVIRONMENT")
			assert.Contains(t, content, "SCW_ACCESS_KEY")
		})
	}
}

// Test_SubcommandsListed checks that non-leaf command pages list their direct
// subcommands in the SUBCOMMANDS section, and leaf pages do not.
func Test_SubcommandsListed(t *testing.T) {
	ctx := core.GetDocGenContext()
	cmds := commands.GetCommands(ctx)

	pages, _, err := manpage.BuildPages(ctx, cmds)
	require.NoError(t, err)

	pagesByName := map[string]*manpage.ManPage{}
	for _, page := range pages {
		pagesByName[page.PageName] = page
	}

	render := func(name string) string {
		t.Helper()
		page := pagesByName[name]
		require.NotNil(t, page, "no page %s", name)

		content, err := manpage.Render(page)
		require.NoError(t, err)

		return content
	}

	// A namespace page lists its resource subcommands.
	instance := render("scw-instance")
	assert.Contains(t, instance, ".SH SUBCOMMANDS")
	assert.Contains(t, instance, "\\fBscw instance server\\fP")
	assert.Contains(t, instance, "\\fBscw instance volume\\fP")

	// A resource page lists its verb subcommands.
	server := render("scw-instance-server")
	assert.Contains(t, server, ".SH SUBCOMMANDS")
	assert.Contains(t, server, "\\fBscw instance server create\\fP")
	assert.Contains(t, server, "\\fBscw instance server list\\fP")

	// A leaf command has no SUBCOMMANDS section.
	create := render("scw-instance-server-create")
	assert.NotContains(t, create, ".SH SUBCOMMANDS")
}

// Test_DeprecatedAnnotation checks that deprecated commands and arguments are
// explicitly flagged in the rendered page.
func Test_DeprecatedAnnotation(t *testing.T) {
	cmds := core.NewCommands(
		&core.Command{
			Namespace:  "demo",
			Resource:   "legacy",
			Verb:       "run",
			Short:      "A deprecated command",
			Long:       "This command will be removed.",
			Deprecated: true,
			ArgsType:   reflect.TypeOf(struct{}{}),
			ArgSpecs: core.ArgSpecs{
				{Name: "old", Short: "An old argument", Deprecated: true},
			},
		},
	)

	pages, _, err := manpage.BuildPages(core.GetDocGenContext(), cmds)
	require.NoError(t, err)

	var page *manpage.ManPage
	for _, p := range pages {
		if p.PageName == "scw-demo-legacy-run" {
			page = p
		}
	}
	require.NotNil(t, page)

	content, err := manpage.Render(page)
	require.NoError(t, err)

	assert.Contains(t, content, ".SH DEPRECATED")
	assert.Contains(t, content, "[Deprecated]")
}
