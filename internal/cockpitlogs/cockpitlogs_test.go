package cockpitlogs_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/cockpitlogs"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ParseValue(t *testing.T) {
	entry, err := cockpitlogs.ParseValue(
		[]string{"1000", `{"resource_instance":"node-1","message":"hello"}`},
	)
	require.NoError(t, err)
	assert.Equal(t, cockpitlogs.Entry{
		Timestamp:        time.Unix(0, 1000),
		ResourceInstance: "node-1",
		Message:          "hello",
	}, entry)

	entry, err = cockpitlogs.ParseValue([]string{"1000", "LOG:  checkpoint starting"})
	require.NoError(t, err)
	assert.Equal(t, "LOG:  checkpoint starting", entry.Message)

	_, err = cockpitlogs.ParseValue([]string{"1000"})
	require.Error(t, err)
}

// Test_Follow runs `-f` against a fake Cockpit API + Loki: the first query returns the newest
// lines, the next one must start right after the last line seen and read forward.
func Test_Follow(t *testing.T) {
	lokiCalls := 0
	tokenDeleted := false

	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		switch {
		case r.URL.Path == "/metadata":
			_, _ = w.Write([]byte(`{"domain":"scw.cloud"}`))
		case r.URL.Path == "/cockpit/v1/regions/fr-par/data-sources":
			_, _ = w.Write([]byte(`{"data_sources":[{"url":"` + srv.URL + `"}],"total_count":1}`))
		case r.Method == http.MethodPost && r.URL.Path == "/cockpit/v1/regions/fr-par/tokens":
			_, _ = w.Write([]byte(`{"id":"token-id","secret_key":"secret","region":"fr-par"}`))
		case r.Method == http.MethodDelete && r.URL.Path == "/cockpit/v1/regions/fr-par/tokens/token-id":
			tokenDeleted = true
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/loki/api/v1/query_range":
			lokiCalls++
			assert.Equal(t, "secret", r.Header.Get("X-Token"))
			assert.Equal(
				t,
				`{resource_type="rdb_instance_postgresql", resource_id="instance-id"}`,
				r.URL.Query().Get("query"),
			)

			if lokiCalls == 1 {
				assert.Empty(t, r.URL.Query().Get("direction"))
				_, _ = w.Write([]byte(lokiBody(`["2000","second"],["1000","first"]`)))

				return
			}

			start, err := time.Parse(time.RFC3339Nano, r.URL.Query().Get("start"))
			assert.NoError(t, err)
			assert.True(t, start.Equal(time.Unix(0, 2001)), "start %s", start)
			assert.Equal(t, "forward", r.URL.Query().Get("direction"))
			_, _ = w.Write([]byte(lokiBody(`["3000","third"]`)))
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL)
		}
	}))
	defer srv.Close()

	client, err := scw.NewClient(
		scw.WithAPIURL(srv.URL),
		scw.WithHTTPClient(srv.Client()),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultProjectID("11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(t.Context())
	defer cancel()

	printed := 0
	resource := cockpitlogs.Resource{
		Type:   "rdb_instance_postgresql",
		ID:     "instance-id",
		Region: scw.RegionFrPar,
	}
	commands := core.NewCommands(&core.Command{
		Namespace:            "test",
		Resource:             "logs",
		AllowAnonymousClient: true,
		DisableAfterChecks:   true,
		ArgsType:             reflect.TypeFor[struct{}](),
		Run:                  func(context.Context, any) (any, error) { return nil, nil },
		FollowFunc: func(ctx context.Context, _ any, emit func(any) error) error {
			return cockpitlogs.Follow(ctx, resource, time.Hour, func(data any) error {
				printed += len(data.(cockpitlogs.Lines))
				if printed == 3 {
					cancel()
				}

				return emit(data)
			})
		},
	})

	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	buildInfo := &core.BuildInfo{Version: version.Must(version.NewSemver("v0.0.0+test"))}
	exitCode, _, err := core.Bootstrap(ctx, &core.BootstrapConfig{
		Args:             []string{"scw", "test", "logs", "-f"},
		Commands:         commands,
		BuildInfo:        buildInfo,
		Stdout:           stdout,
		Stderr:           stderr,
		Client:           client,
		HTTPClient:       srv.Client(),
		DisableTelemetry: true,
		DisableAliases:   true,
		OverrideEnv:      map[string]string{"HOME": t.TempDir()},
		Platform:         terminal.NewPlatform(buildInfo.GetUserAgent()),
	})
	require.NoError(t, err, stderr.String())
	assert.Equal(t, 0, exitCode)
	assert.Equal(t, 2, lokiCalls)
	assert.True(t, tokenDeleted, "the cockpit token must be deleted when follow stops")

	expected := make([]string, 0, 3)
	for i, msg := range []string{"first", "second", "third"} {
		ts := time.Unix(0, int64(i+1)*1000).Format(time.RFC3339)
		expected = append(expected, ts+" node-1 "+msg)
	}
	assert.Equal(t, strings.Join(expected, "\n")+"\n", stdout.String())
}

func lokiBody(values string) string {
	return `{"data":{"result":[{"stream":{"instance":"node-1"},"values":[` + values + `]}]}}`
}
