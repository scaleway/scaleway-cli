package server_test

import (
	"encoding/json"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

const (
	arrayArgName          = "tags.{index}"
	arrayPropertyName     = "tags"
	mapArgName            = "environment-variables.{key}"
	mapPropertyName       = "environment-variables"
	schemaTypeArray       = "array"
	schemaTypeObject      = "object"
	schemaTypeString      = "string"
	schemaPropertiesKey   = "properties"
	schemaTypeKey         = "type"
	schemaItemsKey        = "items"
	schemaAdditionalProps = "additionalProperties"
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

	schema := server.CommandToFlatArgsSchema(cmd)
	rawSchema, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rawSchema, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	properties := decoded[schemaPropertiesKey].(map[string]any)
	if _, ok := properties[arrayArgName]; ok {
		t.Fatalf("Schema should not expose literal placeholder property %q", arrayArgName)
	}
	if _, ok := properties[mapArgName]; ok {
		t.Fatalf("Schema should not expose literal placeholder property %q", mapArgName)
	}

	tags := properties[arrayPropertyName].(map[string]any)
	if tags[schemaTypeKey] != schemaTypeArray {
		t.Fatalf(
			"Expected %q to be an array schema, got %v",
			arrayPropertyName,
			tags[schemaTypeKey],
		)
	}
	tagItems := tags[schemaItemsKey].(map[string]any)
	if tagItems[schemaTypeKey] != schemaTypeString {
		t.Fatalf(
			"Expected %q items to be strings, got %v",
			arrayPropertyName,
			tagItems[schemaTypeKey],
		)
	}

	environmentVariables := properties[mapPropertyName].(map[string]any)
	if environmentVariables[schemaTypeKey] != schemaTypeObject {
		t.Fatalf(
			"Expected %q to be an object schema, got %v",
			mapPropertyName,
			environmentVariables[schemaTypeKey],
		)
	}
	additionalProperties := environmentVariables[schemaAdditionalProps].(map[string]any)
	if additionalProperties[schemaTypeKey] != schemaTypeString {
		t.Fatalf(
			"Expected %q values to be strings, got %v",
			mapPropertyName,
			additionalProperties[schemaTypeKey],
		)
	}
}

func TestCommandToFlatArgsSchemaNestedDynamicArgs(t *testing.T) {
	// Arg specs with nested placeholders (e.g., "pools.{index}.kubelet-args.{key}")
	// should NOT be treated as simple maps or arrays. They must fall through to the
	// default string case so that the literal placeholder name is exposed as the
	// property name, matching pre-existing behavior.
	// See https://github.com/scaleway/scaleway-cli for the original bug report.
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
	rawSchema, err := json.Marshal(schema)
	if err != nil {
		t.Fatalf("Failed to marshal schema: %v", err)
	}

	var decoded map[string]any
	if err := json.Unmarshal(rawSchema, &decoded); err != nil {
		t.Fatalf("Failed to unmarshal schema: %v", err)
	}

	properties := decoded[schemaPropertiesKey].(map[string]any)

	// The literal placeholder name should be the property key (kebab-cased),
	// and it should be a plain string, NOT an object/array schema.
	nestedMapPropName := "pools.{index}.kubelet-args.{key}"
	nestedArrayPropName := "pools.{index}.tags.{index}"

	mapProp, ok := properties[nestedMapPropName].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected nested map arg %q to be exposed as a string property, got %v",
			nestedMapArgName,
			properties[nestedMapPropName],
		)
	}
	if mapProp[schemaTypeKey] != schemaTypeString {
		t.Fatalf(
			"Expected nested map arg %q to be a string, got %v",
			nestedMapArgName,
			mapProp[schemaTypeKey],
		)
	}
	if _, hasAdditionalProps := mapProp[schemaAdditionalProps]; hasAdditionalProps {
		t.Fatalf(
			"Expected nested map arg %q to NOT have additionalProperties, got %v",
			nestedMapArgName,
			mapProp,
		)
	}

	arrayProp, ok := properties[nestedArrayPropName].(map[string]any)
	if !ok {
		t.Fatalf(
			"Expected nested array arg %q to be exposed as a string property, got %v",
			nestedArrayArgName,
			properties[nestedArrayPropName],
		)
	}
	if arrayProp[schemaTypeKey] != schemaTypeString {
		t.Fatalf(
			"Expected nested array arg %q to be a string, got %v",
			nestedArrayArgName,
			arrayProp[schemaTypeKey],
		)
	}
	if _, hasItems := arrayProp[schemaItemsKey]; hasItems {
		t.Fatalf(
			"Expected nested array arg %q to NOT have items, got %v",
			nestedArrayArgName,
			arrayProp,
		)
	}
}
