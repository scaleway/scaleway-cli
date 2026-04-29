package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"slices"
	"strconv"
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

// RegisterResource registers a CLI command as an MCP resource
func (s *MCPServer) RegisterResource(cmd *core.Command) error {
	if !ShouldRegisterCommand(
		cmd,
		s.readOnly,
		s.enabledNamespaces,
		s.enabledResources,
		s.enabledVerbs,
	) {
		return nil
	}

	resource := NewCommandResource(cmd)
	mcpResource := resource.ToMCPResource()

	// Create a handler function for the resource
	handler := func(ctx context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
		// Extract arguments from the request URI
		// URI format: scw://namespace/resource?arg1=value1&arg2=value2
		inputArgs := parseURIToArgs(req.Params.URI)

		result, err := resource.Execute(ctx, inputArgs)
		if err != nil {
			return nil, err
		}

		return result, nil
	}

	// Register with MCP SDK
	s.server.AddResource(mcpResource, handler)

	s.resources = append(s.resources, resource)

	return nil
}

// parseURIToArgs extracts query parameters from a URI and converts them to input args
func parseURIToArgs(uri string) map[string]any {
	args := make(map[string]any)

	// Parse URI query parameters
	// Format: scw://namespace/resource?key1=value1&key2=value2
	parts := strings.SplitN(uri, "?", 2)
	if len(parts) != 2 {
		return args
	}

	query := parts[1]
	paramPairs := strings.SplitSeq(query, "&")

	for pair := range paramPairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 {
			continue
		}

		key := kv[0]
		value := kv[1]

		// Try to parse the value as different types
		if value == "true" {
			args[key] = true
		} else if value == "false" {
			args[key] = false
		} else if num, err := strconv.ParseFloat(value, 64); err == nil {
			args[key] = num
		} else {
			args[key] = value
		}
	}

	return args
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

func RunMCPServer(
	ctx context.Context,
	transportMode string,
	address string,
	readOnly bool,
	enabledNamespaces []string,
	enabledResources []string,
	enabledVerbs []string,
) (any, error) {
	// Get all CLI commands from the meta context
	commands := core.ExtractCommands(ctx)
	cliCommands := commands.GetAll()

	// Get build info for version
	buildInfo := core.ExtractBuildInfo(ctx)
	version := buildInfo.Version.String()

	// Get profile from context (set by global --profile flag)
	profile := core.ExtractProfileName(ctx)
	configPath := core.ExtractConfigPath(ctx)

	// Log startup information to stderr
	fmt.Fprintf(os.Stderr, "Starting MCP server version %s\n", version)
	fmt.Fprintf(os.Stderr, "Transport mode: %s\n", transportMode)
	fmt.Fprintf(os.Stderr, "Read-only mode: %v\n", readOnly)
	if len(enabledNamespaces) > 0 {
		fmt.Fprintf(os.Stderr, "Enabled namespaces: %v\n", enabledNamespaces)
	}
	if len(enabledResources) > 0 {
		fmt.Fprintf(os.Stderr, "Enabled resources: %v\n", enabledResources)
	}
	if len(enabledVerbs) > 0 {
		fmt.Fprintf(os.Stderr, "Enabled verbs: %v\n", enabledVerbs)
	}
	fmt.Fprintf(os.Stderr, "Using profile: %s\n", profile)
	fmt.Fprintf(os.Stderr, "Config path: %s\n", configPath)

	// Reload the client with the profile to ensure proper authentication
	// This is necessary because the bootstrap creates an anonymous client for
	// commands with AllowAnonymousClient: true
	if err := core.ReloadClient(ctx); err != nil {
		return nil, fmt.Errorf("failed to initialize authenticated client: %w", err)
	}

	// Verify client is properly initialized
	client := core.ExtractClient(ctx)
	if client != nil {
		if orgID, ok := client.GetDefaultOrganizationID(); ok {
			fmt.Fprintf(os.Stderr, "Organization ID: %s\n", orgID)
		}
		if projectID, ok := client.GetDefaultProjectID(); ok {
			fmt.Fprintf(os.Stderr, "Project ID: %s\n", projectID)
		}
	} else {
		fmt.Fprintf(os.Stderr, "Warning: No client initialized\n")
	}

	// Create MCP server with all commands
	mcpServer := NewMCPServer(
		version,
		cliCommands,
		readOnly,
		enabledNamespaces,
		enabledResources,
		enabledVerbs,
	)

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
		fmt.Fprintf(os.Stderr, "Running MCP server over stdio\n")
		if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
			return nil, fmt.Errorf("MCP server error: %w", err)
		}

	case "sse":
		fmt.Fprintf(os.Stderr, "Running MCP server with SSE transport on %s\n", address)

		return nil, RunSSEServer(ctx, mcpServer, address)

	case "streamable-http":
		fmt.Fprintf(os.Stderr, "Running MCP server with streamable HTTP transport on %s\n", address)

		return nil, RunStreamableHTTPServer(ctx, mcpServer, address)

	default:
		return nil, fmt.Errorf(
			"unknown transport mode: %s (valid modes: stdio, sse, streamable-http)",
			transportMode,
		)
	}

	return map[string]string{"status": "shutdown"}, nil
}
