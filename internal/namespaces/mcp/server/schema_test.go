package server_test

import (
	"encoding/json"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

const (
	arrayArgName      = "tags.{index}"
	arrayPropertyName = "tags"
	mapArgName        = "environment-variables.{key}"
	mapPropertyName   = "environment-variables"
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

	defaultSchema := server.ArgSpecToJSONSchema(&core.ArgSpec{
		Name:  "default-arg",
		Short: "Default argument",
	})
	if defaultSchema.Type != "string" {
		t.Errorf("Expected default type 'string', got '%s'", defaultSchema.Type)
	}
}

func TestCommandToFlatArgsSchemaDynamicArgs(t *testing.T) {
	cmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "create",
		ArgSpecs: core.ArgSpecs{
			{
				Name:  arrayArgName,
				Short: "Tags",
			},
			{
				Name:  mapArgName,
				Short: "Environment variables",
			},
		},
	}

	// Roundtrip through JSON to pin the wire format: top-level
	// additionalProperties must serialize as false, and map args as an
	// object schema with a string-valued additionalProperties.
	rawSchema, err := json.Marshal(server.CommandToFlatArgsSchema(cmd))
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rawSchema, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	if decoded["additionalProperties"] != false {
		t.Fatalf("Expected additionalProperties to be false, got %v", decoded["additionalProperties"])
	}

	properties := decoded["properties"].(map[string]any)
	if _, ok := properties[arrayArgName]; ok {
		t.Fatalf("Schema should not expose literal placeholder property %q", arrayArgName)
	}
	if _, ok := properties[mapArgName]; ok {
		t.Fatalf("Schema should not expose literal placeholder property %q", mapArgName)
	}

	tags := properties[arrayPropertyName].(map[string]any)
	if tags["type"] != "array" {
		t.Fatalf("Expected %q to be an array schema, got %v", arrayPropertyName, tags["type"])
	}
	if tagItems := tags["items"].(map[string]any); tagItems["type"] != "string" {
		t.Fatalf("Expected %q items to be strings, got %v", arrayPropertyName, tagItems["type"])
	}

	environmentVariables := properties[mapPropertyName].(map[string]any)
	if environmentVariables["type"] != "object" {
		t.Fatalf("Expected %q to be an object schema, got %v", mapPropertyName, environmentVariables["type"])
	}
	if additionalProperties := environmentVariables["additionalProperties"].(map[string]any); additionalProperties["type"] != "string" {
		t.Fatalf("Expected %q values to be strings, got %v", mapPropertyName, additionalProperties["type"])
	}
}

func TestCommandToFlatArgsSchemaNestedDynamicArgs(t *testing.T) {
	// Arg specs with nested placeholders (e.g., "pools.{index}.kubelet-args.{key}")
	// must not be treated as simple maps or arrays. They fall through to the
	// default string case, exposing the literal placeholder name as the property.
	nestedMapArgName := "pools.{index}.kubelet-args.{key}"
	nestedArrayArgName := "pools.{index}.tags.{index}"

	cmd := &core.Command{
		Namespace: "test",
		Resource:  "resource",
		Verb:      "create",
		ArgSpecs: core.ArgSpecs{
			{
				Name:  nestedMapArgName,
				Short: "Kubelet args",
			},
			{
				Name:  nestedArrayArgName,
				Short: "Tags",
			},
		},
	}

	schema := server.CommandToFlatArgsSchema(cmd)

	for _, argName := range []string{nestedMapArgName, nestedArrayArgName} {
		prop, ok := schema.Properties[argName]
		if !ok {
			t.Fatalf("Expected nested arg %q to be exposed as a property", argName)
		}
		if prop.Type != "string" {
			t.Fatalf("Expected nested arg %q to be a string, got %v", argName, prop.Type)
		}
		if prop.AdditionalProperties != nil {
			t.Fatalf("Expected nested arg %q to NOT have additionalProperties, got %v", argName, prop.AdditionalProperties)
		}
		if prop.Items != nil {
			t.Fatalf("Expected nested arg %q to NOT have items, got %v", argName, prop.Items)
		}
	}
}
