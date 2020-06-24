package core

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
)

type BuildInfo struct {
	Version   *version.Version `json:"-"`
	BuildDate string           `json:"build_date"`
	GoVersion string           `json:"go_version"`
	GitBranch string           `json:"git_branch"`
	GitCommit string           `json:"git_commit"`
	GoArch    string           `json:"go_arch"`
	GoOS      string           `json:"go_os"`
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
	scwDisableCheckVersionEnv        = "SCW_DISABLE_CHECK_VERSION"
	latestVersionFileURL             = "https://scw-devtools.s3.nl-ams.scw.cloud/scw-cli-v2-version"
	latestVersionUpdateFileLocalName = "latest-cli-version"
	latestVersionRequestTimeout      = 1 * time.Second
	userAgentPrefix                  = "scaleway-cli"
)

// IsRelease returns true when the version of the CLI is an official release:
// - version must be non-empty (exclude tests)
// - version must not contain metadata (e.g. '+dev')
func (b *BuildInfo) IsRelease() bool {
	return b.Version != nil && b.Version.Metadata() == ""
}

func (b *BuildInfo) GetUserAgent() string {
	if b.Version != nil {
		return userAgentPrefix + "/" + b.Version.String()
	}
	return userAgentPrefix
}

func (b *BuildInfo) checkVersion(ctx context.Context) {
	if !b.IsRelease() || ExtractEnv(ctx, scwDisableCheckVersionEnv) == "true" {
		ExtractLogger(ctx).Debug("skipping check version")
		return
	}

	latestVersionUpdateFilePath := getLatestVersionUpdateFilePath(ExtractCacheDir(ctx))

	// do nothing if last refresh at during the last 24h
	if wasFileModifiedLast24h(latestVersionUpdateFilePath) {
		ExtractLogger(ctx).Debug("version was already checked during past 24 hours")
		return
	}

	// do nothing if we cannot create the file
	err := createAndCloseFile(latestVersionUpdateFilePath)
	if err != nil {
		ExtractLogger(ctx).Debug(err.Error())
		return
	}

	// pull latest version
	latestVersion, err := getLatestVersion(ExtractHTTPClient(ctx))
	if err != nil {
		ExtractLogger(ctx).Debugf("failed to retrieve latest version: %s\n", err)
		return
	}

	if b.Version.LessThan(latestVersion) {
		ExtractLogger(ctx).Warningf("a new version of scw is available (%s), beware that you are currently running %s\n", latestVersion, b.Version)
	} else {
		ExtractLogger(ctx).Debugf("version is up to date (%s)\n", b.Version)
	}
}

func getLatestVersionUpdateFilePath(cacheDir string) string {
	return filepath.Join(cacheDir, latestVersionUpdateFileLocalName)
}

// getLatestVersion attempt to read the latest version of the remote file at latestVersionFileURL.
func getLatestVersion(client *http.Client) (*version.Version, error) {
	ctx, cancelTimeout := context.WithTimeout(context.Background(), latestVersionRequestTimeout)
	defer cancelTimeout()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, latestVersionFileURL, nil)
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
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return version.NewSemver(strings.Trim(string(body), "\n"))
}

// wasFileModifiedLast24h checks whether the file has been updated during last 24 hours.
func wasFileModifiedLast24h(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	yesterday := time.Now().AddDate(0, 0, -1)
	lastUpdate := stat.ModTime()
	return lastUpdate.After(yesterday)
}

// createAndCloseFile creates a file and closes it. It returns true on succeed, false on failure.
func createAndCloseFile(path string) error {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		return fmt.Errorf("failed creating path %s: %s", path, err)
	}
	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return fmt.Errorf("failed creating file %s: %s", path, err)
	}

	return newFile.Close()
}
