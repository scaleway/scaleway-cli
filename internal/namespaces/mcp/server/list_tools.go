package server

import (
	"context"
	"reflect"
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
)

// toolInfo represents information about an MCP tool
type toolInfo struct {
	Namespace string `json:"namespace"`
	Resource  string `json:"resource"`
	Verb      string `json:"verb"`
	ToolName  string `json:"tool_name"`
	Short     string `json:"short"`
}

// ListTools returns a list of available MCP tools based on the server's configuration.
// It returns tools that match the read-only mode and enabled namespaces/resources/verbs.
func (s *MCPServer) ListTools() []toolInfo {
	tools := make([]toolInfo, 0, len(s.commands))

	for _, cmd := range s.commands {
		tools = append(tools, toolInfo{
			Namespace: cmd.Command.Namespace,
			Resource:  cmd.Command.Resource,
			Verb:      cmd.Command.Verb,
			ToolName:  CommandNameToToolName(cmd.Command),
			Short:     cmd.Command.Short,
		})
	}

	// Sort tools by namespaces, resources, verbs for consistent output
	sort.Slice(tools, func(i, j int) bool {
		if tools[i].Namespace != tools[j].Namespace {
			return tools[i].Namespace < tools[j].Namespace
		}
		if tools[i].Resource != tools[j].Resource {
			return tools[i].Resource < tools[j].Resource
		}

		return tools[i].Verb < tools[j].Verb
	})

	return tools
}

func McpServerListTools() *core.Command {
	type listArgs struct {
		Namespaces string `json:"namespaces"`
		Resources  string `json:"resources"`
		Verbs      string `json:"verbs"`
		ReadOnly   bool   `json:"read-only"`
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "list-tools",
		Short:                "List available MCP tools",
		Long:                 "Lists all CLI commands that would be exposed as MCP tools by the server. Use filters to see which commands are available for specific namespaces, resources, or verbs.",
		AllowAnonymousClient: true,
		DisableTelemetry:     true,
		ArgsType:             reflect.TypeOf(listArgs{}),
		Examples: []*core.Example{
			{
				Short: "List all available MCP tools",
				Raw:   `scw mcp server list-tools`,
			},
			{
				Short: "List tools for a specific namespaces",
				Raw:   `scw mcp server list-tools namespaces=instance`,
			},
			{
				Short: "List tools for a specific resources",
				Raw:   `scw mcp server list-tools resources=server`,
			},
			{
				Short: "List only read-only tools (get/list operations)",
				Raw:   `scw mcp server list-tools read-only=true`,
			},
			{
				Short: "List tools with a specific verbs",
				Raw:   `scw mcp server list-tools verbs=get`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespaces",
				Short:      "Filter by namespaces (e.g., instance, iam, object)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "resources",
				Short:      "Filter by resources (e.g., server, volume, bucket)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "verbs",
				Short:      "Filter by verbs (e.g., get, list, create)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "read-only",
				Short:      "Only list read-only tools (get, list operations)",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*listArgs)

			// Get all CLI commands from the meta context
			commands := core.ExtractCommands(ctx)
			cliCommands := commands.GetAll()

			// Step 1: Filter commands based on the given config
			filteredCommands := FilterCommands(cliCommands, CommandFilterConfig{
				ReadOnly:          args.ReadOnly,
				EnabledNamespaces: SplitArg(args.Namespaces),
				EnabledResources:  SplitArg(args.Resources),
				EnabledVerbs:      SplitArg(args.Verbs),
			})

			// Step 2: Create the MCP server with pre-filtered commands
			mcpServer := NewMCPServer(filteredCommands)

			// Step 2: List tools from the MCP server
			return mcpServer.ListTools(), nil
		},
		ExcludeFromMCP: true,
	}
}
