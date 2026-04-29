package server

import (
	"context"
	"reflect"
	"sort"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func McpServerServe() *core.Command {
	type serveArgs struct {
		Transport        string   `json:"transport"`
		Address          string   `json:"address"`
		ReadOnly         bool     `json:"read-only"`
		EnableNamespaces []string `json:"enable-namespaces"`
		EnableResources  []string `json:"enable-resources"`
		EnableVerbs      []string `json:"enable-verbs"`
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "serve",
		Short:                "Start the MCP server",
		Long:                 "Runs the MCP server, exposing all CLI commands as MCP tools for AI assistants. Supports stdio (default), SSE, and streamable HTTP transports.",
		AllowAnonymousClient: true,
		DisableTelemetry:     true,
		ArgsType:             reflect.TypeOf(serveArgs{}),
		Examples: []*core.Example{
			{
				Short: "Start the MCP server with stdio transport (default)",
				Raw:   `scw mcp server serve`,
			},
			{
				Short: "Start the MCP server in read-only mode (only get/list operations)",
				Raw:   `scw mcp server serve --read-only`,
			},
			{
				Short: "Only serve commands from specific namespaces",
				Raw:   `scw mcp server serve --enable-namespaces instance,iam,object`,
			},
			{
				Short: "Only serve commands from specific resources",
				Raw:   `scw mcp server serve --enable-resources server,volume,bucket`,
			},
			{
				Short: "Only serve commands with specific verbs",
				Raw:   `scw mcp server serve --enable-verbs get,list,create`,
			},
			{
				Short: "Combine filters to serve only instance server get/list commands",
				Raw:   `scw mcp server serve --enable-namespaces instance --enable-resources server --enable-verbs get,list`,
			},
			{
				Short: "Start the MCP server with SSE transport on port 8080",
				Raw:   `scw mcp server serve --transport sse --address :8080`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "transport",
				Short:      "Transport mode: stdio (default), sse, or streamable-http",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("stdio"),
			},
			{
				Name:       "address",
				Short:      "Address to bind for SSE and streamable-http transports (e.g., :8080)",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter(":8080"),
			},
			{
				Name:       "read-only",
				Short:      "Only register read-only commands (get, list operations)",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
			{
				Name:       "enable-namespaces",
				Short:      "Only serve commands from specified namespaces (comma-separated)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "enable-resources",
				Short:      "Only serve commands from specified resources (comma-separated)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "enable-verbs",
				Short:      "Only serve commands with specified verbs (comma-separated)",
				Required:   false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*serveArgs)

			return RunMCPServer(
				ctx,
				args.Transport,
				args.Address,
				args.ReadOnly,
				args.EnableNamespaces,
				args.EnableResources,
				args.EnableVerbs,
			)
		},
		ExcludeFromMCP: true, // Skip mcp namespace to avoid recursive server calls
	}
}

func McpServerListTools() *core.Command {
	type listArgs struct {
		Namespace string `json:"namespace"`
		Resource  string `json:"resource"`
		Verb      string `json:"verb"`
		ReadOnly  bool   `json:"read-only"`
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
				Short: "List tools for a specific namespace",
				Raw:   `scw mcp server list-tools namespace=instance`,
			},
			{
				Short: "List tools for a specific resource",
				Raw:   `scw mcp server list-tools resource=server`,
			},
			{
				Short: "List only read-only tools (get/list operations)",
				Raw:   `scw mcp server list-tools read-only=true`,
			},
			{
				Short: "List tools with a specific verb",
				Raw:   `scw mcp server list-tools verb=get`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace",
				Short:      "Filter by namespace (e.g., instance, iam, object)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "resource",
				Short:      "Filter by resource (e.g., server, volume, bucket)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "verb",
				Short:      "Filter by verb (e.g., get, list, create)",
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

			// Build filter arrays from single string args
			var enabledNamespaces, enabledResources, enabledVerbs []string
			if args.Namespace != "" {
				enabledNamespaces = []string{args.Namespace}
			}
			if args.Resource != "" {
				enabledResources = []string{args.Resource}
			}
			if args.Verb != "" {
				enabledVerbs = []string{args.Verb}
			}

			// Collect matching commands
			type toolInfo struct {
				Namespace string `json:"namespace"`
				Resource  string `json:"resource"`
				Verb      string `json:"verb"`
				ToolName  string `json:"tool_name"`
				Short     string `json:"short"`
			}

			var tools []toolInfo
			for _, cmd := range cliCommands {
				if ShouldRegisterCommand(
					cmd,
					args.ReadOnly,
					enabledNamespaces,
					enabledResources,
					enabledVerbs,
				) {
					tools = append(tools, toolInfo{
						Namespace: cmd.Namespace,
						Resource:  cmd.Resource,
						Verb:      cmd.Verb,
						ToolName:  CommandNameToToolName(cmd),
						Short:     cmd.Short,
					})
				}
			}

			// Sort tools by namespace, resource, verb for consistent output
			sort.Slice(tools, func(i, j int) bool {
				if tools[i].Namespace != tools[j].Namespace {
					return tools[i].Namespace < tools[j].Namespace
				}
				if tools[i].Resource != tools[j].Resource {
					return tools[i].Resource < tools[j].Resource
				}

				return tools[i].Verb < tools[j].Verb
			})

			return tools, nil
		},
		ExcludeFromMCP: true,
	}
}

func McpServer() *core.Command {
	return &core.Command{
		Groups:         []string{"utility"},
		Short:          `MCP server management commands`,
		Long:           `Commands for managing the MCP server that exposes Scaleway CLI commands as AI tools.`,
		Namespace:      "mcp",
		Resource:       "server",
		ExcludeFromMCP: true,
	}
}

func McpServerListResources() *core.Command {
	type listResourcesArgs struct {
		Namespace string `json:"namespace"`
		Resource  string `json:"resource"`
		ReadOnly  bool   `json:"read-only"`
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "list-resources",
		Short:                "List available MCP resources",
		Long:                 "Lists all CLI commands that would be exposed as MCP resources by the server. Resources are read-only endpoints for list commands that can be accessed via URI. Use filters to see which resources are available for specific namespaces or resources.",
		AllowAnonymousClient: true,
		DisableTelemetry:     true,
		ArgsType:             reflect.TypeOf(listResourcesArgs{}),
		Examples: []*core.Example{
			{
				Short: "List all available MCP resources",
				Raw:   `scw mcp server list-resources`,
			},
			{
				Short: "List resources for a specific namespace",
				Raw:   `scw mcp server list-resources namespace=instance`,
			},
			{
				Short: "List resources for a specific resource type",
				Raw:   `scw mcp server list-resources resource=server`,
			},
			{
				Short: "List only read-only resources",
				Raw:   `scw mcp server list-resources read-only=true`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "namespace",
				Short:      "Filter by namespace (e.g., instance, iam, object)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "resource",
				Short:      "Filter by resource (e.g., server, volume, bucket)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "read-only",
				Short:      "Only list read-only resources",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("false"),
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*listResourcesArgs)

			// Get all CLI commands from the meta context
			commands := core.ExtractCommands(ctx)
			cliCommands := commands.GetAll()

			// Build filter arrays from single string args
			var enabledNamespaces, enabledResources []string
			if args.Namespace != "" {
				enabledNamespaces = []string{args.Namespace}
			}
			if args.Resource != "" {
				enabledResources = []string{args.Resource}
			}

			// Collect matching resources
			type resourceInfo struct {
				Namespace string `json:"namespace"`
				Resource  string `json:"resource"`
				URI       string `json:"uri"`
				Short     string `json:"short"`
			}

			var resources []resourceInfo
			for _, cmd := range cliCommands {
				// Only list commands are exposed as resources
				if IsListCommand(cmd) &&
					ShouldRegisterCommand(
						cmd,
						args.ReadOnly,
						enabledNamespaces,
						enabledResources,
						nil,
					) {
					resources = append(resources, resourceInfo{
						Namespace: cmd.Namespace,
						Resource:  cmd.Resource,
						URI:       BuildResourceURI(cmd.Namespace, cmd.Resource),
						Short:     cmd.Short,
					})
				}
			}

			// Sort resources by namespace, resource for consistent output
			sort.Slice(resources, func(i, j int) bool {
				if resources[i].Namespace != resources[j].Namespace {
					return resources[i].Namespace < resources[j].Namespace
				}

				return resources[i].Resource < resources[j].Resource
			})

			return resources, nil
		},
		ExcludeFromMCP: true,
	}
}
