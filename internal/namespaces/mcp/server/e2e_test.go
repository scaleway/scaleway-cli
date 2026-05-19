package server_test

import (
	"fmt"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func runServer(t *testing.T) int {
	t.Helper()
	// Find available port on loopback interface only
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	filteredCommands := server.FilterCommands(
		instance.GetCommands().GetSortedCommand(),
		server.CommandFilterConfig{},
	)
	s := server.NewMCPServer(filteredCommands, core.BuildInfo{})

	// Create the streamable HTTP handler.
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return s.Server()
	}, &mcp.StreamableHTTPOptions{
		DisableLocalhostProtection: true,
		JSONResponse:               true,
		Stateless:                  true,
	})

	httpServer := &http.Server{
		Handler: handler,
	}

	// Cleanup: shutdown server and close listener.
	t.Cleanup(func() {
		if err := httpServer.Shutdown(t.Context()); err != nil {
			t.Logf("Server shutdown error: %v", err)
		}
		listener.Close()
	})

	t.Logf("MCP server listening on port %d", port)

	// Start the HTTP server.
	go func() {
		if err := httpServer.Serve(listener); err != http.ErrServerClosed {
			t.Logf("Server error: %v", err)
		}
	}()

	return port
}

func runClient(t *testing.T, port int) {
	t.Helper()
	// Connect to the proxy server (acting as a client).
	client := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1.0.0"}, nil)
	clientSession, err := client.Connect(t.Context(), &mcp.StreamableClientTransport{
		Endpoint: fmt.Sprintf("http://localhost:%d/mcp", port),
	}, nil)
	if err != nil {
		t.Fatalf("Failed to connect: %v", err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(t.Context(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("Failed to call list tools: %v", err)
	}
	t.Logf("Client received result: %v", tools.Tools)

	if len(tools.Tools) == 0 {
		t.Fatal("Expected non-empty tools list")
	}

	resources, err := clientSession.ListResources(t.Context(), &mcp.ListResourcesParams{})
	if err != nil {
		t.Fatalf("Failed to call list resources: %v", err)
	}
	t.Logf("Client received result: %v", resources.Resources)

	if len(resources.Resources) == 0 {
		t.Fatal("Expected non-empty resources list")
	}
}

func Test_E2E(t *testing.T) {
	port := runServer(t)
	// Give the backend a moment to start.
	time.Sleep(100 * time.Millisecond)

	runClient(t, port)
}
