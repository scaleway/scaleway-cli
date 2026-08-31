package server

import (
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

const (
	dynamicArrayPlaceholder = ".{index}"
	dynamicMapPlaceholder   = ".{key}"
	dynamicArgSeparator     = "."
	disallowAdditionalProps = false
	jsonSchemaTypeArray     = "array"
	jsonSchemaTypeObject    = "object"
	jsonSchemaTypeString    = "string"
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
		Description: argSpec.Short,
	}

	// Handle enum values
	if len(argSpec.EnumValues) > 0 {
		schema.Enum = argSpec.EnumValues
		schema.Type = jsonSchemaTypeString

		return schema
	}

	// Default to string for most args
	schema.Type = jsonSchemaTypeString

	return schema
}

// CommandToFlatArgsSchema creates a flat schema for commands that accept all args as strings
func CommandToFlatArgsSchema(cmd *core.Command) *JSONSchema {
	schema := &JSONSchema{
		Type:                 jsonSchemaTypeObject,
		Properties:           make(map[string]*JSONSchema),
		Required:             []string{},
		AdditionalProperties: disallowAdditionalProps,
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

func argSpecToPropertySchema(argSpec *core.ArgSpec) (string, *JSONSchema) {
	if strings.Count(argSpec.Name, dynamicArrayPlaceholder) == 1 &&
		strings.HasSuffix(argSpec.Name, dynamicArrayPlaceholder) {
		propName := strings.TrimSuffix(argSpec.Name, dynamicArrayPlaceholder)

		return strcase.ToKebab(propName), &JSONSchema{
			Type:        jsonSchemaTypeArray,
			Description: argSpec.Short,
			Items: &JSONSchema{
				Type: jsonSchemaTypeString,
			},
		}
	}

	if strings.Count(argSpec.Name, dynamicMapPlaceholder) == 1 &&
		strings.HasSuffix(argSpec.Name, dynamicMapPlaceholder) {
		propName := strings.TrimSuffix(argSpec.Name, dynamicMapPlaceholder)

		return strcase.ToKebab(propName), &JSONSchema{
			Type:        jsonSchemaTypeObject,
			Description: argSpec.Short,
			AdditionalProperties: &JSONSchema{
				Type: jsonSchemaTypeString,
			},
		}
	}

	propName := strcase.ToKebab(argSpec.Name)
	propSchema := &JSONSchema{
		Type:        jsonSchemaTypeString,
		Description: argSpec.Short,
	}

	if len(argSpec.EnumValues) > 0 {
		propSchema.Enum = argSpec.EnumValues
	}

	return propName, propSchema
}
