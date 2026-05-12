package server

import (
	"context"
	"fmt"
	"os"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func McpServerServe() *core.Command {
	type serveArgs struct {
		Transport string `json:"transport"`
		Address   string `json:"address"`
		ReadOnly  bool   `json:"read-only"`
		Namespace string `json:"namespace"`
		Resource  string `json:"resource"`
		Verb      string `json:"verb"`
	}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "serve",
		Short:                "Start the MCP server",
		Long:                 "Runs the MCP server, exposing all CLI commands as MCP tools for AI assistants. Supports stdio (default) and streamable HTTP transports.",
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
				Raw:   `scw mcp server serve namespace=instance,iam,object`,
			},
			{
				Short: "Only serve commands from specific resources",
				Raw:   `scw mcp server serve resource=server,volume,bucket`,
			},
			{
				Short: "Only serve commands with specific verbs",
				Raw:   `scw mcp server serve verb=get,list,create`,
			},
			{
				Short: "Combine filters to serve only instance server get/list commands",
				Raw:   `scw mcp server serve namespace=instance resource=server verb=get,list`,
			},
		},
		ArgSpecs: core.ArgSpecs{
			{
				Name:       "transport",
				Short:      "Transport mode: stdio (default) or streamable-http",
				Required:   false,
				Positional: false,
				Default:    core.DefaultValueSetter("stdio"),
			},
			{
				Name:       "address",
				Short:      "Address to bind for streamable-http transports (e.g., :8080)",
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
				Name:       "namespace",
				Short:      "Only serve commands from specified namespaces (comma-separated)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "resource",
				Short:      "Only serve commands from specified resources (comma-separated)",
				Required:   false,
				Positional: false,
			},
			{
				Name:       "verb",
				Short:      "Only serve commands with specified verbs (comma-separated)",
				Required:   false,
				Positional: false,
			},
		},
		Run: func(ctx context.Context, argsI any) (any, error) {
			args := argsI.(*serveArgs)

			meta := core.ExtractMeta(ctx)

			// Reload the client with the profile to ensure proper authentication
			if err := core.ReloadClient(ctx); err != nil {
				return nil, fmt.Errorf("failed to initialize authenticated client: %w", err)
			}

			// Log startup information to stderr
			fmt.Fprintf(
				os.Stderr,
				"Starting MCP server version %s\n",
				meta.BuildInfo.Version.String(),
			)
			fmt.Fprintf(os.Stderr, "Using profile: %s\n", meta.ProfileFlag)
			fmt.Fprintf(os.Stderr, "Config path: %s\n", meta.ConfigPathFlag)
			fmt.Fprintf(os.Stderr, "Transport mode: %s\n", args.Transport)
			fmt.Fprintf(os.Stderr, "Read-only mode: %v\n", args.ReadOnly)
			if len(SplitArg(args.Namespace)) > 0 {
				fmt.Fprintf(os.Stderr, "Enabled namespaces: %v\n", SplitArg(args.Namespace))
			}
			if len(SplitArg(args.Resource)) > 0 {
				fmt.Fprintf(os.Stderr, "Enabled resources: %v\n", SplitArg(args.Resource))
			}
			if len(SplitArg(args.Verb)) > 0 {
				fmt.Fprintf(os.Stderr, "Enabled verbs: %v\n", SplitArg(args.Verb))
			}

			// Get all CLI commands
			commands := core.ExtractCommands(ctx)
			cliCommands := commands.GetAll()

			// Copy OverrideEnv from context
			for _, envKey := range []string{"HOME", "PATH", scw.ScwAccessKeyEnv, scw.ScwSecretKeyEnv, scw.ScwDefaultOrganizationIDEnv, scw.ScwDefaultProjectIDEnv, scw.ScwDefaultRegionEnv, scw.ScwDefaultZoneEnv} {
				if val := core.ExtractEnv(ctx, envKey); val != "" {
					meta.OverrideEnv[envKey] = val
				}
			}

			// Step 1: Filter commands based on the given config
			filteredCommands := FilterCommands(cliCommands, CommandFilterConfig{
				ReadOnly:          args.ReadOnly,
				EnabledNamespaces: SplitArg(args.Namespace),
				EnabledResources:  SplitArg(args.Resource),
				EnabledVerbs:      SplitArg(args.Verb),
			})

			// Step 2: Inject meta into context for tool/resource execution
			ctx = core.InjectMeta(ctx, meta)

			// Step 3: Create the MCP server with pre-filtered commands
			mcpServer := NewMCPServer(filteredCommands)

			// Step 2: Serve the MCP server with the specified transport
			return mcpServer.Serve(ctx, args.Transport, args.Address)
		},
		ExcludeFromMCP: true, // Skip mcp namespace to avoid recursive server calls
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

			// Step 1: Filter commands based on the given config
			filteredCommands := FilterCommands(cliCommands, CommandFilterConfig{
				ReadOnly:          args.ReadOnly,
				EnabledNamespaces: SplitArg(args.Namespace),
				EnabledResources:  SplitArg(args.Resource),
			})

			// Step 2: Create the MCP server with pre-filtered commands
			mcpServer := NewMCPServer(filteredCommands)

			// Step 2: List resources from the MCP server
			return mcpServer.ListResources(), nil
		},
		ExcludeFromMCP: true,
	}
}
