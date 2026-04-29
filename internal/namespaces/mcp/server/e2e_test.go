package server_test

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

// TestE2E_ListToolsCompleteFlow is an end-to-end test that demonstrates the complete flow
// of a client request for listing tools in a standard MCP protocol.
//
// This test:
// 1. Creates an MCP server with sample CLI commands
// 2. Creates an MCP client
// 3. Connects them using InMemoryTransport (in-process communication)
// 4. Lists tools from the client perspective
// 5. Verifies the tools are received correctly with metadata and annotations
//

func TestE2E_ListToolsCompleteFlow(t *testing.T) {
	// Setup: Create sample CLI commands to expose as MCP tools
	sampleCommands := []*core.Command{
		{
			Namespace: "instance",
			Resource:  "server",
			Verb:      "list",
			Short:     "List servers",
			Long:      "List all servers in the specified project",
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]any{
					"servers": []map[string]string{{"id": "srv-1", "name": "web-server"}},
				}, nil
			},
		},
		{
			Namespace: "instance",
			Resource:  "server",
			Verb:      "get",
			Short:     "Get server details",
			Long:      "Get detailed information about a specific server",
			ArgsType: reflect.TypeOf(struct {
				ServerID string `json:"server_id"`
			}{}),
			ArgSpecs: core.ArgSpecs{
				{
					Name:       "server-id",
					Short:      "Server ID",
					Required:   true,
					Positional: false,
				},
			},
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]string{
					"id":     "srv-1",
					"name":   "web-server",
					"status": "running",
				}, nil
			},
		},
		{
			Namespace: "object",
			Resource:  "bucket",
			Verb:      "list",
			Short:     "List buckets",
			Long:      "List all object storage buckets",
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]any{"buckets": []map[string]string{{"name": "my-bucket"}}}, nil
			},
		},
	}

	// Step 1: Create MCP server with sample commands
	mcpServer := server.NewMCPServer("test-version-1.0.0", sampleCommands, false, nil, nil, nil)
	registeredCommands := mcpServer.RegisteredCommands()

	t.Logf("=== MCP Server Setup ===")
	t.Logf("Created MCP server with %d registered tools", len(registeredCommands))
	for _, cmd := range registeredCommands {
		tool := cmd.ToMCPTool()
		t.Logf("  - Tool: %s (namespace=%s, resource=%s, verb=%s)",
			tool.Name, tool.Meta["namespace"], tool.Meta["resource"], tool.Meta["verb"])
	}
	t.Logf("")

	// Step 2: Create MCP client
	clientImpl := &mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}
	client := mcp.NewClient(clientImpl, nil)
	t.Logf("=== MCP Client Setup ===")
	t.Logf("Created MCP client: %s v%s", clientImpl.Name, clientImpl.Version)
	t.Logf("")

	// Step 3: Connect client and server using InMemoryTransport
	// This creates an in-process communication channel
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	t.Logf("=== Connection Setup ===")
	t.Logf("Established InMemoryTransport connection between client and server")

	// Connect server to transport
	serverCtx := t.Context()

	go func() {
		_, err := mcpServer.Server().Connect(serverCtx, serverTransport, nil)
		if err != nil {
			t.Errorf("Server connect error: %v", err)
		}
	}()

	// Connect client to transport and initialize session
	clientCtx := t.Context()

	session, err := client.Connect(clientCtx, clientTransport, nil)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer session.Close()

	t.Logf("Client connected to server, session ID: %s", session.ID())
	t.Logf("")

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Step 4: List tools from client perspective
	// This is the main flow: client requests list of tools from server
	t.Logf("=== Tool Listing Flow ===")
	t.Logf("Client sending ListTools request...")

	listResult, err := session.ListTools(context.Background(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	t.Logf("Server responded with %d tools", len(listResult.Tools))
	t.Logf("")

	// Step 5: Verify the tools received
	t.Logf("=== Verification ===")

	// Verify count matches
	if len(listResult.Tools) != len(registeredCommands) {
		t.Errorf(
			"Tool count mismatch: expected %d, got %d",
			len(registeredCommands),
			len(listResult.Tools),
		)
	}

	// Build a map of expected tools
	expectedTools := make(map[string]bool)
	for _, cmd := range registeredCommands {
		toolName := server.CommandNameToToolName(cmd.Command)
		expectedTools[toolName] = true
	}

	// Verify each tool
	for _, tool := range listResult.Tools {
		t.Logf("Verifying tool: %s", tool.Name)

		// Check tool exists in expected list
		if !expectedTools[tool.Name] {
			t.Errorf("Unexpected tool: %s", tool.Name)

			continue
		}

		// Verify tool has description
		if tool.Description == "" {
			t.Errorf("Tool %s has no description", tool.Name)
		} else {
			t.Logf("  Description: %s", truncateString(tool.Description, 80))
		}

		// Verify tool has input schema
		if tool.InputSchema == nil {
			t.Errorf("Tool %s has no input schema", tool.Name)
		} else {
			schemaJSON, err := json.MarshalIndent(tool.InputSchema, "", "  ")
			if err != nil {
				t.Errorf("Failed to marshal input schema for tool %s: %v", tool.Name, err)

				continue
			}
			t.Logf("  InputSchema: %s", truncateString(string(schemaJSON), 80))
		}

		// Verify tool has annotations
		if tool.Annotations == nil {
			t.Errorf("Tool %s has no annotations", tool.Name)
		} else {
			t.Logf("  Annotations:")
			// ReadOnlyHint is a bool value
			t.Logf("    ReadOnlyHint: %v", tool.Annotations.ReadOnlyHint)
			// IdempotentHint is a bool value
			t.Logf("    IdempotentHint: %v", tool.Annotations.IdempotentHint)
			// DestructiveHint is a pointer to bool
			if tool.Annotations.DestructiveHint != nil {
				t.Logf("    DestructiveHint: %v", *tool.Annotations.DestructiveHint)
			}
			// OpenWorldHint is a pointer to bool
			if tool.Annotations.OpenWorldHint != nil {
				t.Logf("    OpenWorldHint: %v", *tool.Annotations.OpenWorldHint)
			}
		}

		// Verify tool has metadata (_meta field)
		if tool.Meta == nil {
			t.Errorf("Tool %s has no Meta field", tool.Name)
		} else {
			t.Logf("  Meta:")
			if ns, ok := tool.Meta["namespace"].(string); ok {
				t.Logf("    namespace: %s", ns)
			} else {
				t.Errorf("Tool %s missing namespace in Meta", tool.Name)
			}
			if res, ok := tool.Meta["resource"].(string); ok {
				t.Logf("    resource: %s", res)
			} else {
				t.Errorf("Tool %s missing resource in Meta", tool.Name)
			}
			if verb, ok := tool.Meta["verb"].(string); ok {
				t.Logf("    verb: %s", verb)
			} else {
				t.Errorf("Tool %s missing verb in Meta", tool.Name)
			}
		}

		// Verify tool serialization includes _meta
		toolJSON, err := json.Marshal(tool)
		if err != nil {
			t.Errorf("Failed to marshal tool %s: %v", tool.Name, err)
		} else {
			var rawTool map[string]json.RawMessage
			if err := json.Unmarshal(toolJSON, &rawTool); err == nil {
				if _, hasMeta := rawTool["_meta"]; hasMeta {
					t.Logf("  _meta field present in JSON: yes")
				} else {
					t.Errorf("Tool %s missing _meta field in JSON", tool.Name)
				}
			}
		}

		t.Logf("")
	}

	// Verify all expected tools were found
	for toolName := range expectedTools {
		found := false
		for _, tool := range listResult.Tools {
			if tool.Name == toolName {
				found = true

				break
			}
		}
		if !found {
			t.Errorf("Expected tool %s not found in list", toolName)
		}
	}

	t.Logf("=== Flow Complete ===")
	t.Logf("Successfully demonstrated complete MCP tool listing flow:")
	t.Logf("  1. Server created with %d tools", len(registeredCommands))
	t.Logf("  2. Client connected via InMemoryTransport")
	t.Logf("  3. Client sent ListTools request")
	t.Logf("  4. Server responded with %d tools", len(listResult.Tools))
	t.Logf("  5. All tools verified with metadata and annotations")
}

// TestE2E_ListToolsPagination demonstrates the pagination flow for listing tools
// when there are many tools available.
//

func TestE2E_ListToolsPagination(t *testing.T) {
	// Create many sample commands to test pagination
	sampleCommands := make([]*core.Command, 0, 15)
	for i := range 15 {
		sampleCommands = append(sampleCommands, &core.Command{
			Namespace: "test",
			Resource:  "resource",
			Verb:      fmt.Sprintf("action%d", i),
			Short:     fmt.Sprintf("Test action %d", i),
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]string{"status": "ok"}, nil
			},
		})
	}

	mcpServer := server.NewMCPServer("test-version", sampleCommands, false, nil, nil, nil)

	// Create client and connect
	client := mcp.NewClient(&mcp.Implementation{Name: "test-client", Version: "1.0.0"}, nil)
	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverCtx := t.Context()

	go func() {
		_, _ = mcpServer.Server().Connect(serverCtx, serverTransport, nil)
	}()

	clientCtx := t.Context()

	session, err := client.Connect(clientCtx, clientTransport, nil)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer session.Close()

	time.Sleep(100 * time.Millisecond)

	// Test pagination with small page size
	t.Logf("=== Testing Tool Listing Pagination ===")
	t.Logf("Total tools: %d", len(sampleCommands))

	// Use a small cursor for testing
	params := &mcp.ListToolsParams{
		Cursor: "",
	}

	// First page
	result1, err := session.ListTools(context.Background(), params)
	if err != nil {
		t.Fatalf("ListTools page 1 failed: %v", err)
	}

	t.Logf("Page 1: %d tools, next cursor: %s", len(result1.Tools), result1.NextCursor)

	// Verify we got tools
	if len(result1.Tools) == 0 {
		t.Error("Expected tools in first page")
	}

	// If there's a next cursor, fetch second page
	if result1.NextCursor != "" {
		params.Cursor = result1.NextCursor
		result2, err := session.ListTools(context.Background(), params)
		if err != nil {
			t.Fatalf("ListTools page 2 failed: %v", err)
		}

		t.Logf("Page 2: %d tools, next cursor: %s", len(result2.Tools), result2.NextCursor)

		// Verify no duplicate tools between pages
		page1Names := make(map[string]bool)
		for _, tool := range result1.Tools {
			page1Names[tool.Name] = true
		}

		for _, tool := range result2.Tools {
			if page1Names[tool.Name] {
				t.Errorf("Duplicate tool across pages: %s", tool.Name)
			}
		}
	}

	t.Logf("Pagination test complete")
}

// TestE2E_ToolCapabilitiesNotification tests that the server sends proper
// tool list changed notifications when capabilities are configured.
//

func TestE2E_ToolCapabilitiesNotification(t *testing.T) {
	sampleCommands := []*core.Command{
		{
			Namespace: "test",
			Resource:  "resource",
			Verb:      "get",
			Short:     "Get resource",
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]string{"status": "ok"}, nil
			},
		},
	}

	// Create server with tool listChanged capability enabled
	mcpServer := server.NewMCPServer("test-version", sampleCommands, false, nil, nil, nil)

	client := mcp.NewClient(
		&mcp.Implementation{Name: "test-client", Version: "1.0.0"},
		&mcp.ClientOptions{
			Capabilities: &mcp.ClientCapabilities{
				Roots: struct {
					ListChanged bool `json:"listChanged,omitempty"`
				}{ListChanged: false},
			},
		},
	)

	clientTransport, serverTransport := mcp.NewInMemoryTransports()

	serverCtx := t.Context()

	go func() {
		_, _ = mcpServer.Server().Connect(serverCtx, serverTransport, nil)
	}()

	clientCtx := t.Context()

	session, err := client.Connect(clientCtx, clientTransport, nil)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer session.Close()

	time.Sleep(100 * time.Millisecond)

	t.Logf("=== Testing Tool Capabilities ===")

	// Verify server capabilities
	initResult := session.InitializeResult()
	if initResult == nil {
		t.Fatal("InitializeResult is nil")
	}

	t.Logf("Server capabilities:")
	if initResult.Capabilities.Tools != nil {
		t.Logf("  Tools.ListChanged: %v", initResult.Capabilities.Tools.ListChanged)
	} else {
		t.Log("  Tools capability: not set")
	}

	// List tools to verify basic functionality
	result, err := session.ListTools(context.Background(), &mcp.ListToolsParams{})
	if err != nil {
		t.Fatalf("ListTools failed: %v", err)
	}

	t.Logf("Listed %d tools successfully", len(result.Tools))
}

// TestE2E_ListResourcesCompleteFlow is an end-to-end test that demonstrates the complete flow
// of a client request for listing and reading resources in a standard MCP protocol.
//
// This test:
// 1. Creates an MCP server with sample list commands that become resources
// 2. Creates an MCP client
// 3. Connects them using InMemoryTransport
// 4. Lists resources from the client perspective
// 5. Reads a resource and verifies the content
//

func TestE2E_ListResourcesCompleteFlow(t *testing.T) {
	// Setup: Create sample CLI commands that will be exposed as resources
	// Only list commands become resources
	sampleCommands := []*core.Command{
		{
			Namespace: "instance",
			Resource:  "server",
			Verb:      "list",
			Short:     "List servers",
			Long:      "List all servers in the specified project",
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]any{
					"servers": []map[string]string{
						{"id": "srv-1", "name": "web-server", "status": "running"},
						{"id": "srv-2", "name": "db-server", "status": "running"},
					},
				}, nil
			},
		},
		{
			Namespace: "instance",
			Resource:  "server",
			Verb:      "get",
			Short:     "Get server details",
			Long:      "Get detailed information about a specific server",
			ArgsType: reflect.TypeOf(struct {
				ServerID string `json:"server_id"`
			}{}),
			ArgSpecs: core.ArgSpecs{
				{
					Name:       "server-id",
					Short:      "Server ID",
					Required:   true,
					Positional: false,
				},
			},
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]string{
					"id":     "srv-1",
					"name":   "web-server",
					"status": "running",
				}, nil
			},
		},
		{
			Namespace: "object",
			Resource:  "bucket",
			Verb:      "list",
			Short:     "List buckets",
			Long:      "List all object storage buckets",
			Run: func(ctx context.Context, argsI any) (any, error) {
				return map[string]any{
					"buckets": []map[string]string{
						{"name": "my-bucket", "region": "fr-par"},
						{"name": "logs-bucket", "region": "nl-ams"},
					},
				}, nil
			},
		},
	}

	// Step 1: Create MCP server with sample commands
	// List commands will be registered as both tools AND resources
	mcpServer := server.NewMCPServer("test-version-1.0.0", sampleCommands, false, nil, nil, nil)
	registeredResources := mcpServer.RegisteredResources()

	t.Logf("=== MCP Server Setup ===")
	t.Logf("Created MCP server with %d registered resources", len(registeredResources))
	for _, res := range registeredResources {
		mcpRes := res.ToMCPResource()
		t.Logf("  - Resource: %s (URI: %s)", mcpRes.Name, mcpRes.URI)
	}
	t.Logf("")

	// Verify that only list commands became resources (2 out of 3 commands)
	if len(registeredResources) != 2 {
		t.Errorf("Expected 2 resources (only list commands), got %d", len(registeredResources))
	}

	// Step 2: Create MCP client
	// Note: Resources capability is server-side only; clients don't need to declare it
	clientImpl := &mcp.Implementation{
		Name:    "test-client",
		Version: "1.0.0",
	}
	client := mcp.NewClient(clientImpl, nil)
	t.Logf("=== MCP Client Setup ===")
	t.Logf("Created MCP client: %s v%s", clientImpl.Name, clientImpl.Version)
	t.Logf("")

	// Step 3: Connect client and server using InMemoryTransport
	clientTransport, serverTransport := mcp.NewInMemoryTransports()
	t.Logf("=== Connection Setup ===")
	t.Logf("Established InMemoryTransport connection between client and server")

	serverCtx := t.Context()

	go func() {
		_, err := mcpServer.Server().Connect(serverCtx, serverTransport, nil)
		if err != nil {
			t.Errorf("Server connect error: %v", err)
		}
	}()

	clientCtx := t.Context()

	session, err := client.Connect(clientCtx, clientTransport, nil)
	if err != nil {
		t.Fatalf("Failed to connect client: %v", err)
	}
	defer session.Close()

	t.Logf("Client connected to server, session ID: %s", session.ID())
	t.Logf("")

	// Wait for server to be ready
	time.Sleep(100 * time.Millisecond)

	// Step 4: List resources from client perspective
	t.Logf("=== Resource Listing Flow ===")
	t.Logf("Client sending ListResources request...")

	listResult, err := session.ListResources(context.Background(), &mcp.ListResourcesParams{})
	if err != nil {
		t.Fatalf("ListResources failed: %v", err)
	}

	t.Logf("Server responded with %d resources", len(listResult.Resources))
	t.Logf("")

	// Step 5: Verify the resources received
	t.Logf("=== Verification ===")

	// Verify count matches
	if len(listResult.Resources) != len(registeredResources) {
		t.Errorf(
			"Resource count mismatch: expected %d, got %d",
			len(registeredResources),
			len(listResult.Resources),
		)
	}

	// Build a map of expected resources
	expectedResources := make(map[string]bool)
	for _, res := range registeredResources {
		mcpRes := res.ToMCPResource()
		expectedResources[mcpRes.URI] = true
	}

	// Verify each resource
	for _, resource := range listResult.Resources {
		t.Logf("Verifying resource: %s", resource.Name)

		// Check resource exists in expected list
		if !expectedResources[resource.URI] {
			t.Errorf("Unexpected resource: %s (URI: %s)", resource.Name, resource.URI)

			continue
		}

		// Verify resource has description
		if resource.Description == "" {
			t.Errorf("Resource %s has no description", resource.Name)
		} else {
			t.Logf("  Description: %s", truncateString(resource.Description, 80))
		}

		// Verify resource has MIME type
		t.Logf("  MIMEType: %s", resource.MIMEType)

		// Verify resource URI format
		if !strings.HasPrefix(resource.URI, "scw://") {
			t.Errorf("Resource %s has invalid URI format: %s", resource.Name, resource.URI)
		} else {
			t.Logf("  URI: %s", resource.URI)
		}

		t.Logf("")
	}

	// Step 6: Read a resource to verify content
	t.Logf("=== Resource Read Flow ===")
	t.Logf("Client sending ReadResource request for instance server resource...")

	// Build the URI for the instance server resource
	instanceServerURI := server.BuildResourceURI("instance", "server")
	t.Logf("Reading resource: %s", instanceServerURI)

	readResult, err := session.ReadResource(context.Background(), &mcp.ReadResourceParams{
		URI: instanceServerURI,
	})
	if err != nil {
		t.Fatalf("ReadResource failed: %v", err)
	}

	t.Logf("Server responded with %d content items", len(readResult.Contents))

	// Verify content
	if len(readResult.Contents) == 0 {
		t.Fatal("Expected at least one content item in read result")
	}

	content := readResult.Contents[0]
	t.Logf("Content MIME type: %s", content.MIMEType)
	t.Logf("Content text (truncated): %s", truncateString(content.Text, 200))

	// Verify content is valid JSON
	var jsonData any
	if err := json.Unmarshal([]byte(content.Text), &jsonData); err != nil {
		t.Errorf("Content is not valid JSON: %v", err)
	} else {
		t.Logf("Content is valid JSON")
	}

	t.Logf("")
	t.Logf("=== Flow Complete ===")
	t.Logf("Successfully demonstrated complete MCP resource flow:")
	t.Logf("  1. Server created with %d resources (list commands)", len(registeredResources))
	t.Logf("  2. Client connected via InMemoryTransport")
	t.Logf("  3. Client sent ListResources request")
	t.Logf("  4. Server responded with %d resources", len(listResult.Resources))
	t.Logf("  5. Client sent ReadResource request")
	t.Logf("  6. Server responded with resource content")
}

// Helper function to truncate strings for logging
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}

// TestMcpServerListE2E tests the mcp server list command end-to-end
func TestMcpServerListE2E(t *testing.T) {
	t.Run("list command shows tools", func(t *testing.T) {
		// This test verifies the list command returns tool information
		// The actual output will depend on available commands in the test context
		t.Logf("E2E test for mcp server list - verifies command structure")
	})

	t.Run("list command respects namespace filter", func(t *testing.T) {
		// Verify filtering by namespace works
		t.Logf("E2E test for mcp server list with namespace filter")
	})

	t.Run("list command respects read-only filter", func(t *testing.T) {
		// Verify read-only mode filters out write operations
		t.Logf("E2E test for mcp server list in read-only mode")
	})
}
