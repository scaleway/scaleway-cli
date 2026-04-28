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
	server            *mcp.Server
	commands          []*CommandTool
	readOnly          bool
	enabledNamespaces []string
	enabledResources  []string
	enabledVerbs      []string
}

// NewMCPServer creates a new MCP server that exposes CLI commands as tools
func NewMCPServer(
	version string,
	cliCommands []*core.Command,
	readOnly bool,
	enabledNamespaces, enabledResources, enabledVerbs []string,
) *MCPServer {
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
		server:            mcpServer,
		commands:          make([]*CommandTool, 0, len(cliCommands)),
		readOnly:          readOnly,
		enabledNamespaces: enabledNamespaces,
		enabledResources:  enabledResources,
		enabledVerbs:      enabledVerbs,
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

var (
	ExcludedNamespaces = []string{
		// Skip config namespace to avoid security risks
		"config",
		// Skip shell-centric namespaces
		"alias",
		"autocomplete",
		"shell",
		"login",
		"init",
	}

	ExcludedVerbs = []string{
		"edit", // Shell-centric verb
	}

	ExcludedResources = []string{}
)

// ShouldRegisterCommand returns true if the command should be registered as an MCP tool.
// It filters out:
// - Hidden commands
// - Commands with ExcludeFromMCP flag set
// - Commands without a Run function (namespace/resource containers)
// - Commands in excluded namespaces
// - Commands with excluded verbs
// - When readOnly is true, only commands with get/list verbs are registered
// - When enabledNamespaces/ Resources/ Verbs are set, only matching commands are registered
func ShouldRegisterCommand(
	cmd *core.Command,
	readOnly bool,
	enabledNamespaces, enabledResources, enabledVerbs []string,
) bool {
	// Skip hidden commands
	if cmd.Hidden {
		return false
	}

	if cmd.ExcludeFromMCP {
		return false
	}

	// Skip commands without a Run function (namespace/resource containers)
	if cmd.Run == nil {
		return false
	}

	if slices.Contains(ExcludedNamespaces, cmd.Namespace) {
		return false
	}

	if slices.Contains(ExcludedVerbs, cmd.Verb) {
		return false
	}

	if slices.Contains(ExcludedResources, cmd.Resource) {
		return false
	}

	// If enabled namespaces are specified, only allow those namespaces
	if len(enabledNamespaces) > 0 && !slices.Contains(enabledNamespaces, cmd.Namespace) {
		return false
	}

	// If enabled resources are specified, only allow those resources
	if len(enabledResources) > 0 && !slices.Contains(enabledResources, cmd.Resource) {
		return false
	}

	// If enabled verbs are specified, only allow those verbs
	if len(enabledVerbs) > 0 && !slices.Contains(enabledVerbs, cmd.Verb) {
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
	if !ShouldRegisterCommand(
		cmd,
		s.readOnly,
		s.enabledNamespaces,
		s.enabledResources,
		s.enabledVerbs,
	) {
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

// RegisteredCommands returns the list of commands registered as MCP tools
func (s *MCPServer) RegisteredCommands() []*CommandTool {
	return s.commands
}
