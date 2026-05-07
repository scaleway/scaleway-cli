package server_test

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func runServer(t *testing.T) int {
	t.Helper()
	// Find available port
	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatalf("Failed to find available port: %v", err)
	}
	port := listener.Addr().(*net.TCPAddr).Port

	s := server.NewMCPServer(
		"test",
		instance.GetCommands().GetSortedCommand(),
		server.CommandFilterConfig{},
		nil, // No baseMeta for test
	)

	// Create the streamable HTTP handler.
	handler := mcp.NewStreamableHTTPHandler(func(req *http.Request) *mcp.Server {
		return s.Server()
	}, &mcp.StreamableHTTPOptions{
		DisableLocalhostProtection: true,
		JSONResponse:               true,
		Stateless:                  true,
	})

	log.Printf("MCP server listening on port %d", port)

	// Start the HTTP server.
	go func() {
		if err := http.Serve(listener, handler); err != nil {
			log.Printf("Server error: %v", err)
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
		log.Fatalf("Failed to connect: %v", err)
	}
	defer clientSession.Close()

	tools, err := clientSession.ListTools(t.Context(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("Failed to call list tools: %v", err)
	}
	t.Logf("Client received result: %v", tools.Tools)

	resources, err := clientSession.ListResources(t.Context(), &mcp.ListResourcesParams{})
	if err != nil {
		t.Fatalf("Failed to call list resources: %v", err)
	}
	t.Logf("Client received result: %v", resources.Resources)
}

func Test_E2E(t *testing.T) {
	port := runServer(t)
	// Give the backend a moment to start.
	time.Sleep(100 * time.Millisecond)

	runClient(t, port)
}
