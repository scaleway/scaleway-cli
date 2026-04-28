package server

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// CommandResource wraps a CLI command for MCP resource exposure
type CommandResource struct {
	Command *core.Command
}

// NewCommandResource creates a new CommandResource from a core.Command
func NewCommandResource(cmd *core.Command) *CommandResource {
	return &CommandResource{
		Command: cmd,
	}
}

// ToMCPResource converts the CommandResource to an MCP Resource
func (cr *CommandResource) ToMCPResource() *mcp.Resource {
	// Build URI: scw://{namespace}/{resource}
	uri := BuildResourceURI(cr.Command.Namespace, cr.Command.Resource)

	resource := &mcp.Resource{
		URI:         uri,
		Name:        ResourceName(cr.Command),
		Description: cr.Command.Short,
		MIMEType:    "application/json",
	}

	// Add command metadata in description for debugging
	if cr.Command.Long != "" {
		resource.Description = cr.Command.Short + "\n\n" + truncateString(cr.Command.Long, 200)
	}

	return resource
}

// ResourceName generates a human-readable name for the resource
func ResourceName(cmd *core.Command) string {
	parts := []string{}

	if cmd.Namespace != "" {
		parts = append(parts, cmd.Namespace)
	}
	if cmd.Resource != "" {
		parts = append(parts, cmd.Resource)
	}

	return strings.Join(parts, " ")
}

// BuildResourceURI creates a URI for a namespace/resource combination
func BuildResourceURI(namespace, resource string) string {
	return fmt.Sprintf("scw://%s/%s", namespace, resource)
}

// Execute runs the CLI command with the provided arguments and returns resource contents
func (cr *CommandResource) Execute(
	ctx context.Context,
	inputArgs map[string]any,
) (*mcp.ReadResourceResult, error) {
	// Skip commands without a Run function
	if cr.Command.Run == nil {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  BuildResourceURI(cr.Command.Namespace, cr.Command.Resource),
					Text: "Command has no execution logic",
				},
			},
		}, nil
	}

	// Create a new instance of the expected args type
	var cmdArgs any
	if cr.Command.ArgsType != nil {
		cmdArgs = reflect.New(cr.Command.ArgsType).Interface()

		// Convert map[string]any to []string format expected by args.UnmarshalStruct
		// Format: ["arg1=value1", "arg2=value2"]
		rawArgs := make([]string, 0, len(inputArgs))
		for argName, argValue := range inputArgs {
			// Convert from kebab-case back to original arg spec name
			originalName := argName
			for _, spec := range cr.Command.ArgSpecs {
				if strcase.ToKebab(spec.Name) == argName {
					originalName = spec.Name

					break
				}
			}

			// Convert value to string
			var valueStr string
			switch v := argValue.(type) {
			case string:
				valueStr = v
			case bool:
				valueStr = strconv.FormatBool(v)
			case float64:
				valueStr = fmt.Sprintf("%v", v)
			case int:
				valueStr = strconv.Itoa(v)
			default:
				// For complex types, marshal to JSON
				if b, err := json.Marshal(v); err == nil {
					valueStr = string(b)
				} else {
					valueStr = fmt.Sprintf("%v", v)
				}
			}

			rawArgs = append(rawArgs, originalName+"="+valueStr)
		}

		// Unmarshal the args into the struct
		if err := args.UnmarshalStruct(rawArgs, cmdArgs); err != nil {
			return &mcp.ReadResourceResult{
				Contents: []*mcp.ResourceContents{
					{
						URI:  BuildResourceURI(cr.Command.Namespace, cr.Command.Resource),
						Text: fmt.Sprintf("Error parsing arguments: %v", err),
					},
				},
			}, nil
		}
	} else {
		// Fallback for commands without ArgsType
		cmdArgs = inputArgs
	}

	// Execute the command's Run function
	result, err := cr.Command.Run(ctx, cmdArgs)
	if err != nil {
		return &mcp.ReadResourceResult{
			Contents: []*mcp.ResourceContents{
				{
					URI:  BuildResourceURI(cr.Command.Namespace, cr.Command.Resource),
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
		}, nil
	}

	// Format the result
	var resultStr string
	if result != nil {
		if data, err := json.MarshalIndent(result, "", "  "); err == nil {
			resultStr = string(data)
		} else {
			resultStr = fmt.Sprintf("%v", result)
		}
	} else {
		resultStr = "Command executed successfully (no output)"
	}

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      BuildResourceURI(cr.Command.Namespace, cr.Command.Resource),
				MIMEType: "application/json",
				Text:     resultStr,
			},
		},
	}, nil
}
