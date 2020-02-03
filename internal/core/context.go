package core

import (
	"context"
	"io"

	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// meta store globally available variables like sdk client or global Flags.
type meta struct {
	ProfileFlag     string
	DebugModeFlag   bool
	PrinterTypeFlag printer.Type

	BuildInfo *BuildInfo
	Client    *scw.Client
	Printer   printer.Printer
	Commands  *Commands

	command *Command
	stdout  io.Writer
	stderr  io.Writer
	result  interface{}
}

type contextKey int

const (
	metaContextKey contextKey = iota
)

// newMetaContext creates a new ctx with injected meta and returns it.
func newMetaContext(meta *meta) context.Context {
	return context.WithValue(context.Background(), metaContextKey, meta)
}

// extractMeta extracts meta from a given context.
func extractMeta(ctx context.Context) *meta {
	return ctx.Value(metaContextKey).(*meta)
}

func ExtractCommands(ctx context.Context) *Commands {
	return extractMeta(ctx).Commands
}

func GetOrganizationIDFromContext(ctx context.Context) (organizationID string) {
	client := ExtractClient(ctx)
	organizationID, exists := client.GetDefaultOrganizationID()
	if !exists {
		panic("no default organization ID found")
	}
	return organizationID
}

func ExtractClient(ctx context.Context) *scw.Client {
	return extractMeta(ctx).Client
}

func ExtractBuildInfo(ctx context.Context) *BuildInfo {
	return extractMeta(ctx).BuildInfo
}
