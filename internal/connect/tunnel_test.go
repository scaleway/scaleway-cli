package connect_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test_SSHLocalForwardEarlyExit checks that a tunnel process that dies before
// the local port is open (e.g. an ssh authentication failure) is reported
// quickly, instead of after the full readiness timeout. The ready path and
// the Stop behavior are exercised by the redis cluster connect test against
// a fake ssh binary.
func Test_SSHLocalForwardEarlyExit(t *testing.T) {
	fakeDir := t.TempDir()
	script := filepath.Join(fakeDir, "ssh")
	require.NoError(t, os.WriteFile(script, []byte("#!/bin/sh\nexit 21\n"), 0o755))
	t.Setenv("PATH", fakeDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	start := time.Now()
	tunnel, err := connect.SSHLocalForward(os.Stderr, "1.2.3.4", 61000, "172.16.12.2", 6379)
	elapsed := time.Since(start)

	require.ErrorContains(t, err, "exited before becoming ready")
	assert.Nil(t, tunnel)
	assert.Less(
		t,
		elapsed,
		5*time.Second,
		"the failure must be reported quickly, not after the readiness timeout",
	)
}
