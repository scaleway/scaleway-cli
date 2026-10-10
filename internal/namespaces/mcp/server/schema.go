package server

import (
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// JSONSchema represents a JSON Schema object
type JSONSchema struct {
	Type                 string                 `json:"type,omitempty"`
	Description          string                 `json:"description,omitempty"`
	Properties           map[string]*JSONSchema `json:"properties,omitempty"`
	Required             []string               `json:"required,omitempty"`
	AdditionalProperties any                    `json:"additionalProperties,omitempty"`
	Enum                 []string               `json:"enum,omitempty"`
	Default              any                    `json:"default,omitempty"`
	Items                *JSONSchema            `json:"items,omitempty"`
}

// ArgSpecToJSONSchema converts a core.ArgSpec to JSON Schema
func ArgSpecToJSONSchema(argSpec *core.ArgSpec) *JSONSchema {
	schema := &JSONSchema{
		Type:        "string",
		Description: argSpec.Short,
	}

	if len(argSpec.EnumValues) > 0 {
		schema.Enum = argSpec.EnumValues
	}

	return schema
}

// CommandToFlatArgsSchema creates a flat schema for commands that accept all args as strings
func CommandToFlatArgsSchema(cmd *core.Command) *JSONSchema {
	schema := &JSONSchema{
		Type:                 "object",
		Properties:           make(map[string]*JSONSchema),
		Required:             []string{},
		AdditionalProperties: false,
	}

	for _, argSpec := range cmd.ArgSpecs {
		propName, propSchema := argSpecToPropertySchema(argSpec)

		schema.Properties[propName] = propSchema

		if argSpec.Required {
			schema.Required = append(schema.Required, propName)
		}
	}

	return schema
}

// dynamicPrefix reports whether argName ends with a single dynamic placeholder
// ("tags.{index}", "environment-variables.{key}") and returns its prefix.
// Nested placeholders ("pools.{index}.tags.{index}") do not match.
func dynamicPrefix(argName string) (prefix string, isMap bool) {
	if p, ok := strings.CutSuffix(argName, ".{index}"); ok && !strings.Contains(p, ".{") {
		return p, false
	}

	if p, ok := strings.CutSuffix(argName, ".{key}"); ok && !strings.Contains(p, ".{") {
		return p, true
	}

	return "", false
}

func argSpecToPropertySchema(argSpec *core.ArgSpec) (string, *JSONSchema) {
	if prefix, isMap := dynamicPrefix(argSpec.Name); prefix != "" {
		if isMap {
			return strcase.ToKebab(prefix), &JSONSchema{
				Type:                 "object",
				Description:          argSpec.Short,
				AdditionalProperties: &JSONSchema{Type: "string"},
			}
		}

		return strcase.ToKebab(prefix), &JSONSchema{
			Type:        "array",
			Description: argSpec.Short,
			Items:       &JSONSchema{Type: "string"},
		}
	}

	propName := strcase.ToKebab(argSpec.Name)
	propSchema := &JSONSchema{
		Type:        "string",
		Description: argSpec.Short,
	}

	if len(argSpec.EnumValues) > 0 {
		propSchema.Enum = argSpec.EnumValues
	}

	return propName, propSchema
}
