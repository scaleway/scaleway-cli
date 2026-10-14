package semver

import (
	"errors"

	"github.com/hashicorp/go-version"
)

// LatestVersion returns the highest version string in versions.
// Entries that do not parse as a version are ignored.
func LatestVersion(versions []string) (string, error) {
	latest, _ := version.NewVersion("0.0.0")
	latestStr := ""
	for _, s := range versions {
		parsed, err := version.NewVersion(s)
		if err != nil {
			continue
		}
		if parsed.GreaterThan(latest) {
			latest = parsed
			latestStr = s
		}
	}
	if latestStr == "" {
		return "", errors.New("no available version found")
	}

	return latestStr, nil
}
