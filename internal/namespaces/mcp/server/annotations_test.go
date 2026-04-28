package server_test

import (
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

// TestInventoryNilAnnotations prints a detailed inventory of all commands
// where any annotation field is nil.
//
// This test creates a real MCPServer instance and checks the annotations
// of all actually registered commands, ensuring the test matches production behavior.
//

func TestInventoryNilAnnotations(t *testing.T) {
	// Get all commands and create a real MCP server instance
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false, nil, nil, nil)
	registeredCommands := mcpServer.RegisteredCommands()

	type nilAnnotation struct {
		commandLine string
		namespace   string
		resource    string
		verb        string
		annotations *mcp.ToolAnnotations
	}

	var nilAnnotations []nilAnnotation

	for _, tool := range registeredCommands {
		mcpTool := tool.ToMCPTool()

		if mcpTool.Annotations == nil {
			nilAnnotations = append(nilAnnotations, nilAnnotation{
				commandLine: tool.Command.GetCommandLine("scw"),
				namespace:   tool.Command.Namespace,
				resource:    tool.Command.Resource,
				verb:        tool.Command.Verb,
				annotations: nil,
			})

			continue
		}

		// Check which fields are nil
		ann := mcpTool.Annotations

		// OpenWorldHint is a pointer - check if nil
		if ann.OpenWorldHint == nil {
			nilAnnotations = append(nilAnnotations, nilAnnotation{
				commandLine: tool.Command.GetCommandLine("scw"),
				namespace:   tool.Command.Namespace,
				resource:    tool.Command.Resource,
				verb:        tool.Command.Verb,
				annotations: ann,
			})

			continue
		}

		// DestructiveHint is a pointer - check if nil
		if ann.DestructiveHint == nil {
			nilAnnotations = append(nilAnnotations, nilAnnotation{
				commandLine: tool.Command.GetCommandLine("scw"),
				namespace:   tool.Command.Namespace,
				resource:    tool.Command.Resource,
				verb:        tool.Command.Verb,
				annotations: ann,
			})
		}
	}

	// Print detailed inventory
	t.Logf("=== INVENTORY: Commands with nil annotation fields ===")
	t.Logf("Total commands analyzed: %d", len(allCommands))
	t.Logf("Commands registered to MCP server: %d", len(registeredCommands))
	t.Logf("Commands with nil annotation fields: %d", len(nilAnnotations))
	t.Logf("")

	if len(nilAnnotations) > 0 {
		// Group by namespace for readability
		byNamespace := make(map[string][]nilAnnotation)
		for _, m := range nilAnnotations {
			byNamespace[m.namespace] = append(byNamespace[m.namespace], m)
		}

		for ns, cmds := range byNamespace {
			t.Logf("Namespace: %s (%d commands)", ns, len(cmds))

			for _, m := range cmds {
				nilFields := []string{}
				if m.annotations == nil {
					nilFields = append(nilFields, "ALL_ANNOTATIONS_NIL")
				} else if m.annotations.DestructiveHint == nil {
					nilFields = append(nilFields, "DestructiveHint")
				}

				t.Logf("  - %s (verb=%s): nil fields [%s]",
					m.commandLine, m.verb, fmt.Sprint(nilFields))
			}
			t.Logf("")
		}
	}
}

// TestAllRegisteredCommandsHaveAnnotations verifies that all commands registered
// by the MCP server have non-nil annotations with all fields set.
//
// This test creates a real MCPServer instance to ensure it tests the actual
// production behavior.
//

func TestAllRegisteredCommandsHaveAnnotations(t *testing.T) {
	// Get all commands and create a real MCP server instance
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false, nil, nil, nil)
	registeredCommands := mcpServer.RegisteredCommands()

	var failedCommands []string

	for _, tool := range registeredCommands {
		mcpTool := tool.ToMCPTool()

		// Check annotations exist
		if mcpTool.Annotations == nil {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: annotations is nil",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))

			continue
		}

		// Check OpenWorldHint is set
		if mcpTool.Annotations.OpenWorldHint == nil {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: OpenWorldHint is nil",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))
		}

		// Check DestructiveHint is set
		if mcpTool.Annotations.DestructiveHint == nil {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: DestructiveHint is nil",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))
		}
	}

	if len(failedCommands) > 0 {
		t.Errorf("%d commands have missing annotations:", len(failedCommands))

		for i, failure := range failedCommands {
			if i >= 20 {
				t.Logf("... and %d more", len(failedCommands)-20)

				break
			}
			t.Logf("  - %s", failure)
		}
	}
}

// TestAllToolsHaveCommandMetadata verifies that all MCP tools have
// namespace, resource, and verb metadata set in their Meta field.
//

func TestAllToolsHaveCommandMetadata(t *testing.T) {
	// Get all commands and create a real MCP server instance
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false, nil, nil, nil)
	registeredCommands := mcpServer.RegisteredCommands()

	var failedCommands []string

	for _, tool := range registeredCommands {
		mcpTool := tool.ToMCPTool()

		// Check Meta field exists
		if mcpTool.Meta == nil {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: Meta field is nil",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))

			continue
		}

		// Check namespace is set
		if ns, ok := mcpTool.Meta["namespace"].(string); !ok || ns == "" {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: namespace metadata not set or empty",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))
		}

		// Check resource is set (can be empty for namespace-level commands)
		if _, ok := mcpTool.Meta["resource"].(string); !ok {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: resource metadata not set",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))
		}

		// Check verb is set (can be empty for namespace/resource-level commands)
		if _, ok := mcpTool.Meta["verb"].(string); !ok {
			failedCommands = append(failedCommands, fmt.Sprintf(
				"%s/%s/%s: verb metadata not set",
				tool.Command.Namespace, tool.Command.Resource, tool.Command.Verb))
		}
	}

	if len(failedCommands) > 0 {
		t.Errorf("%d commands have missing metadata:", len(failedCommands))

		for i, failure := range failedCommands {
			if i >= 20 {
				t.Logf("... and %d more", len(failedCommands)-20)

				break
			}
			t.Logf("  - %s", failure)
		}
	}
}

// TestToolMetadataValues verifies that metadata values match the command's
// actual namespace, resource, and verb fields.
func TestToolMetadataValues(t *testing.T) {
	testCases := []struct {
		name     string
		command  *core.Command
		wantNS   string
		wantRes  string
		wantVerb string
	}{
		{
			name: "full command",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
			},
			wantNS:   "instance",
			wantRes:  "server",
			wantVerb: "list",
		},
		{
			name: "namespace only",
			command: &core.Command{
				Namespace: "config",
			},
			wantNS:   "config",
			wantRes:  "",
			wantVerb: "",
		},
		{
			name: "namespace and resource",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "volume",
			},
			wantNS:   "instance",
			wantRes:  "volume",
			wantVerb: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			tool := server.NewCommandTool(tc.command)
			mcpTool := tool.ToMCPTool()

			if mcpTool.Meta == nil {
				t.Fatal("Meta field is nil")
			}

			if got, ok := mcpTool.Meta["namespace"].(string); !ok || got != tc.wantNS {
				t.Errorf("namespace: got %q, want %q", got, tc.wantNS)
			}

			if got, ok := mcpTool.Meta["resource"].(string); !ok || got != tc.wantRes {
				t.Errorf("resource: got %q, want %q", got, tc.wantRes)
			}

			if got, ok := mcpTool.Meta["verb"].(string); !ok || got != tc.wantVerb {
				t.Errorf("verb: got %q, want %q", got, tc.wantVerb)
			}
		})
	}
}
