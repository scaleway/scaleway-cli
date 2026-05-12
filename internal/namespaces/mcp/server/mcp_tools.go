package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"os"
	"reflect"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-sdk-go/scw"
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
		Title:       CommandNameToToolName(ct.Command),
		Name:        CommandNameToToolName(ct.Command),
		Description: ct.Command.Short,
		InputSchema: CommandToFlatArgsSchema(ct.Command),
	}

	// Add command metadata in description for debugging
	if ct.Command.Long != "" {
		tool.Description = ct.Command.Short + "\n\n" + truncateString(ct.Command.Long, 200)
	}

	// Add command metadata (namespace, resource, verb) to the tool's Meta field
	// This allows clients to access structured information about the command
	// https://modelcontextprotocol.io/specification/2025-11-25/basic#_meta
	tool.Meta = mcp.Meta{
		"namespace": ct.Command.Namespace,
		"resource":  ct.Command.Resource,
		"verb":      ct.Command.Verb,
	}

	// https://modelcontextprotocol.io/specification/2025-11-25/server/resources#annotations

	tool.Annotations = &mcp.ToolAnnotations{
		OpenWorldHint:   new(true), // We interact with Scaleway for all registered commands
		DestructiveHint: new(ct.Command.IsDestructive()),
		IdempotentHint:  ct.Command.IsIdempotent(),
		ReadOnlyHint:    ct.Command.IsReadOnly(),
	}

	return tool
}

// Execute runs the CLI command with the provided arguments
func (ct *CommandTool) Execute(
	ctx context.Context,
	inputArgs map[string]any,
) (*mcp.CallToolResult, error) {
	// Ensure meta is available in context
	// If meta is already in context (from wrapper), this is a no-op
	// If not, this creates meta from environment/config (useful for tests)
	ctx, err := ensureMetaInContext(ctx, nil)
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error initializing client: %v", err),
				},
			},
			IsError: true,
		}, nil
	}

	// Verify client was successfully injected
	if core.ExtractClient(ctx) == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Error: client not initialized - check SCW credentials (access key, secret key, or profile)",
				},
			},
			IsError: true,
		}, nil
	}

	// Skip commands without a Run function
	if ct.Command.Run == nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: "Command has no execution logic",
				},
			},
		}, nil
	}

	// Create a new instance of the expected args type
	var cmdArgs any
	if ct.Command.ArgsType != nil {
		cmdArgs = reflect.New(ct.Command.ArgsType).Interface()

		// Convert map[string]any to []string format expected by args.UnmarshalStruct
		// Format: ["arg1=value1", "arg2=value2"]
		rawArgs := make([]string, 0, len(inputArgs))
		for argName, argValue := range inputArgs {
			// Convert from kebab-case back to original arg spec name
			originalName := argName
			for _, spec := range ct.Command.ArgSpecs {
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
				if b, marshalErr := json.Marshal(v); marshalErr == nil {
					valueStr = string(b)
				} else {
					valueStr = fmt.Sprintf("%v", v)
				}
			}

			rawArgs = append(rawArgs, originalName+"="+valueStr)
		}

		// Unmarshal the args into the struct
		if unmarshalErr := args.UnmarshalStruct(rawArgs, cmdArgs); unmarshalErr != nil {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{
						Text: fmt.Sprintf("Error parsing arguments: %v", unmarshalErr),
					},
				},
				IsError: true,
			}, unmarshalErr
		}
	} else {
		// Fallback for commands without ArgsType
		cmdArgs = inputArgs
	}

	// Execute the command's Run function using the interceptor chain
	// This ensures custom request wrappers (like customListServersRequest) are properly handled
	var result any
	var execErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				execErr = fmt.Errorf("panic recovered during command execution: %v", r)
			}
		}()

		// Create a runner that wraps the command's Run function
		var runner core.CommandRunner = func(ctx context.Context, argsI any) (i any, err error) {
			return ct.Command.Run(ctx, argsI)
		}

		// Apply command interceptor if present
		if ct.Command.Interceptor != nil {
			result, execErr = ct.Command.Interceptor(ctx, cmdArgs, runner)
		} else {
			result, execErr = runner(ctx, cmdArgs)
		}
	}()
	if execErr != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Error: %v", execErr),
				},
			},
			IsError: true,
		}, execErr
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

// ensureMetaInContext ensures meta is available in the context for command execution.
// If meta is already present, it returns the context unchanged.
// If meta is provided, it injects it into the context.
// If no meta is provided, it creates one from environment/config.
// Returns the updated context and any error that occurred.
func ensureMetaInContext(ctx context.Context, meta *core.Meta) (context.Context, error) {
	// Meta already in context - nothing to do
	if core.ExtractMeta(ctx) != nil {
		return ctx, nil
	}

	// Use provided meta
	if meta != nil {
		m := *meta
		m.OverrideEnv = make(map[string]string)
		maps.Copy(m.OverrideEnv, meta.OverrideEnv)

		return core.InjectMeta(ctx, &m), nil
	}

	// No meta provided - create one from environment/config
	configPath := scw.GetConfigPath()
	config, err := scw.LoadConfigFromPath(configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ctx, fmt.Errorf("loading config: %w", err)
	}

	profileName := scw.DefaultProfileName
	if config != nil && config.ActiveProfile != nil {
		profileName = *config.ActiveProfile
	}
	if envProfile := os.Getenv(scw.ScwActiveProfileEnv); envProfile != "" {
		profileName = envProfile
	}

	options := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent("scaleway-mcp-server"),
		scw.WithEnv(),
	}

	if config != nil {
		profile, err := config.GetProfile(profileName)
		if err != nil {
			return ctx, fmt.Errorf("getting profile %q: %w", profileName, err)
		}
		if profile != nil {
			options = append(options, scw.WithProfile(profile))
		}
	}

	client, err := scw.NewClient(options...)
	if err != nil {
		return ctx, fmt.Errorf("creating client: %w", err)
	}

	return core.InjectMeta(ctx, &core.Meta{
		Client:      client,
		OverrideEnv: map[string]string{},
		BinaryName:  "scw-mcp",
	}), nil
}
