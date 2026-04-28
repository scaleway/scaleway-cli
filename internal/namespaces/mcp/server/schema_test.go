package server_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func TestCommandToFlatArgsSchema(t *testing.T) {
	cmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "list",
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "zone",
				Short:      "Zone to target",
				Required:   true,
				EnumValues: []string{"fr-par-1", "nl-ams-1"},
			},
			{
				Name:     "project-id",
				Short:    "Project ID",
				Required: false,
			},
		},
	}

	schema := server.CommandToFlatArgsSchema(cmd)

	if schema.Type != "object" {
		t.Errorf("Expected type 'object', got '%s'", schema.Type)
	}

	if len(schema.Properties) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(schema.Properties))
	}

	if len(schema.Required) != 1 {
		t.Errorf("Expected 1 required field, got %d", len(schema.Required))
	}

	if schema.Required[0] != "zone" {
		t.Errorf("Expected 'zone' to be required, got '%s'", schema.Required[0])
	}
}

func TestArgSpecToJSONSchema(t *testing.T) {
	argSpec := &core.ArgSpec{
		Name:       "test-arg",
		Short:      "Test argument",
		EnumValues: []string{"value1", "value2"},
	}

	schema := server.ArgSpecToJSONSchema(argSpec)

	if schema.Type != "string" {
		t.Errorf("Expected type 'string', got '%s'", schema.Type)
	}

	if len(schema.Enum) != 2 {
		t.Errorf("Expected 2 enum values, got %d", len(schema.Enum))
	}
}
