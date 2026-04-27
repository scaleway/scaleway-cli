package server

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// CommandTool wraps a CLI command for MCP tool execution
type CommandTool struct {
	Command *core.Command
}

// NewCommandTool creates a new CommandTool from a core.Command
func NewCommandTool(cmd *core.Command) *CommandTool {
	return &CommandTool{
		Command: cmd,
	}
}

// ToMCPTool converts the CommandTool to an MCP Tool
func (ct *CommandTool) ToMCPTool() *mcp.Tool {
	tool := &mcp.Tool{
		Name:        CommandNameToToolName(ct.Command),
		Description: ct.Command.Short,
		InputSchema: CommandToFlatArgsSchema(ct.Command),
	}

	// Add command metadata in description for debugging
	if ct.Command.Long != "" {
		tool.Description = ct.Command.Short + "\n\n" + truncateString(ct.Command.Long, 200)
	}

	return tool
}

// Execute runs the CLI command with the provided arguments
func (ct *CommandTool) Execute(ctx context.Context, args map[string]any) (*mcp.CallToolResult, error) {
	// Build command arguments for the CLI runner
	cmdArgs := make(map[string]interface{})

	for argName, argValue := range args {
		// Convert from kebab-case back to original arg spec name
		originalName := argName
		for _, spec := range ct.Command.ArgSpecs {
			if strcase.ToKebab(spec.Name) == argName {
				originalName = spec.Name
				break
			}
		}

		cmdArgs[originalName] = argValue
	}

	// Execute the command's Run function
	result, err := ct.Command.Run(ctx, cmdArgs)

	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: %v", err),
				},
			},
			IsError: true,
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

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: resultStr,
			},
		},
	}, nil
}

// CommandNameToToolName converts a command to an MCP tool name
func CommandNameToToolName(cmd *core.Command) string {
	parts := []string{}

	if cmd.Namespace != "" {
		parts = append(parts, cmd.Namespace)
	}
	if cmd.Resource != "" {
		parts = append(parts, cmd.Resource)
	}
	if cmd.Verb != "" {
		parts = append(parts, cmd.Verb)
	}

	// Join with underscores for MCP tool naming convention
	toolName := strings.Join(parts, "_")

	// Ensure tool name doesn't exceed MCP limits (typically 64 chars)
	if len(toolName) > 64 {
		toolName = toolName[:64]
	}

	return toolName
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}

	return s[:maxLen-3] + "..."
}
