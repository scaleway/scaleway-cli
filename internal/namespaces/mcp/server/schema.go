package server

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// JSONSchema represents a JSON Schema object
type JSONSchema struct {
	Type                 string                 `json:"type,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties *bool                  `json:"additionalProperties,omitempty"`
	Enum                 []string               `json:"enum,omitempty"`
	Default              any                    `json:"default,omitempty"`
}

// ArgSpecToJSONSchema converts a core.ArgSpec to JSON Schema
func ArgSpecToJSONSchema(argSpec *core.ArgSpec) *JSONSchema {
	schema := &JSONSchema{
		Description: argSpec.Short,
	}

	// Handle enum values
	if len(argSpec.EnumValues) > 0 {
		schema.Enum = argSpec.EnumValues
		schema.Type = "string"

		return schema
	}

	// Default to string for most args
	schema.Type = "string"

	return schema
}

// CommandToFlatArgsSchema creates a flat schema for commands that accept all args as strings
func CommandToFlatArgsSchema(cmd *core.Command) *JSONSchema {
	schema := &JSONSchema{
		Type:                 "object",
		Properties:           make(map[string]*JSONSchema),
		Required:             []string{},
		AdditionalProperties: new(false),
	}

	for _, argSpec := range cmd.ArgSpecs {
		propName := strcase.ToKebab(argSpec.Name)
		propSchema := &JSONSchema{
			Type:        "string",
			Description: argSpec.Short,
		}

		if len(argSpec.EnumValues) > 0 {
			propSchema.Enum = argSpec.EnumValues
		}

		schema.Properties[propName] = propSchema

		if argSpec.Required {
			schema.Required = append(schema.Required, propName)
		}
	}

	return schema
}
