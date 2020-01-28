package core

import (
	"context"
	"io"

	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/spf13/cobra"
)

// meta store globally available variables like sdk client or global Flags.
type meta struct {
	AccessKeyFlag   string
	SecretKeyFlag   string
	ProfileFlag     string
	DebugModeFlag   bool
	PrinterTypeFlag printer.Type

	BuildInfo  *BuildInfo
	Client     *scw.Client
	Printer    printer.Printer
	Commands   *Commands
	RunCommand *cobra.Command

	stdout io.Writer
	stderr io.Writer
	result *interface{}
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
// It should not be used directly.
// Use ExtractClient(), ExtractBuildInfo() instead,
// or create new Extract__() methods if necessary.
//
// TODO: remove usage from cobraPreRunInitMeta()
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

func extractPrinter(ctx context.Context) printer.Printer {
	return extractMeta(ctx).Printer
}

func setContextResult(ctx context.Context, result interface{}) {
	*extractMeta(ctx).result = result
}
