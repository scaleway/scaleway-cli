package k8s_test

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const fakeKubectlScript = `#!/bin/sh
[ "$1" = "__complete" ] || exit 1
[ -n "${KUBECTL_FAKE_ARGS_FILE:-}" ] && printf '%s\n' "$@" > "$KUBECTL_FAKE_ARGS_FILE"
echo "pods"
echo "nodes"
echo "deployments"
echo ":0"
`

// runKubectlAutoComplete calls core.AutoComplete for
// "scw k8s cluster kubectl <args...> <prefix>" with a fake kubectl on the PATH.
func runKubectlAutoComplete(t *testing.T, args ...string) (core.AutocompleteSuggestions, []string) {
	t.Helper()

	tmpDir := t.TempDir()
	kubectlPath := filepath.Join(tmpDir, "kubectl")
	require.NoError(t, os.WriteFile(kubectlPath, []byte(fakeKubectlScript), 0o755))
	t.Setenv("PATH", tmpDir+string(os.PathListSeparator)+os.Getenv("PATH"))

	argsFile := filepath.Join(tmpDir, "args")
	t.Setenv("KUBECTL_FAKE_ARGS_FILE", argsFile)

	ctx := core.InjectMeta(context.Background(), &core.Meta{
		Commands: k8s.GetCommands(),
	})

	words := append([]string{"scw", "k8s", "cluster", "kubectl"}, args...)
	prefix := words[len(words)-1]
	leftWords := words[:len(words)-1]

	resp := core.AutoComplete(ctx, leftWords, prefix, nil)

	var received []string
	if data, err := os.ReadFile(argsFile); err == nil {
		received = strings.Split(strings.TrimSuffix(string(data), "\n"), "\n")
	}

	return resp.Suggestions, received
}

func Test_KubectlAutoComplete(t *testing.T) {
	const clusterID = "3313184d-3767-4535-902a-2bba7b7a4414"

	t.Run("simple", func(t *testing.T) {
		suggestions, received := runKubectlAutoComplete(t, clusterID, "get", "")
		assert.Equal(t, core.AutocompleteSuggestions{"deployments", "nodes", "pods"}, suggestions)
		// The cluster ID must not be forwarded, only the kubectl args and the prefix
		assert.Equal(t, []string{"__complete", "get", ""}, received)
	})

	t.Run("separator_and_flags", func(t *testing.T) {
		suggestions, received := runKubectlAutoComplete(t, clusterID, "--", "get", "pods", "-o", "")
		// "pods" was already typed, so scw's positional dedup drops it
		assert.Equal(t, core.AutocompleteSuggestions{"deployments", "nodes"}, suggestions)
		assert.Equal(t, []string{"__complete", "get", "pods", "-o", ""}, received)
	})

	t.Run("region_not_forwarded", func(t *testing.T) {
		_, received := runKubectlAutoComplete(t, "region=fr-par", clusterID, "get", "")
		assert.Equal(t, []string{"__complete", "get", ""}, received)
	})

	t.Run("no_cluster_id", func(t *testing.T) {
		suggestions, received := runKubectlAutoComplete(t, "")
		assert.Empty(t, suggestions)
		assert.Nil(t, received)
	})
}
