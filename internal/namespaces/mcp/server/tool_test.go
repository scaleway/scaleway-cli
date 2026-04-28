package server_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
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

	result, err := tool.Execute(context.Background(), inputArgs)
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

	result, err := tool.Execute(context.Background(), inputArgs)
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
