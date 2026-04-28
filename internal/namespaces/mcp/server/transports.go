package server

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func RunSSEServer(ctx context.Context, mcpServer *MCPServer, address string) error {
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
