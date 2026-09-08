package main

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	testInjectedDate     = "2025-06-26T08:16:28AM"
	testInjectedCommit   = "abcdef1"
	testVCSTime          = "2026-09-07T15:09:20Z"
	testVCSRevision      = "bede29c38b66ef43eb38ef3c08aea75c601dd7fa"
	testVCSShortRevision = "bede29c"
	testShortRevision    = "abc"
)

func vcsBuildInfo() *debug.BuildInfo {
	return &debug.BuildInfo{
		Settings: []debug.BuildSetting{
			{Key: "vcs", Value: "git"},
			{Key: vcsRevisionSetting, Value: testVCSRevision},
			{Key: vcsTimeSetting, Value: testVCSTime},
			{Key: "vcs.modified", Value: "false"},
		},
	}
}

func Test_buildDate(t *testing.T) {
	for _, tt := range []struct {
		name       string
		injected   string
		buildInfos *debug.BuildInfo
		expected   string
	}{
		{
			name:       "injected at link time wins over the embedded stamp",
			injected:   testInjectedDate,
			buildInfos: vcsBuildInfo(),
			expected:   testInjectedDate,
		},
		{
			name:       "falls back to the embedded commit date",
			injected:   unknownBuildValue,
			buildInfos: vcsBuildInfo(),
			expected:   testVCSTime,
		},
		{
			name:       "empty injected value falls back too",
			injected:   "",
			buildInfos: vcsBuildInfo(),
			expected:   testVCSTime,
		},
		{
			name:       "stays unknown without build info",
			injected:   unknownBuildValue,
			buildInfos: nil,
			expected:   unknownBuildValue,
		},
		{
			name:       "stays unknown when the binary carries no VCS stamp",
			injected:   unknownBuildValue,
			buildInfos: &debug.BuildInfo{},
			expected:   unknownBuildValue,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, buildDate(tt.injected, tt.buildInfos))
		})
	}
}

func Test_gitCommit(t *testing.T) {
	for _, tt := range []struct {
		name       string
		injected   string
		buildInfos *debug.BuildInfo
		expected   string
	}{
		{
			name:       "injected at link time wins over the embedded stamp",
			injected:   testInjectedCommit,
			buildInfos: vcsBuildInfo(),
			expected:   testInjectedCommit,
		},
		{
			name:       "falls back to the embedded revision, shortened",
			injected:   unknownBuildValue,
			buildInfos: vcsBuildInfo(),
			expected:   testVCSShortRevision,
		},
		{
			name:       "empty injected value falls back too",
			injected:   "",
			buildInfos: vcsBuildInfo(),
			expected:   testVCSShortRevision,
		},
		{
			name:       "stays unknown without build info",
			injected:   unknownBuildValue,
			buildInfos: nil,
			expected:   unknownBuildValue,
		},
		{
			name:       "stays unknown when the binary carries no VCS stamp",
			injected:   unknownBuildValue,
			buildInfos: &debug.BuildInfo{},
			expected:   unknownBuildValue,
		},
		{
			name:     "a revision shorter than the short form is kept as is",
			injected: unknownBuildValue,
			buildInfos: &debug.BuildInfo{
				Settings: []debug.BuildSetting{{Key: vcsRevisionSetting, Value: testShortRevision}},
			},
			expected: testShortRevision,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, gitCommit(tt.injected, tt.buildInfos))
		})
	}
}

func Test_newBuildInfo(t *testing.T) {
	t.Run("resolves the embedded stamps into the reported build info", func(t *testing.T) {
		buildInfo := newBuildInfo(vcsBuildInfo())

		assert.Equal(t, testVCSTime, buildInfo.BuildDate)
		assert.Equal(t, testVCSShortRevision, buildInfo.GitCommit)
		// The Go toolchain embeds no branch name, so this one stays unknown.
		assert.Equal(t, unknownBuildValue, buildInfo.GitBranch)
		assert.Equal(t, GoVersion, buildInfo.GoVersion)
		assert.Equal(t, userAgentPrefix, buildInfo.UserAgentPrefix)
	})

	t.Run("stays unknown when the binary carries no VCS stamp", func(t *testing.T) {
		buildInfo := newBuildInfo(&debug.BuildInfo{})

		assert.Equal(t, unknownBuildValue, buildInfo.BuildDate)
		assert.Equal(t, unknownBuildValue, buildInfo.GitCommit)
		assert.Equal(t, unknownBuildValue, buildInfo.GitBranch)
	})
}
