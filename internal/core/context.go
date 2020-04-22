package core

import (
	"context"
	"io"
	"os"

	"github.com/scaleway/scaleway-cli/internal/printer"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// meta store globally available variables like sdk client or global Flags.
type meta struct {
	BinaryName string

	ProfileFlag     string
	DebugModeFlag   bool
	PrinterTypeFlag printer.Type

	BuildInfo   *BuildInfo
	Client      *scw.Client
	Printer     printer.Printer
	Commands    *Commands
	OverrideEnv map[string]string

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

func ExtractClientOrCreate(ctx context.Context) (*scw.Client, error) {
	client := ExtractClient(ctx)
	if client == nil {
		var err error
		client, err = createClient(nil)
		if err != nil {
			return nil, err
		}
	}
	return client, nil
}

func ExtractBuildInfo(ctx context.Context) *BuildInfo {
	return extractMeta(ctx).BuildInfo
}

func ExtractEnv(ctx context.Context, envKey string) string {
	meta := extractMeta(ctx)
	if value, exist := meta.OverrideEnv[envKey]; exist {
		return value
	}

	if envKey == "HOME" {
		homeDir, _ := os.UserHomeDir()
		return homeDir
	}

	return os.Getenv(envKey)
}

func ExtractUserHomeDir(ctx context.Context) string {
	return ExtractEnv(ctx, "HOME")
}

func ExtractBinaryName(ctx context.Context) string {
	return extractMeta(ctx).BinaryName
}
