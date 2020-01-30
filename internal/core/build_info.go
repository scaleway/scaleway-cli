package core

import (
	"strings"
)

type BuildInfo struct {
	Version   string
	BuildDate string
	GoVersion string
	GitBranch string
	GitCommit string
	GoArch    string
	GoOS      string
}

// IsRelease returns true when the version of the CLI is an official release:
// - version must be non-empty (exclude tests)
// - version must not contain label (e.g. '+dev')
func (b *BuildInfo) IsRelease() bool {
	return b.Version != "" && !strings.Contains(b.Version, "+")
}
