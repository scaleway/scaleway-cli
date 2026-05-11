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
	// meta from bootstrap context for HTTP transport
	meta *core.Meta
}

// NewCommandTool creates a new CommandTool from a core.Command with optional baseMeta
func NewCommandTool(cmd *core.Command, baseMeta *core.Meta) *CommandTool {
	return &CommandTool{
		Command: cmd,
		meta:    baseMeta,
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
	ctx, err := injectMetaIfMissing(ctx, ct.meta)
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
			}, nil
		}
	} else {
		// Fallback for commands without ArgsType
		cmdArgs = inputArgs
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

// injectMetaIfMissing initializes context with meta if not already present.
// This is required when running as an MCP server (especially HTTP streamable)
// where the context doesn't go through the normal CLI bootstrap process.
// Returns the updated context and any error that occurred during initialization.
func injectMetaIfMissing(ctx context.Context, meta *core.Meta) (context.Context, error) {
	if core.ExtractClient(ctx) != nil {
		return ctx, nil
	}
	if meta != nil {
		m := *meta
		m.OverrideEnv = make(map[string]string)
		maps.Copy(meta.OverrideEnv, meta.OverrideEnv)

		return core.InjectMeta(ctx, &m), nil
	}
	// baseMeta is nil - create an authenticated client from environment/config
	// Load config to get active profile and credentials
	configPath := scw.GetConfigPath()
	config, err := scw.LoadConfigFromPath(configPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		return ctx, fmt.Errorf("loading config: %w", err)
	}

	// Get active profile name (from config or env)
	profileName := scw.DefaultProfileName
	if config != nil && config.ActiveProfile != nil {
		profileName = *config.ActiveProfile
	}
	if envProfile := os.Getenv(scw.ScwActiveProfileEnv); envProfile != "" {
		profileName = envProfile
	}

	// Build client options with defaults and profile
	options := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithUserAgent("scaleway-mcp-server"),
		scw.WithEnv(), // Load credentials from environment variables
	}

	// Load profile from config file if available
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
