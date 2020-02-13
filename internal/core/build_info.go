package core

import (
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type BuildInfo struct {
	Version   *version.Version
	BuildDate string
	GoVersion string
	GitBranch string
	GitCommit string
	GoArch    string
	GoOS      string
}

const (
	scwDisableCheckVersionEnv        = "SCW_DISABLE_CHECK_VERSION"
	latestVersionFileURL             = "https://scw-devtools.s3.nl-ams.scw.cloud/scw-cli-v2-version"
	latestVersionUpdateFileLocalName = "latest-cli-version"
	latestVersionRequestTimeout      = 1 * time.Second
)

// IsRelease returns true when the version of the CLI is an official release:
// - version must be non-empty (exclude tests)
// - version must not contain metadata (e.g. '+dev')
func (b *BuildInfo) IsRelease() bool {
	return b.Version != nil && b.Version.Metadata() == ""
}

func (b *BuildInfo) checkVersion() {
	if !b.IsRelease() || os.Getenv(scwDisableCheckVersionEnv) == "true" {
		logger.Debugf("skipping check version")
		return
	}

	latestVersionUpdateFilePath := getLatestVersionUpdateFilePath()

	// do nothing if last refresh at during the last 24h
	if wasFileModifiedLast24h(latestVersionUpdateFilePath) {
		logger.Debugf("version was already checked during past 24 hours")
		return
	}

	// do nothing if we cannot create the file
	if !createAndCloseFile(latestVersionUpdateFilePath) {
		return
	}

	// pull latest version
	latestVersion, err := getLatestVersion()
	if err != nil {
		logger.Debugf("failed to retrieve latest version: %s", err)
		return
	}

	if b.Version.LessThan(latestVersion) {
		logger.Warningf("a new version of scw is available (%s), beware that you are currently running %v", latestVersion, b.Version)
	} else {
		logger.Debugf("version is up to date (%s)", b.Version)
	}
}

func getLatestVersionUpdateFilePath() string {
	return filepath.Join(scw.GetCacheDirectory(), latestVersionUpdateFileLocalName)
}

// getLatestVersion attempt to read the latest version of the remote file at latestVersionFileURL.
func getLatestVersion() (*version.Version, error) {
	resp, err := (&http.Client{
		Timeout: latestVersionRequestTimeout,
	}).Get(latestVersionFileURL)
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
func createAndCloseFile(path string) bool {
	err := os.MkdirAll(filepath.Dir(path), 0700)
	if err != nil {
		logger.Debugf("failed creating path %s: %s", path, err)
	}
	newFile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		logger.Debugf("failed creating file %s: %s", path, err)
		return false
	}

	newFile.Close()
	return true
}
