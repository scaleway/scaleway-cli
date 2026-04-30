package server_test

import (
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	mcp "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
	"github.com/stretchr/testify/assert"
)

func TestMcpServerList(t *testing.T) {
	cmds := core.NewCommandsMerge(mcp.GetCommands(), k8s.GetCommands(), account.GetCommands())

	t.Run("Basic", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With Namespace Filter", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools namespace=k8s",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With Resource Filter", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools resource=project",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With Verb Filter", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools verb=list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With ReadOnly Filter", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools read-only=true",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With Namespace Filter That Matches Nothing", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools namespace=nonexistent",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("With Combined Filters", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd:      "scw mcp server list-tools namespace=k8s resource=cluster verb=list",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func TestShouldRegisterCommand(t *testing.T) {
	testCmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "get",
		Short:     "Test command",
		Hidden:    false,
		Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
		ArgSpecs:  core.ArgSpecs{},
		Examples:  []*core.Example{},
	}

	t.Run("Should register visible command", func(t *testing.T) {
		result := server.ShouldLoadCommand(testCmd, server.CommandFilterConfig{})
		assert.True(t, result)
	})

	t.Run("Should not register hidden command", func(t *testing.T) {
		hiddenCmd := &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "get",
			Short:     "Test command",
			Hidden:    true,
			Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(hiddenCmd, server.CommandFilterConfig{})
		assert.False(t, result)
	})

	t.Run("Should not register command without Run function", func(t *testing.T) {
		noRunCmd := &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "get",
			Short:     "Test command",
			Hidden:    false,
			Run:       nil,
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(noRunCmd, server.CommandFilterConfig{})
		assert.False(t, result)
	})

	t.Run("Should not register excluded namespace", func(t *testing.T) {
		configCmd := &core.Command{
			Namespace: "config",
			Resource:  "resource",
			Verb:      "get",
			Short:     "Test command",
			Hidden:    false,
			Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(configCmd, server.CommandFilterConfig{})
		assert.False(t, result)
	})

	t.Run("Should not register excluded verb", func(t *testing.T) {
		editCmd := &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "edit",
			Short:     "Test command",
			Hidden:    false,
			Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(editCmd, server.CommandFilterConfig{})
		assert.False(t, result)
	})

	t.Run(
		"Should not register command in readOnly mode if not read-only operation",
		func(t *testing.T) {
			createCmd := &core.Command{
				Namespace: "test",
				Resource:  "resource",
				Verb:      "create",
				Short:     "Test command",
				Hidden:    false,
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
				ArgSpecs:  core.ArgSpecs{},
				Examples:  []*core.Example{},
			}
			result := server.ShouldLoadCommand(
				createCmd,
				server.CommandFilterConfig{ReadOnly: true},
			)
			assert.False(t, result)
		},
	)

	t.Run("Should register get command in readOnly mode", func(t *testing.T) {
		getCmd := &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "get",
			Short:     "Test command",
			Hidden:    false,
			Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(getCmd, server.CommandFilterConfig{ReadOnly: true})
		assert.True(t, result)
	})

	t.Run("Should register list command in readOnly mode", func(t *testing.T) {
		listCmd := &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "list",
			Short:     "Test command",
			Hidden:    false,
			Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			ArgSpecs:  core.ArgSpecs{},
			Examples:  []*core.Example{},
		}
		result := server.ShouldLoadCommand(listCmd, server.CommandFilterConfig{ReadOnly: true})
		assert.True(t, result)
	})

	t.Run("Should filter by enabled namespace", func(t *testing.T) {
		result := server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledNamespaces: []string{"other"}},
		)
		assert.False(t, result)

		result = server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledNamespaces: []string{"test"}},
		)
		assert.True(t, result)
	})

	t.Run("Should filter by enabled resource", func(t *testing.T) {
		result := server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledResources: []string{"other"}},
		)
		assert.False(t, result)

		result = server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledResources: []string{"resource"}},
		)
		assert.True(t, result)
	})

	t.Run("Should filter by enabled verb", func(t *testing.T) {
		result := server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledVerbs: []string{"other"}},
		)
		assert.False(t, result)

		result = server.ShouldLoadCommand(
			testCmd,
			server.CommandFilterConfig{EnabledVerbs: []string{"get"}},
		)
		assert.True(t, result)
	})
}
