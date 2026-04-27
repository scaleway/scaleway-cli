package mcp

import (
	"testing"
)

func TestMcpCommands(t *testing.T) {
	cmds := GetCommands()
	if cmds == nil {
		t.Fatal("GetCommands() returned nil")
	}

	all := cmds.GetAll()
	if len(all) == 0 {
		t.Fatal("Expected at least one MCP command")
	}
}
