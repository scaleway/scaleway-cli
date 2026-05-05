package core

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/scaleway/scaleway-cli/v2/internal/alias"
	cliConfig "github.com/scaleway/scaleway-cli/v2/internal/config"
	"github.com/scaleway/scaleway-cli/v2/internal/platform"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Meta store globally available variables like sdk client or global Flags.
type Meta struct {
	BinaryName string

	ProfileFlag    string
	ConfigPathFlag string
	Logger         *Logger

	BuildInfo    *BuildInfo
	Client       *scw.Client
	Commands     *Commands
	OverrideEnv  map[string]string
	OverrideExec OverrideExecFunc
	CliConfig    *cliConfig.Config
	Platform     platform.Platform

	command                     *Command
	stdout                      io.Writer
	stderr                      io.Writer
	stdin                       io.Reader
	result                      any
	httpClient                  *http.Client
	isClientFromBootstrapConfig bool
	BetaMode                    bool
}

type contextKey int

const (
	metaContextKey contextKey = iota
)

// InjectMeta creates a new ctx based on the given one with injected meta and returns it.
func InjectMeta(ctx context.Context, meta *Meta) context.Context {
	return context.WithValue(ctx, metaContextKey, meta)
}

// extractMeta extracts Meta from a given context.
func extractMeta(ctx context.Context) *Meta {
	if meta, ok := ctx.Value(metaContextKey).(*Meta); ok {
		return meta
	}
	return nil
}

// injectSDKConfig add config to a Meta context
func InjectConfig(ctx context.Context, config *scw.Config) {
	meta := extractMeta(ctx)
	if meta != nil && meta.Platform != nil {
		meta.Platform.SetScwConfig(config)
	}
}

func extractConfig(ctx context.Context) *scw.Config {
	m := extractMeta(ctx)
	if m != nil && m.Platform != nil {
		return m.Platform.ScwConfig()
	}

	return nil
}

func ExtractCommands(ctx context.Context) *Commands {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.Commands
}

func ExtractCliConfig(ctx context.Context) *cliConfig.Config {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.CliConfig
}

func ExtractAliases(ctx context.Context) *alias.Config {
	cliConfig := ExtractCliConfig(ctx)
	if cliConfig == nil {
		return nil
	}
	return cliConfig.Alias
}

func GetOrganizationIDFromContext(ctx context.Context) string {
	client := ExtractClient(ctx)
	if client == nil {
		return ""
	}
	organizationID, _ := client.GetDefaultOrganizationID()

	return organizationID
}

func GetProjectIDFromContext(ctx context.Context) string {
	client := ExtractClient(ctx)
	if client == nil {
		return ""
	}
	projectID, _ := client.GetDefaultProjectID()

	return projectID
}

func ExtractClient(ctx context.Context) *scw.Client {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.Client
}

func ExtractLogger(ctx context.Context) *Logger {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.Logger
}

func ExtractBuildInfo(ctx context.Context) *BuildInfo {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.BuildInfo
}

func ExtractBetaMode(ctx context.Context) bool {
	meta := extractMeta(ctx)
	if meta == nil {
		return false
	}
	return meta.BetaMode
}

func ExtractEnv(ctx context.Context, envKey string) string {
	meta := extractMeta(ctx)
	if meta != nil {
		if value, exist := meta.OverrideEnv[envKey]; exist {
			return value
		}
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
	meta := extractMeta(ctx)
	if meta == nil {
		return ""
	}
	return meta.BinaryName
}

func ExtractStdin(ctx context.Context) io.Reader {
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.stdin
}

func ExtractProfileName(ctx context.Context) string {
	meta := extractMeta(ctx)
	// Handle profile flag -p
	if meta != nil && meta.ProfileFlag != "" {
		return meta.ProfileFlag
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
	meta := extractMeta(ctx)
	if meta == nil {
		return nil
	}
	return meta.httpClient
}

func ExtractConfigPath(ctx context.Context) string {
	meta := extractMeta(ctx)
	if meta == nil {
		return scw.GetConfigPath()
	}
	if meta.ConfigPathFlag != "" {
		return meta.ConfigPathFlag
	}
	// This is only useful for test when we override home environment variable
	if home := meta.OverrideEnv["HOME"]; home != "" {
		return path.Join(home, ".config", "scw", "config.yaml")
	}

	return scw.GetConfigPath()
}

func ExtractCliConfigPath(ctx context.Context) string {
	meta := extractMeta(ctx)
	if meta == nil {
		configPath, _ := cliConfig.FilePath()
		return configPath
	}
	// This is only useful for test when we override home environment variable
	if home := meta.OverrideEnv["HOME"]; home != "" {
		return path.Join(home, ".config", "scw", cliConfig.DefaultConfigFileName)
	}
	configPath, _ := cliConfig.FilePath()

	return configPath
}

func ReloadClient(ctx context.Context) error {
	meta := extractMeta(ctx)
	if meta == nil {
		return fmt.Errorf("cannot reload client: meta not found in context")
	}
	var err error
	meta.Client, err = meta.Platform.CreateClient(
		meta.httpClient,
		ExtractConfigPath(ctx),
		ExtractProfileName(ctx),
	)

	return err
}

func ExtractConfigPathFlag(ctx context.Context) string {
	meta := extractMeta(ctx)
	if meta == nil {
		return ""
	}
	return meta.ConfigPathFlag
}

func ExtractProfileFlag(ctx context.Context) string {
	meta := extractMeta(ctx)
	if meta == nil {
		return ""
	}
	return meta.ProfileFlag
}

// GetDocGenContext returns a minimal context that can be used by scw-doc-gen
func GetDocGenContext() context.Context {
	ctx := context.Background()
	client, _ := scw.NewClient(
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)

	return InjectMeta(ctx, &Meta{
		Client: client,
	})
}
