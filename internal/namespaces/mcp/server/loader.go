package server

import (
	"context"
	"log"
	"slices"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
)

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
)

// ShouldLoadCommand returns true if the command should be registered as an MCP tool.
// It filters out:
// - Hidden commands
// - Commands with ExcludeFromMCP flag set
// - Commands without a Run function (namespace/resource containers)
// - Commands in excluded namespaces
// - Commands with excluded verbs
// - When readOnly is true, only commands with get/list verbs are registered
// - When enabledNamespaces/ Resources/ Verbs are set, only matching commands are registered
func ShouldLoadCommand(
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
	if readOnly && !cmd.IsReadOnlyCommand() {
		return false
	}

	return true
}

// LoadCommand loads a CLI command as an MCP tool and optionally as a resource
func (s *MCPServer) LoadCommand(cmd *core.Command) error {
	if !ShouldLoadCommand(
		cmd,
		s.readOnly,
		s.enabledNamespaces,
		s.enabledResources,
		s.enabledVerbs,
	) {
		return nil
	}

	// Register as a tool
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

	// Register as a resource if it's a list command
	if cmd.IsListCommand() {
		if err := s.LoadResource(cmd); err != nil {
			log.Printf(
				"Warning: failed to load resource %s: %v\n",
				cmd.GetCommandLine("scw"),
				err,
			)
		}
	}

	return nil
}
