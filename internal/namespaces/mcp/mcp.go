package mcp

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"syscall"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mcp/server"
)

func GetCommands() *core.Commands {
	return core.NewCommands(
		mcpRoot(),
		mcpServerCommand(),
	)
}

func mcpRoot() *core.Command {
	return &core.Command{
		Groups:    []string{"utility"},
		Namespace: "mcp",
		Short:     "MCP (Model Context Protocol) server management",
		Long:      "Commands for managing the MCP server that exposes Scaleway CLI commands as AI tools.",
	}
}

func mcpServerCommand() *core.Command {
	type serveArgs struct{}

	return &core.Command{
		Groups:               []string{"utility"},
		Namespace:            "mcp",
		Resource:             "server",
		Verb:                 "serve",
		Short:                "Start the MCP server",
		Long:                 "Runs the MCP server over stdio, exposing all CLI commands as MCP tools for AI assistants.",
		AllowAnonymousClient: true,
		DisableTelemetry:     true,
		ArgsType:             reflect.TypeOf(serveArgs{}),
		Run: func(ctx context.Context, argsI any) (any, error) {
			return runMCPServer(ctx)
		},
	}
}

func runMCPServer(ctx context.Context) (any, error) {
	// Get all CLI commands from the meta context
	commands := core.ExtractCommands(ctx)
	cliCommands := commands.GetAll()

	// Get build info for version
	buildInfo := core.ExtractBuildInfo(ctx)
	version := buildInfo.Version.String()

	// Create MCP server with all commands
	mcpServer := server.NewMCPServer(version, cliCommands)

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		cancel()
	}()

	// Run server over stdio
	if err := mcpServer.Run(ctx, &mcp.StdioTransport{}); err != nil {
		return nil, fmt.Errorf("MCP server error: %w", err)
	}

	return map[string]string{"status": "shutdown"}, nil
}
