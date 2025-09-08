package core

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

func cassetteRequestFilter(i *cassette.Interaction) error {
	delete(i.Request.Headers, "x-auth-token")
	delete(i.Request.Headers, "X-Auth-Token")
	orgIDRegex := regexp.MustCompile(`(.+)organization_id=[0-9a-f-]{36}(.+)`)
	tokenRegex := regexp.MustCompile(`^https://api\.scaleway\.com/account/v1/tokens/[0-9a-f-]{36}$`)
	accessKeyRegex := regexp.MustCompile(`(.+)?SCW[0-9A-Z]{17}(.+)?`)

	i.Request.URL = orgIDRegex.ReplaceAllString(
		i.Request.URL,
		"${1}organization_id=11111111-1111-1111-1111-111111111111${2}")
	i.Request.URL = tokenRegex.ReplaceAllString(
		i.Request.URL,
		"api.scaleway.com/account/v1/tokens/11111111-1111-1111-1111-111111111111")
	i.Request.URL = accessKeyRegex.ReplaceAllString(
		i.Request.URL,
		"${1}SCWXXXXXXXXXXXXXXXXX${2}")

	return nil
}

func cassetteResponseFilter(i *cassette.Interaction) error {
	i.Response.Body = regexp.MustCompile(`"secret_key":"[0-9a-f-]{36}"`).
		ReplaceAllString(i.Response.Body, `"secret_key":"11111111-1111-1111-1111-111111111111"`)

	// Buildpacks
	i.Request.URL = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).
		ReplaceAllString(i.Request.URL, "pack.local%2Fbuilder%2F11111111111111111111")
	i.Request.URL = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Request.URL, "pack.local/builder/11111111111111111111")

	i.Request.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")
	i.Response.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")

	return nil
}

const (
	windowDockerEngine      = "//./pipe/docker_engine"
	unixDockerEngine        = "/var/run/docker.sock"
	escapedUnixDockerEngine = "%2Fvar%2Frun%2Fdocker.sock"
)

func cassetteMatcher(r *http.Request, i cassette.Request) bool {
	// Docker
	if r.URL.Host == windowDockerEngine || r.URL.Host == "npipe://"+windowDockerEngine {
		r.URL.Host = unixDockerEngine
	}

	r.URL.RawQuery = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).
		ReplaceAllString(r.URL.RawQuery, "pack.local%2Fbuilder%2F11111111111111111111")
	r.URL.Path = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(r.URL.Path, "pack.local/builder/11111111111111111111")

	// Read body
	if r.Body != nil && r.Body != http.NoBody {
		reqBody, err := io.ReadAll(r.Body)
		if err != nil {
			log.Fatal("failed to read request body")
		}
		r.Body.Close()
		r.Body = io.NopCloser(bytes.NewBuffer(reqBody))
	}

	// Specific handling of s3 URLs
	// Url format is https://test-acc-scaleway-object-bucket-lifecycle-8445817190507446251.s3.fr-par.scw.cloud/?lifecycle=
	if strings.HasSuffix(r.URL.Host, "scw.cloud") {
		return customS3Matcher(r, i)
	}

	// Specific handling of Docker URLs
	// URLs are stored unescaped in the cassette but the matcher expects an escaped URL
	if r.URL.Host == unixDockerEngine {
		return customDockerMatcher(r, i)
	}

	return cassette.DefaultMatcher(r, i)
}

func customS3Matcher(r *http.Request, i cassette.Request) bool {
	expectedURL, _ := url.Parse(i.URL)
	actualURL := r.URL
	if !strings.HasSuffix(expectedURL.Host, "scw.cloud") {
		return false
	}

	actualS3Host := strings.Split(actualURL.Host, ".")
	expectedS3Host := strings.Split(expectedURL.Host, ".")
	if len(actualS3Host) < 1 || len(expectedS3Host) < 1 {
		return false
	}
	actualBucket := actualS3Host[0]
	expectedBucket := expectedS3Host[0]

	// Compare bucket names without the random number at the end
	if strings.Contains(actualBucket, "-") {
		actualBucket = actualBucket[:strings.LastIndex(actualBucket, "-")]
	}
	if strings.Contains(expectedBucket, "-") {
		expectedBucket = expectedBucket[:strings.LastIndex(expectedBucket, "-")]
	}
	if actualBucket != expectedBucket {
		return false
	}

	// Compare queries
	expectedURLValues := expectedURL.Query()
	actualURLValues := actualURL.Query()
	expectedURL.RawQuery = expectedURLValues.Encode()
	actualURL.RawQuery = actualURLValues.Encode()

	return r.Method == i.Method && r.URL.Path == expectedURL.Path &&
		actualURL.RawQuery == expectedURL.RawQuery
}

func customDockerMatcher(r *http.Request, i cassette.Request) bool {
	escapedRecordedURL := regexp.MustCompile(`http://`+unixDockerEngine+`(.+)?`).
		ReplaceAllString(
			i.URL,
			"http://"+escapedUnixDockerEngine+"${1}")

	return r.URL.String() == escapedRecordedURL
}

func unescapeDockerURL(i *cassette.Interaction) error {
	i.Request.URL = regexp.MustCompile(`http://`+escapedUnixDockerEngine+`(.+)?`).
		ReplaceAllString(
			i.Request.URL,
			"http://"+unixDockerEngine+"${1}")

	return nil
}

// getHTTPRecoder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP request must be made to run the tests.
//
// It is important to call add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func getHTTPRecoder(t *testing.T, update bool) (client *http.Client, cleanup func(), err error) {
	t.Helper()
	recorderMode := recorder.ModeReplayOnly
	if update {
		recorderMode = recorder.ModeRecordOnly
	}

	// Setup recorder and scw client
	r, err := recorder.NewWithOptions(&recorder.Options{
		CassetteName:       getTestFilePath(t, ".cassette"),
		Mode:               recorderMode,
		RealTransport:      &SocketPassthroughTransport{},
		SkipRequestLatency: true,
	})
	if err != nil {
		return nil, nil, err
	}

	// Starting with v3, go-vcr now calls net/url.Parse to build the interaction which results in an error for escaped
	// Docker URLs on Test_Deploy (container), so we need to unescape these paths on this step.
	r.AddHook(unescapeDockerURL, recorder.AfterCaptureHook)

	// Add a filter which removes Authorization headers from all requests:
	r.AddHook(cassetteRequestFilter, recorder.BeforeSaveHook)

	// Remove secrets from response
	r.AddHook(cassetteResponseFilter, recorder.BeforeSaveHook)

	r.SetMatcher(cassetteMatcher)

	return &http.Client{Transport: &retryableHTTPTransport{transport: r}}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
