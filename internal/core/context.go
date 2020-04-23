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

	ProfileFlag    string
	ConfigPathFlag string

	BuildInfo    *BuildInfo
	Client       *scw.Client
	Printer      printer.Printer
	Commands     *Commands
	OverrideEnv  map[string]string
	OverrideExec OverrideExecFunc

	command                     *Command
	stdout                      io.Writer
	stderr                      io.Writer
	result                      interface{}
	isClientFromBootstrapConfig bool
}

type contextKey int

const (
	metaContextKey contextKey = iota
)

// injectMeta creates a new ctx based on the given one with injected meta and returns it.
func injectMeta(ctx context.Context, meta *meta) context.Context {
	return context.WithValue(ctx, metaContextKey, meta)
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

func ExtractProfileName(ctx context.Context) string {
	if extractMeta(ctx).ProfileFlag != "" {
		return extractMeta(ctx).ProfileFlag
	}
	return ExtractEnv(ctx, scw.ScwActiveProfileEnv)
}

func ExtractConfigPath(ctx context.Context) string {
	if extractMeta(ctx).ConfigPathFlag != "" {
		return extractMeta(ctx).ConfigPathFlag
	}
	return scw.GetConfigPath()
}

func ReloadClient(ctx context.Context) error {
	var err error
	meta := extractMeta(ctx)
	meta.Client, err = createClient(meta.BuildInfo, "")
	return err
}

func ExtractProfileName(ctx context.Context) string {
	return extractMeta(ctx).ProfileFlag
}
