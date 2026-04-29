package server

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
)

// MCPServer wraps the MCP server with CLI command integration
type MCPServer struct {
	server       *mcp.Server
	commands     []*CommandTool
	resources    []*CommandResource
	filterConfig CommandFilterConfig
}

// NewMCPServer creates a new MCP server that exposes CLI commands as tools and resources
func NewMCPServer(
	version string,
	cliCommands []*core.Command,
	filterConfig CommandFilterConfig,
) *MCPServer {
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:       "scaleway-mcp",
		Title:      "Scaleway MCP Server",
		WebsiteURL: "https://cli.scaleway.com",
		Version:    version,
	}, &mcp.ServerOptions{
		// Enable tools and resources capabilities with listChanged notifications
		Capabilities: &mcp.ServerCapabilities{
			Tools:     &mcp.ToolCapabilities{ListChanged: true},
			Resources: &mcp.ResourceCapabilities{ListChanged: true},
		},
	})

	s := &MCPServer{
		server:       mcpServer,
		commands:     make([]*CommandTool, 0, len(cliCommands)),
		resources:    make([]*CommandResource, 0),
		filterConfig: filterConfig,
	}

	// Register all commands during initialization
	for _, cmd := range cliCommands {
		if err := s.LoadCommand(cmd); err != nil {
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
