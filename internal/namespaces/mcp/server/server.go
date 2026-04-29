package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"sort"
	"strings"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
)

// MCPServer wraps the MCP server with CLI command integration
type MCPServer struct {
	server            *mcp.Server
	commands          []*CommandTool
	resources         []*CommandResource
	readOnly          bool
	enabledNamespaces []string
	enabledResources  []string
	enabledVerbs      []string
}

// NewMCPServer creates a new MCP server that exposes CLI commands as tools and resources
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
		// Enable tools and resources capabilities with listChanged notifications
		Capabilities: &mcp.ServerCapabilities{
			Tools:     &mcp.ToolCapabilities{ListChanged: true},
			Resources: &mcp.ResourceCapabilities{ListChanged: true},
		},
	})

	s := &MCPServer{
		server:            mcpServer,
		commands:          make([]*CommandTool, 0, len(cliCommands)),
		resources:         make([]*CommandResource, 0),
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

// RegisterCommand registers a CLI command as an MCP tool and optionally as a resource
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
	if IsListCommand(cmd) {
		if err := s.RegisterResource(cmd); err != nil {
			log.Printf(
				"Warning: failed to register resource %s: %v\n",
				cmd.GetCommandLine("scw"),
				err,
			)
		}
	}

	return nil
}

// IsListCommand returns true if the command is a list operation
func IsListCommand(cmd *core.Command) bool {
	if cmd.Verb == "" {
		return false
	}

	// Direct match for "list"
	if cmd.Verb == "list" {
		return true
	}

	return false
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

// RegisteredResources returns the list of commands registered as MCP resources
func (s *MCPServer) RegisteredResources() []*CommandResource {
	return s.resources
}

// ListTools returns a list of available MCP tools based on the server's configuration.
// It returns tools that match the read-only mode and enabled namespaces/resources/verbs.
func (s *MCPServer) ListTools() []toolInfo {
	var tools []toolInfo
	for _, cmd := range s.commands {
		if ShouldRegisterCommand(
			cmd.Command,
			s.readOnly,
			s.enabledNamespaces,
			s.enabledResources,
			s.enabledVerbs,
		) {
			tools = append(tools, toolInfo{
				Namespace: cmd.Command.Namespace,
				Resource:  cmd.Command.Resource,
				Verb:      cmd.Command.Verb,
				ToolName:  CommandNameToToolName(cmd.Command),
				Short:     cmd.Command.Short,
			})
		}
	}

	// Sort tools by namespace, resource, verb for consistent output
	sort.Slice(tools, func(i, j int) bool {
		if tools[i].Namespace != tools[j].Namespace {
			return tools[i].Namespace < tools[j].Namespace
		}
		if tools[i].Resource != tools[j].Resource {
			return tools[i].Resource < tools[j].Resource
		}
		return tools[i].Verb < tools[j].Verb
	})

	return tools
}

// toolInfo represents information about an MCP tool
type toolInfo struct {
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	Verb      string `json:"verb"`
	ToolName  string `json:"tool_name"`
	Short     string `json:"short"`
}

// ListResources returns a list of available MCP resources based on the server's configuration.
// Resources are read-only endpoints for list commands that can be accessed via URI.
func (s *MCPServer) ListResources() []resourceInfo {
	var resources []resourceInfo
	for _, cmd := range s.commands {
		// Only list commands are exposed as resources
		if IsListCommand(cmd.Command) &&
			ShouldRegisterCommand(
				cmd.Command,
				s.readOnly,
				s.enabledNamespaces,
				s.enabledResources,
				s.enabledVerbs,
			) {
			resources = append(resources, resourceInfo{
				Namespace: cmd.Command.Namespace,
				Resource:  cmd.Command.Resource,
				URI:       BuildResourceURI(cmd.Command.Namespace, cmd.Command.Resource),
				Short:     cmd.Command.Short,
			})
		}
	}

	// Sort resources by namespace, resource for consistent output
	sort.Slice(resources, func(i, j int) bool {
		if resources[i].Namespace != resources[j].Namespace {
			return resources[i].Namespace < resources[j].Namespace
		}

		return resources[i].Resource < resources[j].Resource
	})

	return resources
}

// resourceInfo represents information about an MCP resource
type resourceInfo struct {
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	URI       string `json:"uri"`
	Short     string `json:"short"`
}

// Serve runs the MCP server with the specified transport.
// It handles graceful shutdown and transport selection.
func (s *MCPServer) Serve(
	ctx context.Context,
	transportMode string,
	address string,
) (any, error) {
	// Log transport information
	fmt.Fprintf(os.Stderr, "Running MCP server with %s transport", transportMode)
	if transportMode != "stdio" {
		fmt.Fprintf(os.Stderr, " on %s", address)
	}
	fmt.Fprintf(os.Stderr, "\n")

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Run server based on transport mode
	switch transportMode {
	case "stdio":
		if err := s.Run(ctx, &mcp.StdioTransport{}); err != nil {
			return nil, fmt.Errorf("MCP server error: %w", err)
		}

	case "sse":
		return nil, RunSSEServer(ctx, s, address)

	case "streamable-http":
		return nil, RunStreamableHTTPServer(ctx, s, address)

	default:
		return nil, fmt.Errorf(
			"unknown transport mode: %s (valid modes: stdio, sse, streamable-http)",
			transportMode,
		)
	}

	return map[string]string{"status": "shutdown"}, nil
}
