package server_test

import (
	"fmt"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

// TestInventoryNilAnnotations prints a detailed inventory of all commands
// where any annotation field is nil.
//
// This test creates a real MCPServer instance and checks the annotations
// of all actually registered commands, ensuring the test matches production behavior.
//
//nolint:revive // This is a test function
func TestInventoryNilAnnotations(t *testing.T) {
	// Get all commands and create a real MCP server instance
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false)
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
				} else {
					if m.annotations.DestructiveHint == nil {
						nilFields = append(nilFields, "DestructiveHint")
					}
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
//nolint:revive // This is a test function
func TestAllRegisteredCommandsHaveAnnotations(t *testing.T) {
	// Get all commands and create a real MCP server instance
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false)
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
