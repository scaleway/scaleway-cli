package platform

import (
	"net/http"

	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// Platform defines an environment running the CLI
// It can be the implementation to run in a terminal
// Or the implementation to run in a browser (used for wasm/js build)
type Platform interface {
	// CreateClient returns a valid client for the current platform
	CreateClient(client *http.Client, configPath string, profileName string) (*scw.Client, error)

	// ScwConfig returns a scaleway config if available, can be nil
	// TODO: remove if possible, currently used in profile completion
	ScwConfig() *scw.Config
	// SetScwConfig set the stored config, useful for testing purpose
	SetScwConfig(cfg *scw.Config)
}

func NewDefault(useragent string) *terminal.Platform {
	return &terminal.Platform{
		UserAgent: useragent,
	}
}
