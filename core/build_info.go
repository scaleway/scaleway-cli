package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

type BuildInfo struct {
	Version         *version.Version `json:"-"`
	BuildDate       string           `json:"build_date"`
	GoVersion       string           `json:"go_version"`
	GitBranch       string           `json:"git_branch"`
	GitCommit       string           `json:"git_commit"`
	GoArch          string           `json:"go_arch"`
	GoOS            string           `json:"go_os"`
	UserAgentPrefix string           `json:"user_agent_prefix"`
}

func (b *BuildInfo) MarshalJSON() ([]byte, error) {
	type Tmp BuildInfo

	return json.Marshal(
		struct {
			Tmp
			Version string `json:"version"`
		}{Tmp(*b), b.Version.String()},
	)
}

const (
	scwDisableCheckVersionEnv   = "SCW_DISABLE_CHECK_VERSION"
	latestGithubReleaseURL      = "https://api.github.com/repos/scaleway/scaleway-cli/releases/latest"
	latestVersionRequestTimeout = 1 * time.Second
)

// IsRelease returns true when the version of the CLI is an official release:
// - version must be non-empty (exclude tests)
// - version must not contain metadata (e.g. '+dev')
func (b *BuildInfo) IsRelease() bool {
	return b.Version != nil && b.Version.Metadata() == ""
}

func (b *BuildInfo) GetUserAgent() string {
	if b.Version != nil {
		return b.UserAgentPrefix + "/" + b.Version.String()
	}

	return b.UserAgentPrefix
}

func (b *BuildInfo) Tags() map[string]string {
	return map[string]string{
		"version":    b.Version.String(),
		"go_arch":    b.GoArch,
		"go_os":      b.GoOS,
		"go_version": b.GoVersion,
	}
}

func (b *BuildInfo) checkVersion(ctx context.Context) {
	if !b.IsRelease() || ExtractEnv(ctx, scwDisableCheckVersionEnv) == "true" {
		ExtractLogger(ctx).Debug("skipping check version")

		return
	}

	// pull latest version
	latestVersion, err := getLatestVersion(ExtractHTTPClient(ctx))
	if err != nil {
		ExtractLogger(ctx).Debugf("failed to retrieve latest version: %s\n", err)

		return
	}

	if b.Version.LessThan(latestVersion) {
		ExtractLogger(
			ctx,
		).Warningf("A new version of scw is available (%s), beware that you are currently running %s\n", latestVersion, b.Version)
	} else {
		ExtractLogger(ctx).Debugf("version is up to date (%s)\n", b.Version)
	}
}

// getLatestVersion attempt to read the latest version of the remote file at latestVersionFileURL.
func getLatestVersion(client *http.Client) (*version.Version, error) {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), latestVersionRequestTimeout)
	defer cancelTimeout()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestGithubReleaseURL, nil)
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	jsonBody := struct {
		TagName string `json:"tag_name"`
	}{}

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal version from remote: %w", err)
	}

	return version.NewSemver(strings.TrimPrefix(jsonBody.TagName, "v"))
}
