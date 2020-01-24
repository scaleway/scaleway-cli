package core

import (
	"context"
	"io"

	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// meta store globally available variables like sdk client or global Flags.
type meta struct {
	AccessKeyFlag   string
	SecretKeyFlag   string
	ProfileFlag     string
	DebugModeFlag   bool
	PrinterTypeFlag printer.Type

	BuildInfo *BuildInfo
	Client    *scw.Client
	Printer   printer.Printer

	stdout io.Writer
	stderr io.Writer
}

type contextKey int

const (
	metaContextKey contextKey = iota
	commandsContextKey
	resultContextKey
)

// injectMeta creates a child of ctx with injected meta and returns it.
func injectMeta(ctx context.Context, meta *meta) context.Context {
	return context.WithValue(ctx, metaContextKey, meta)
}

// extractMeta extracts meta from a given context.
// It should not be used directly.
// Use ExtractClient(), ExtractBuildInfo() instead,
// or create new Extract__() methods if necessary.
//
// TODO: remove usage from cobraPreRunInitMeta()
func extractMeta(ctx context.Context) *meta {
	return ctx.Value(metaContextKey).(*meta)
}

func injectCommands(ctx context.Context, cmds *Commands) context.Context {
	return context.WithValue(ctx, commandsContextKey, cmds)
}

func ExtractCommands(ctx context.Context) *Commands {
	return ctx.Value(commandsContextKey).(*Commands)
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

func extractPrinter(ctx context.Context) printer.Printer {
	return extractMeta(ctx).Printer
}

func injectResultSetter(ctx context.Context, result *interface{}) context.Context {
	return context.WithValue(ctx, resultContextKey, func(r interface{}) {
		*result = r
	})
}

func setContextResult(ctx context.Context, result interface{}) {
	ctx.Value(resultContextKey).(func(interface{}))(result)
}
