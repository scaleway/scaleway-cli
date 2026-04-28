package server_test

import (
	"encoding/json"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

// TestToolMetaSerialization verifies that the Meta field is properly
// serialized to JSON as _meta when tools are listed.
func TestToolMetaSerialization(t *testing.T) {
	allCommands := commands.GetCommands().GetAll()

	mcpServer := server.NewMCPServer("test-version", allCommands, false, nil, nil, nil)
	registeredCommands := mcpServer.RegisteredCommands()

	if len(registeredCommands) == 0 {
		t.Fatal("No commands registered")
	}

	// Test the first command
	tool := registeredCommands[0]
	mcpTool := tool.ToMCPTool()

	// Marshal to JSON to see what the client actually receives
	jsonBytes, err := json.MarshalIndent(mcpTool, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal tool: %v", err)
	}

	t.Logf("Tool name: %s", mcpTool.Name)
	t.Logf("JSON representation:\n%s", string(jsonBytes))

	// Verify _meta field exists in JSON
	var rawJSON map[string]interface{}
	if err := json.Unmarshal(jsonBytes, &rawJSON); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	meta, ok := rawJSON["_meta"]
	if !ok {
		t.Error("_meta field not present in JSON")
	} else {
		t.Logf("_meta content: %v", meta)

		// Check for expected keys
		metaMap, ok := meta.(map[string]interface{})
		if !ok {
			t.Error("_meta is not a map")
		} else {
			if _, ok := metaMap["namespace"]; !ok {
				t.Error("_meta.namespace missing")
			}
			if _, ok := metaMap["resource"]; !ok {
				t.Error("_meta.resource missing")
			}
			if _, ok := metaMap["verb"]; !ok {
				t.Error("_meta.verb missing")
			}
		}
	}
}
