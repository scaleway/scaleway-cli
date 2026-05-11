package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RunStreamableHTTPServer(
	ctx context.Context,
	mcpServer *MCPServer,
	address string,
) error {
	// Create the streamable HTTP handler, which manages sessions and transports
	handler := mcp.NewStreamableHTTPHandler(
		func(req *http.Request) *mcp.Server {
			return mcpServer.Server()
		},
		&mcp.StreamableHTTPOptions{
			JSONResponse: true,
			Stateless:    true,
		},
	)

	server := &http.Server{
		Addr:    address,
		Handler: handler,
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
