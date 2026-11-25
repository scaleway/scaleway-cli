//go:build unix

package instance_test

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func TestSSHFetchKeysInterrupted(t *testing.T) {
	api := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprint(w, `{"total_count":2,"servers":[
   {"name":"first","state":"running","public_ips":[{"address":"192.0.2.1"}]},
   {"name":"second","state":"running","public_ips":[{"address":"192.0.2.2"}]}
  ]}`)
	}))
	defer api.Close()
	client, err := scw.NewClient(
		scw.WithoutAuth(),
		scw.WithAPIURL(api.URL),
		scw.WithDefaultZone(scw.ZoneFrPar1),
	)
	require.NoError(t, err)

	for _, tc := range []struct {
		signal string
		code   int
	}{{"INT", 130}, {"TERM", 143}} {
		t.Run(tc.signal, func(t *testing.T) {
			dir := t.TempDir()
			trace := filepath.Join(dir, "trace")
			// The child signals only the CLI. Finishing the trap proves the CLI
			// forwarded that signal and waited for SSH to shut down.
			const ssh = `#!/bin/sh
echo "$1" >> "$SSH_TEST_TRACE"
trap 'echo finished >> "$SSH_TEST_TRACE"; exit 255' "$SSH_TEST_SIGNAL"
kill -s "$SSH_TEST_SIGNAL" "$PPID"
sleep 1
exit 1
`
			require.NoError(t, os.WriteFile(filepath.Join(dir, "ssh"), []byte(ssh), 0o700))
			t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
			t.Setenv("SSH_TEST_SIGNAL", tc.signal)
			t.Setenv("SSH_TEST_TRACE", trace)
			var stderr bytes.Buffer
			code, _, err := core.Bootstrap(t.Context(), &core.BootstrapConfig{
				Args:             []string{"scw", "instance", "ssh", "fetch-keys", "zone=fr-par-1"},
				Commands:         instance.GetCommands(),
				Client:           client,
				Stdout:           io.Discard,
				Stderr:           &stderr,
				Platform:         terminal.NewPlatform("test"),
				BuildInfo:        &core.BuildInfo{},
				DisableTelemetry: true,
				OverrideEnv: map[string]string{
					"HOME":                      dir,
					"SCW_DISABLE_CHECK_VERSION": "true",
				},
			})
			require.Error(t, err)
			require.Equal(t, tc.code, code, stderr.String())
			got, err := os.ReadFile(trace)
			require.NoError(t, err)
			require.Equal(t, "192.0.2.1\nfinished\n", string(got))
		})
	}
}
