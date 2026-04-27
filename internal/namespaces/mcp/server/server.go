package server

import (
	"context"
	"log"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
)

// MCPServer wraps the MCP server with CLI command integration
type MCPServer struct {
	server   *mcp.Server
	commands []*CommandTool
}

// NewMCPServer creates a new MCP server that exposes CLI commands as tools
func NewMCPServer(version string, cliCommands []*core.Command) *MCPServer {
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

// RegisterCommand registers a CLI command as an MCP tool
func (s *MCPServer) RegisterCommand(cmd *core.Command) error {
	// Skip hidden commands
	if cmd.Hidden {
		return nil
	}

	// Skip commands without a Run function (namespace/resource containers)
	if cmd.Run == nil && cmd.Verb == "" {
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
