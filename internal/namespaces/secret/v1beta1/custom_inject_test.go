package secret_test

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
	secretSDK "github.com/scaleway/scaleway-sdk-go/api/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const injectSecretID = "11111111-1111-1111-1111-111111111111"

// newSecretAPIMock serves two secret versions: one by ID (latest) and one by path /my-app, name db, revision 2.
// It also returns the number of requests it received.
func newSecretAPIMock(t *testing.T) (*scw.Client, *atomic.Int32) {
	t.Helper()

	var requests atomic.Int32

	const base = "/secret-manager/v1beta1/regions/fr-par/"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)

		var data string

		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.URL.Path == base+"secrets/"+injectSecretID+"/versions/latest/access":
			data = "s3cr3t"
		case r.URL.Path == base+"secrets-by-path/versions/2/access" &&
			r.URL.Query().Get("secret_path") == "/my-app" &&
			r.URL.Query().Get("secret_name") == "db":
			data = `{"password":"hunter2"}`
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = io.WriteString(
				w,
				`{"type":"not_found","resource":"secret","message":"resource is not found"}`,
			)

			return
		}

		_ = json.NewEncoder(w).Encode(&secretSDK.AccessSecretVersionResponse{Data: []byte(data)})
	}))
	t.Cleanup(srv.Close)

	client, err := scw.NewClient(
		scw.WithHTTPClient(srv.Client()),
		scw.WithAPIURL(srv.URL),
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", injectSecretID),
	)
	require.NoError(t, err)

	return client, &requests
}

func Test_SecretInject(t *testing.T) {
	tests := []struct {
		name         string
		template     string
		fromFile     bool
		wantOut      string
		wantErr      string
		wantRequests int32
	}{
		{
			name:         "by-id-fetched-once",
			template:     `{{ secret "` + injectSecretID + `" }},{{ with true }}{{ secret "` + injectSecretID + `" }}{{ end }}`,
			wantOut:      "s3cr3t,s3cr3t",
			wantRequests: 1,
		},
		{
			name:         "unused-unknown-secret-is-prefetched-without-failing",
			template:     `{{ if false }}{{ secret "missing" }}{{ end }}ok`,
			wantOut:      "ok",
			wantRequests: 1,
		},
		{
			name:         "by-path-with-revision-and-json-field",
			template:     `{{ secret "my-app/db@2" | jsonField "password" }}`,
			fromFile:     true,
			wantOut:      "hunter2",
			wantRequests: 1,
		},
		{
			name:     "empty-name",
			template: `{{ secret "my-app/@2" }}`,
			wantErr:  "invalid secret reference",
		},
		{
			name:     "unknown-secret-outputs-nothing",
			template: `before {{ secret "missing" }}`,
			wantErr:  "accessing secret \"missing\"",
		},
		{
			name:     "invalid-template",
			template: `{{ secret `,
			wantErr:  "parsing template",
		},
	}

	for _, tt := range tests {
		client, requests := newSecretAPIMock(t)
		config := &core.TestConfig{
			Commands: secret.GetCommands(),
			Client:   client,
			Args:     []string{"scw", "secret", "inject"},
			Stdin:    strings.NewReader(tt.template),
			Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()

				if tt.wantErr != "" {
					require.ErrorContains(t, ctx.Err, tt.wantErr)
					assert.Empty(t, ctx.Stdout)

					return
				}

				require.NoError(t, ctx.Err)
				assert.Equal(t, core.RawResult(tt.wantOut), ctx.Result)
				assert.Equal(t, tt.wantRequests, requests.Load())
			},
		}
		// Args are rendered as Go templates by core.Test, so templates go through stdin or a file.
		if tt.fromFile {
			path := filepath.Join(t.TempDir(), "config.tpl")
			require.NoError(t, os.WriteFile(path, []byte(tt.template), 0o600))
			config.Args = append(config.Args, "template=@"+path)
		}

		t.Run(tt.name, core.Test(config))
	}
}
