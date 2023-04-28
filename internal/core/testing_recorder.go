package core

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"regexp"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
)

func cassetteRequestFilter(i *cassette.Interaction) error {
	delete(i.Request.Headers, "x-auth-token")
	delete(i.Request.Headers, "X-Auth-Token")
	i.Request.URL = regexp.MustCompile("organization_id=[0-9a-f-]{36}").ReplaceAllString(i.Request.URL, "organization_id=11111111-1111-1111-1111-111111111111")
	i.Request.URL = regexp.MustCompile(`api\.scaleway\.com/account/v1/tokens/[0-9a-f-]{36}`).ReplaceAllString(i.Request.URL, "api.scaleway.com/account/v1/tokens/11111111-1111-1111-1111-111111111111")
	i.Request.URL = regexp.MustCompile(`api\.scaleway\.com/iam/v1alpha1/api-keys/SCW[0-9A-Z]{17}`).ReplaceAllString(i.Request.URL, "api.scaleway.com/iam/v1alpha1/api-keys/SCWXXXXXXXXXXXXXXXXX")

	return nil
}

func cassetteResponseFilter(i *cassette.Interaction) error {
	i.Response.Body = regexp.MustCompile(`"secret_key":"[0-9a-f-]{36}"`).ReplaceAllString(i.Response.Body, `"secret_key":"11111111-1111-1111-1111-111111111111"`)

	// Buildpacks
	i.Request.URL = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).ReplaceAllString(i.Request.URL, "pack.local%2Fbuilder%2F11111111111111111111")
	i.Request.URL = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).ReplaceAllString(i.Request.URL, "pack.local/builder/11111111111111111111")

	i.Request.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")
	i.Response.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")

	return nil
}

const (
	windowDockerEngine = "//./pipe/docker_engine"
	unixDockerEngine   = "/var/run/docker.sock"
)

func cassetteMatcher(r *http.Request, i cassette.Request) bool {
	// Docker
	if r.URL.Host == windowDockerEngine {
		r.URL.Host = unixDockerEngine
	}

	r.URL.Host = regexp.MustCompile(`/docker_engine`).ReplaceAllString(r.URL.Host, "/docker.sock")

	if r.URL.Scheme == "npipe" {
		r.URL.Scheme = "http"
	}

	r.URL.RawQuery = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).ReplaceAllString(r.URL.RawQuery, "pack.local%2Fbuilder%2F11111111111111111111")
	r.URL.Path = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).ReplaceAllString(r.URL.Path, "pack.local/builder/11111111111111111111")

	// Read body
	if r.Body != nil && r.Body != http.NoBody {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("failed to read request body")
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	return cassette.DefaultMatcher(r, i)
}

// getHTTPRecoder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP request must be made to run the tests.
//
// It is important to call add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func getHTTPRecoder(t *testing.T, update bool) (client *http.Client, cleanup func(), err error) {
	recorderMode := recorder.ModeReplaying
	if update {
		recorderMode = recorder.ModeRecording
	}

	// Setup recorder and scw client
	r, err := recorder.NewAsMode(getTestFilePath(t, ".cassette"), recorderMode, &SocketPassthroughTransport{})
	if err != nil {
		return nil, nil, err
	}

	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(cassetteRequestFilter)

	// Remove secrets from response
	r.AddSaveFilter(cassetteResponseFilter)

	r.SetMatcher(cassetteMatcher)

	return &http.Client{Transport: &retryableHTTPTransport{transport: r}}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
