package server_test

import (
	"context"
	"encoding/json"
	"reflect"
	"strings"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func TestCommandNameToToolName(t *testing.T) {
	tests := []struct {
		name     string
		command  *core.Command
		expected string
	}{
		{
			name: "namespace only",
			command: &core.Command{
				Namespace: "config",
			},
			expected: "config",
		},
		{
			name: "namespace and resource",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
			},
			expected: "instance_server",
		},
		{
			name: "full command",
			command: &core.Command{
				Namespace: "instance",
				Resource:  "server",
				Verb:      "list",
			},
			expected: "instance_server_list",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := server.CommandNameToToolName(tt.command)
			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

// createTestContextWithMeta creates a context with a properly initialized meta and client
// This is needed to avoid config file loading errors in CI environments
func createTestContextWithMeta(t *testing.T) context.Context {
	t.Helper()
	client, err := scw.NewClient(
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultProjectID("11111111-1111-1111-1111-111111111111"),
		scw.WithUserAgent("cli-test"),
	)
	require.NoError(t, err)

	return core.InjectMeta(context.Background(), &core.Meta{
		Client:      client,
		OverrideEnv: map[string]string{},
		BinaryName:  "scw-test",
	})
}

func TestCommandToolExecute(t *testing.T) {
	type testArgs struct {
		Name  string `json:"name"`
		Value int    `json:"value"`
	}

	executedArgs := &testArgs{}
	cmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "get",
		ArgsType:  reflect.TypeOf(testArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "name",
				Short:      "Name parameter",
				Required:   true,
				Positional: false,
			},
			{
				Name:       "value",
				Short:      "Value parameter",
				Required:   false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			executedArgs = args.(*testArgs)

			return map[string]string{"status": "ok"}, nil
		},
	}

	tool := server.NewCommandTool(cmd)

	inputArgs := map[string]any{
		"name":  "test-name",
		"value": float64(42),
	}

	ctx := createTestContextWithMeta(t)
	result, err := tool.Execute(ctx, inputArgs)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if len(result.Content) == 0 {
		t.Fatal("Expected content in result")
	}

	if executedArgs.Name != "test-name" {
		t.Errorf("Expected name 'test-name', got '%s'", executedArgs.Name)
	}
	if executedArgs.Value != 42 {
		t.Errorf("Expected value 42, got %d", executedArgs.Value)
	}
}

func TestCommandToolExecuteWithKebabCase(t *testing.T) {
	type testArgs struct {
		ProjectID string `json:"project_id"`
		Zone      string `json:"zone"`
	}

	executedArgs := &testArgs{}
	cmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "list",
		ArgsType:  reflect.TypeOf(testArgs{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "project-id",
				Short:      "Project ID",
				Required:   true,
				Positional: false,
			},
			{
				Name:       "zone",
				Short:      "Zone",
				Required:   false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, args any) (i any, e error) {
			executedArgs = args.(*testArgs)

			return map[string]string{"status": "ok"}, nil
		},
	}

	tool := server.NewCommandTool(cmd)

	// MCP sends kebab-case keys
	inputArgs := map[string]any{
		"project-id": "12345",
		"zone":       "fr-par-1",
	}

	ctx := createTestContextWithMeta(t)
	result, err := tool.Execute(ctx, inputArgs)
	if err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	if executedArgs.ProjectID != "12345" {
		t.Errorf("Expected project_id '12345', got '%s'", executedArgs.ProjectID)
	}
	if executedArgs.Zone != "fr-par-1" {
		t.Errorf("Expected zone 'fr-par-1', got '%s'", executedArgs.Zone)
	}

	_ = result
}

// TestToolMetaSerialization verifies that the Meta field is properly
// serialized to JSON as _meta when tools are listed.
func TestToolMetaSerialization(t *testing.T) {
	allCommands := commands.GetCommands(context.Background()).GetAll()

	filteredCommands := server.FilterCommands(allCommands, server.CommandFilterConfig{})
	mcpServer := server.NewMCPServer(filteredCommands, core.BuildInfo{})
	registeredCommands := mcpServer.RegisteredCommands()

	if len(registeredCommands) == 0 {
		t.Fatal("No commands registered")
	}

	// Test the first command
	tool := registeredCommands[0]
	mcpTool := tool.ToMCPTool()

	// Marshal to JSON to see what the client actually receives
	jsonBytes, err := json.MarshalIndent(mcpTool, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal tool: %v", err)
	}

	t.Logf("Tool name: %s", mcpTool.Name)
	t.Logf("JSON representation:\n%s", string(jsonBytes))

	// Verify _meta field exists in JSON
	var rawJSON map[string]any
	if err := json.Unmarshal(jsonBytes, &rawJSON); err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	meta, ok := rawJSON["_meta"]
	if !ok {
		t.Error("_meta field not present in JSON")
	} else {
		t.Logf("_meta content: %v", meta)

		// Check for expected keys
		metaMap, ok := meta.(map[string]any)
		if !ok {
			t.Error("_meta is not a map")
		} else {
			if _, ok := metaMap["namespace"]; !ok {
				t.Error("_meta.namespace missing")
			}
			if _, ok := metaMap["resource"]; !ok {
				t.Error("_meta.resource missing")
			}
			if _, ok := metaMap["verb"]; !ok {
				t.Error("_meta.verb missing")
			}
		}
	}
}

// TestCommandToolExecutePanicRecovery verifies that panics during command
// execution are recovered and returned as proper error responses.
func TestCommandToolExecutePanicRecovery(t *testing.T) {
	cmd := &core.Command{
		Namespace: "test",
		Resource:  "panic",
		Verb:      "trigger",
		ArgsType:  nil,
		Run: func(ctx context.Context, args any) (i any, e error) {
			// Simulate a panic like the one in instanceServerList
			panic("nil pointer dereference")
		},
	}

	tool := server.NewCommandTool(cmd)

	inputArgs := map[string]any{}

	ctx := createTestContextWithMeta(t)
	result, err := tool.Execute(ctx, inputArgs)
	// Should not return an error (panic is recovered)
	if err != nil {
		t.Fatalf("Execute should not return error, got: %v", err)
	}

	// Should have content
	if len(result.Content) == 0 {
		t.Fatal("Expected content in result")
	}

	// Should be marked as error
	if !result.IsError {
		t.Error("Expected IsError to be true")
	}

	// Content should mention panic recovery
	tc, ok := result.Content[0].(*mcp.TextContent)
	if !ok {
		t.Fatalf("Expected TextContent, got %T", result.Content[0])
	}

	if !strings.Contains(tc.Text, "panic recovered") {
		t.Errorf("Expected panic recovery message, got: %s", tc.Text)
	}

	t.Logf("Recovered panic message: %s", tc.Text)
}
