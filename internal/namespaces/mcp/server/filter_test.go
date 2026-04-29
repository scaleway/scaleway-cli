package server_test

import (
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func TestShouldRegisterCommand_WithEnabledNamespaces(t *testing.T) {
	tests := []struct {
		name              string
		command           *core.Command
		enabledNamespaces []string
		expected          bool
	}{
		{
			name: "command in enabled namespace",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance", "iam"},
			expected:          true,
		},
		{
			name: "command not in enabled namespace",
			command: &core.Command{
				Namespace: "object",
				Resource:  "bucket",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance", "iam"},
			expected:          false,
		},
		{
			name: "no enabled namespaces allows all",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{},
			expected:          true,
		},
		{
			name: "excluded namespace takes precedence",
			command: &core.Command{
				Namespace: "config",
				Resource:  "path",
				Verb:      "get",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"config"},
			expected:          false, // config is in ExcludedNamespaces
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.ShouldLoadCommand(
				tt.command,
				false,
				tt.enabledNamespaces,
				nil,
				nil,
			)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestShouldRegisterCommand_WithEnabledResources(t *testing.T) {
	tests := []struct {
		name             string
		command          *core.Command
		enabledResources []string
		expected         bool
	}{
		{
			name: "command with enabled resource",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledResources: []string{"server", "volume"},
			expected:         true,
		},
		{
			name: "command with non-enabled resource",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "image",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledResources: []string{"server", "volume"},
			expected:         false,
		},
		{
			name: "no enabled resources allows all",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledResources: []string{},
			expected:         true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.ShouldLoadCommand(tt.command, false, nil, tt.enabledResources, nil)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestShouldRegisterCommand_WithEnabledVerbs(t *testing.T) {
	tests := []struct {
		name         string
		command      *core.Command
		enabledVerbs []string
		expected     bool
	}{
		{
			name: "command with enabled verb",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledVerbs: []string{"get", "list"},
			expected:     true,
		},
		{
			name: "command with non-enabled verb",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "delete",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledVerbs: []string{"get", "list"},
			expected:     false,
		},
		{
			name: "no enabled verbs allows all",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledVerbs: []string{},
			expected:     true,
		},
		{
			name: "excluded verb takes precedence",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "edit",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledVerbs: []string{"edit", "list"},
			expected:     false, // edit is in ExcludedVerbs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.ShouldLoadCommand(tt.command, false, nil, nil, tt.enabledVerbs)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestShouldRegisterCommand_WithCombinedFilters(t *testing.T) {
	tests := []struct {
		name              string
		command           *core.Command
		enabledNamespaces []string
		enabledResources  []string
		enabledVerbs      []string
		expected          bool
	}{
		{
			name: "all filters match",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance"},
			enabledResources:  []string{"server"},
			enabledVerbs:      []string{"list"},
			expected:          true,
		},
		{
			name: "namespace does not match",
			command: &core.Command{
				Namespace: "object",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance"},
			enabledResources:  []string{"server"},
			enabledVerbs:      []string{"list"},
			expected:          false,
		},
		{
			name: "resource does not match",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "volume",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance"},
			enabledResources:  []string{"server"},
			enabledVerbs:      []string{"list"},
			expected:          false,
		},
		{
			name: "verb does not match",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "create",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance"},
			enabledResources:  []string{"server"},
			enabledVerbs:      []string{"list"},
			expected:          false,
		},
		{
			name: "only namespace filter set",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{"instance"},
			enabledResources:  []string{},
			enabledVerbs:      []string{},
			expected:          true,
		},
		{
			name: "only resource filter set",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{},
			enabledResources:  []string{"server"},
			enabledVerbs:      []string{},
			expected:          true,
		},
		{
			name: "only verb filter set",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			enabledNamespaces: []string{},
			enabledResources:  []string{},
			enabledVerbs:      []string{"list"},
			expected:          true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.ShouldLoadCommand(
				tt.command,
				false,
				tt.enabledNamespaces,
				tt.enabledResources,
				tt.enabledVerbs,
			)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestShouldRegisterCommand_WithReadOnlyAndFilters(t *testing.T) {
	tests := []struct {
		name              string
		command           *core.Command
		readOnly          bool
		enabledNamespaces []string
		enabledVerbs      []string
		expected          bool
	}{
		{
			name: "read-only mode with matching verb and namespace",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			readOnly:          true,
			enabledNamespaces: []string{"instance"},
			enabledVerbs:      []string{"list"},
			expected:          true,
		},
		{
			name: "read-only mode with non-read-only verb",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "delete",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			readOnly:          true,
			enabledNamespaces: []string{"instance"},
			enabledVerbs:      []string{"delete"},
			expected:          false, // read-only mode blocks delete
		},
		{
			name: "read-only mode with namespace not matching",
			command: &core.Command{
				Namespace: "object",
				Resource:  "bucket",
				Verb:      "list",
				Run:       func(ctx context.Context, argsI any) (any, error) { return nil, nil },
			},
			readOnly:          true,
			enabledNamespaces: []string{"instance"},
			enabledVerbs:      []string{"list"},
			expected:          false, // namespace filter blocks it
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.ShouldLoadCommand(
				tt.command,
				tt.readOnly,
				tt.enabledNamespaces,
				nil,
				tt.enabledVerbs,
			)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
