package server

import (
	"context"
	"log"
	"slices"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
)

// MCPServer wraps the MCP server with CLI command integration
type MCPServer struct {
	server   *mcp.Server
	commands []*CommandTool
	readOnly bool
}

// NewMCPServer creates a new MCP server that exposes CLI commands as tools
func NewMCPServer(version string, cliCommands []*core.Command, readOnly bool) *MCPServer {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "scaleway-cli",
		Version: version,
	}, &mcp.ServerOptions{
		// Explicitly enable tools capability with listChanged notifications
		Capabilities: &mcp.ServerCapabilities{
			Tools: &mcp.ToolCapabilities{ListChanged: true},
		},
	})

	s := &MCPServer{
		server:   mcpServer,
		commands: make([]*CommandTool, 0, len(cliCommands)),
		readOnly: readOnly,
	}

	// Register all commands during initialization
	for _, cmd := range cliCommands {
		if err := s.RegisterCommand(cmd); err != nil {
			// Log but don't fail - some commands might not be compatible
			log.Printf(
				"Warning: failed to register command %s: %v\n",
				cmd.GetCommandLine("scw"),
				err,
			)
		}
	}

	return s
}

// shouldRegisterCommand returns true if the command should be registered as an MCP tool.
// It filters out:
// - Hidden commands
// - Commands without a Run function (namespace/resource containers)
// - All commands from the "mcp" namespace to avoid recursive server calls
// - All commands that are shell-centric
// - When readOnly is true, only commands with get/list verbs are registered
func shouldRegisterCommand(cmd *core.Command, readOnly bool) bool {
	// Skip hidden commands
	if cmd.Hidden {
		return false
	}

	// Skip commands without a Run function (namespace/resource containers)
	if cmd.Run == nil {
		return false
	}

	excludedNamespaces := []string{
		// Skip mcp namespace to avoid recursive server calls
		"mcp",
		// Skip config namespace to avoid security risks
		"config",
		// Skip shell-centric namespaces
		"alias",
		"autocomplete",
		"shell",
	}

	if slices.Contains(excludedNamespaces, cmd.Namespace) {
		return false
	}

	// In read-only mode, only allow get/list operations
	if readOnly && !isReadOnlyCommand(cmd) {
		return false
	}

	return true
}

// isReadOnlyCommand returns true if the command is a read-only operation
// (get, list, or get-* verbs)
func isReadOnlyCommand(cmd *core.Command) bool {
	if cmd.Verb == "" {
		return false
	}

	// Direct match for "get" or "list"
	if cmd.Verb == "get" || cmd.Verb == "list" {
		return true
	}

	// Match compound verbs that start with "get-" (e.g., "get-account", "get-credentials")
	if strings.HasPrefix(cmd.Verb, "get-") {
		return true
	}

	return false
}

// RegisterCommand registers a CLI command as an MCP tool
func (s *MCPServer) RegisterCommand(cmd *core.Command) error {
	if !shouldRegisterCommand(cmd, s.readOnly) {
		return nil
	}

	tool := NewCommandTool(cmd)
	mcpTool := tool.ToMCPTool()

	// Create a wrapper function for the tool using the correct MCP SDK signature
	wrapper := func(ctx context.Context, req *mcp.CallToolRequest, input map[string]any) (*mcp.CallToolResult, map[string]any, error) {
		result, err := tool.Execute(ctx, input)
		var output map[string]any
		if err != nil {
			// Return error - MCP SDK will wrap it
			return result, nil, err
		}
		// Extract text content for structured output
		if len(result.Content) > 0 {
			if tc, ok := result.Content[0].(*mcp.TextContent); ok {
				output = map[string]any{"result": tc.Text}
			}
		}

		return result, output, nil
	}

	// Register with MCP SDK
	mcp.AddTool(s.server, mcpTool, wrapper)

	s.commands = append(s.commands, tool)

	return nil
}

// Run starts the MCP server using the specified transport
func (s *MCPServer) Run(ctx context.Context, transport mcp.Transport) error {
	return s.server.Run(ctx, transport)
}

// Server returns the underlying MCP server for direct transport connections
func (s *MCPServer) Server() *mcp.Server {
	return s.server
}
