package mcp_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp"
)

func TestMcpCommands(t *testing.T) {
	cmds := mcp.GetCommands()
	if cmds == nil {
		t.Fatal("GetCommands() returned nil")
	}

	all := cmds.GetAll()
	if len(all) == 0 {
		t.Fatal("Expected at least one MCP command")
	}
}
