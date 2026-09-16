package mcp

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		mcpRoot(),
		server.McpServer(),
		server.McpServerServe(),
		server.McpServerListTools(),
		server.McpServerListResources(),
	)
}

func mcpRoot() *core.Command {
	return &core.Command{
		Groups:         []string{"utility"},
		Namespace:      "mcp",
		Short:          "MCP (Model Context Protocol) server management",
		Long:           "Commands for managing the MCP server that exposes Scaleway CLI commands as AI tools.",
		ExcludeFromMCP: true, // Skip mcp namespace to avoid recursive server calls
	}
}
