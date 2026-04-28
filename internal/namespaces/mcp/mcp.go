package mcp

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		mcpRoot(),
		mcpServerCommand(),
	)
}

func mcpRoot() *core.Command {
	return &core.Command{
		Groups:    []string{"utility"},
		Namespace: "mcp",
		Short:     "MCP (Model Context Protocol) server management",
		Long:      "Commands for managing the MCP server that exposes Scaleway CLI commands as AI tools.",
	}
}

func mcpServerCommand() *core.Command {
	type serveArgs struct {
		Transport string `json:"transport"`
		Address   string `json:"address"`
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "serve",
		Short:                "Start the MCP server",
		Long:                 "Runs the MCP server, exposing all CLI commands as MCP tools for AI assistants. Supports stdio (default), SSE, and streamable HTTP transports.",
		AllowAnonymousClient: true,
		DisableTelemetry:     true,
		ArgsType:             reflect.TypeOf(serveArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "transport",
				Short:      "Transport mode: stdio (default), sse, or streamable-http",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("stdio"),
			},
			{
				Name:       "address",
				Short:      "Address to bind for SSE and streamable-http transports (e.g., :8080)",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter(":8080"),
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*serveArgs)

			return runMCPServer(ctx, args.Transport, args.Address)
		},
	}
}

func runMCPServer(ctx context.Context, transportMode string, address string) (any, error) {
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
	mcpServer := server.NewMCPServer(version, cliCommands)

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

		return nil, runSSEServer(ctx, mcpServer, address)

	case "streamable-http":
		fmt.Fprintf(os.Stderr, "Running MCP server with streamable HTTP transport on %s\n", address)

		return nil, runStreamableHTTPServer(ctx, mcpServer, address)

	default:
		return nil, fmt.Errorf(
			"unknown transport mode: %s (valid modes: stdio, sse, streamable-http)",
			transportMode,
		)
	}

	return map[string]string{"status": "shutdown"}, nil
}

func runSSEServer(ctx context.Context, mcpServer *server.MCPServer, address string) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		// Create SSE transport for this connection
		transport := &mcp.SSEServerTransport{
			Endpoint: "/message",
			Response: w,
		}

		// Connect the server to this transport
		session, err := mcpServer.Server().Connect(ctx, transport, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Connection error: %v", err), http.StatusInternalServerError)

			return
		}

		// Handle messages until session ends
		_ = session.Wait()
	})

	mux.HandleFunc("/message", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)

			return
		}

		// Handle incoming messages
		// Note: This is a simplified handler - in production you'd want
		// to track sessions and route messages appropriately
		w.WriteHeader(http.StatusAccepted)
	})

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
		}
	}()

	fmt.Fprintf(os.Stderr, "SSE server listening on %s\n", address)
	fmt.Fprintf(os.Stderr, "Connect to: http://%s/sse\n", address)

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	return server.Shutdown(shutdownCtx)
}

func runStreamableHTTPServer(
	ctx context.Context,
	mcpServer *server.MCPServer,
	address string,
) error {
	mux := http.NewServeMux()

	mux.HandleFunc("/mcp", func(w http.ResponseWriter, r *http.Request) {
		// Create streamable transport for this connection
		transport := &mcp.StreamableServerTransport{
			SessionID: "", // Empty for stateless mode
		}

		// Connect the server to this transport
		session, err := mcpServer.Server().Connect(ctx, transport, nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("Connection error: %v", err), http.StatusInternalServerError)

			return
		}

		// Handle messages until session ends
		_ = session.Wait()
	})

	server := &http.Server{
		Addr:    address,
		Handler: mux,
	}

	// Start server in goroutine
	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Fprintf(os.Stderr, "HTTP server error: %v\n", err)
		}
	}()

	fmt.Fprintf(os.Stderr, "Streamable HTTP server listening on %s\n", address)
	fmt.Fprintf(os.Stderr, "Connect to: http://%s/mcp\n", address)

	// Wait for shutdown signal
	<-ctx.Done()

	// Graceful shutdown with timeout
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	return server.Shutdown(shutdownCtx)
}
