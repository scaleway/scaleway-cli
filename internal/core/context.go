package core

import (
	"context"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// meta store globally available variables like sdk client or global Flags.
type meta struct {
	BinaryName string

	ProfileFlag    string
	ConfigPathFlag string
	Logger         *Logger

	BuildInfo    *BuildInfo
	Client       *scw.Client
	Commands     *Commands
	OverrideEnv  map[string]string
	OverrideExec OverrideExecFunc

	command                     *Command
	stdout                      io.Writer
	stderr                      io.Writer
	stdin                       io.Reader
	result                      interface{}
	httpClient                  *http.Client
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

func GetProjectIDFromContext(ctx context.Context) (projectID string) {
	client := ExtractClient(ctx)
	projectID, exists := client.GetDefaultProjectID()
	if !exists {
		panic("no default project ID found")
	}
	return projectID
}

func ExtractClient(ctx context.Context) *scw.Client {
	return extractMeta(ctx).Client
}

func ExtractLogger(ctx context.Context) *Logger {
	return extractMeta(ctx).Logger
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

func ExtractCacheDir(ctx context.Context) string {
	env := ExtractEnv(ctx, scw.ScwCacheDirEnv)
	if env != "" {
		return env
	}
	return scw.GetCacheDirectory()
}

func ExtractBinaryName(ctx context.Context) string {
	return extractMeta(ctx).BinaryName
}

func ExtractStdin(ctx context.Context) io.Reader {
	return extractMeta(ctx).stdin
}

func ExtractProfileName(ctx context.Context) string {
	// Handle profile flag -p
	if extractMeta(ctx).ProfileFlag != "" {
		return extractMeta(ctx).ProfileFlag
	}

	// Handle SCW_PROFILE env variable
	if env := ExtractEnv(ctx, scw.ScwActiveProfileEnv); env != "" {
		return env
	}

	// Handle active_profile in config file
	configPath := ExtractConfigPath(ctx)
	config, err := scw.LoadConfigFromPath(configPath)
	if err == nil && config.ActiveProfile != nil {
		return *config.ActiveProfile
	}

	// Return default profile name
	return scw.DefaultProfileName
}

func ExtractHTTPClient(ctx context.Context) *http.Client {
	return extractMeta(ctx).httpClient
}

func ExtractConfigPath(ctx context.Context) string {
	meta := extractMeta(ctx)
	if meta.ConfigPathFlag != "" {
		return extractMeta(ctx).ConfigPathFlag
	}
	// This is only useful for test when we override home environment variable
	if home := meta.OverrideEnv["HOME"]; home != "" {
		return path.Join(home, ".config", "scw", "config.yaml")
	}
	return scw.GetConfigPath()
}

func ReloadClient(ctx context.Context) error {
	var err error
	meta := extractMeta(ctx)
	// if client is from bootstrap we are probably running test
	// if we reload the client we loose the cassette recorder
	if meta.isClientFromBootstrapConfig {
		return nil
	}
	meta.Client, err = createClient(meta.httpClient, meta.BuildInfo, ExtractProfileName(ctx))
	return err
}

func ExtractConfigPathFlag(ctx context.Context) string {
	return extractMeta(ctx).ConfigPathFlag
}

func ExtractProfileFlag(ctx context.Context) string {
	return extractMeta(ctx).ProfileFlag
}

// GetDocGenContext reuturn a minimal context that can be used by scw-doc-gen
func GetDocGenContext() context.Context {
	ctx := context.Background()
	client, _ := scw.NewClient(
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	return injectMeta(ctx, &meta{
		Client: client,
	})
}
