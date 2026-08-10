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

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-sdk-go/vcr"
	"github.com/stretchr/testify/assert"
)

func cassetteRequestFilter(i *cassette.Interaction) error {
	delete(i.Request.Headers, "x-auth-token")
	delete(i.Request.Headers, "X-Auth-Token")
	orgIDRegex := regexp.MustCompile(`(.+)organization_id=[0-9a-f-]{36}(.+)`)
	tokenRegex := regexp.MustCompile(`^https://api\.scaleway\.com/account/v1/tokens/[0-9a-f-]{36}$`)

	i.URL = orgIDRegex.ReplaceAllString(
		i.URL,
		"${1}organization_id=11111111-1111-1111-1111-111111111111${2}")
	i.URL = tokenRegex.ReplaceAllString(
		i.URL,
		"api.scaleway.com/account/v1/tokens/11111111-1111-1111-1111-111111111111")

	return nil
}

func cassetteResponseFilter(i *cassette.Interaction) error {
	i.Response.Body = regexp.MustCompile(`"secret_key":"[0-9a-f-]{36}"`).
		ReplaceAllString(i.Response.Body, `"secret_key":"11111111-1111-1111-1111-111111111111"`)

	// Buildpacks
	i.URL = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).
		ReplaceAllString(i.URL, "pack.local%2Fbuilder%2F11111111111111111111")
	i.URL = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.URL, "pack.local/builder/11111111111111111111")

	i.Request.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")
	i.Response.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")

	return nil
}

const (
	windowDockerEngine = "//./pipe/docker_engine"
	unixDockerEngine   = "/var/run/docker.sock"
)

func applyQueryMatcherIgnore(r *http.Request, i *cassette.Request) bool {
	queryMatcherIgnore := []string(nil)

	if strings.Contains(r.URL.RawQuery, "project_id=11111111-1111-1111-1111-111111111111") {
		match, err := regexp.MatchString("project_id=[0-9a-f-]{36}", i.URL)
		if err != nil || !match {
			return false
		}

		queryMatcherIgnore = append(queryMatcherIgnore, "project_id")
	}

	if strings.Contains(r.URL.RawQuery, "end=") {
		if !strings.Contains(i.URL, "end=") {
			return false
		}

		queryMatcherIgnore = append(queryMatcherIgnore, "end")
	}

	if strings.Contains(r.URL.RawQuery, "start=") {
		if !strings.Contains(i.URL, "start=") {
			return false
		}

		queryMatcherIgnore = append(queryMatcherIgnore, "start")
	}

	// Return early if no further operation is needed (nothing to be ignored from the query)
	if len(queryMatcherIgnore) == 0 {
		return true
	}

	cassetteURL, _ := url.Parse(i.URL)
	requestURL := r.URL

	requestURLValues := requestURL.Query()
	cassetteURLValues := cassetteURL.Query()

	for _, query := range queryMatcherIgnore {
		requestURLValues.Del(query)
		cassetteURLValues.Del(query)
	}

	r.URL.RawQuery = requestURLValues.Encode()
	cassetteURL.RawQuery = cassetteURLValues.Encode()
	i.URL = cassetteURL.String()

	return true
}

func cassetteMatcher(r *http.Request, i cassette.Request) bool {
	// Docker
	if r.URL.Host == windowDockerEngine || r.URL.Host == "npipe://"+windowDockerEngine {
		r.URL.Host = unixDockerEngine
	}

	// Buildpacks
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
	match, err := regexp.MatchString(".s3.[a-z]{2}-[a-z]{3}.scw.cloud", r.URL.Host)
	if err == nil && match {
		return customS3Matcher(r, i)
	}

	// Temporary QueryMatcherIgnore for container tests (should be removed when #5580 is merged and the new vcr is used)
	// The cassette stores real IDs and data, but the fake client used to replay interactions only has placeholder values.
	// We need to ignore those fields in the request URL (after verifying that they are also present in the cassette).
	if !applyQueryMatcherIgnore(r, &i) {
		return false
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

// getHTTPRecoder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP request must be made to run the tests.
//
// It is important to call add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func getHTTPRecoder(t *testing.T, update bool) (client *http.Client, cleanup func(), err error) {
	t.Helper()
	recorderMode := recorder.ModeReplaying
	if update {
		recorderMode = recorder.ModeRecording
	}

	// Setup recorder and scw client
	r, err := recorder.NewAsMode(
		getTestFilePath(t, ".cassette"),
		recorderMode,
		&SocketPassthroughTransport{},
	)
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

func newHTTPRecorder(t *testing.T, folder string, update bool) (*http.Client, func(), error) {
	t.Helper()

	r, err := vcr.NewHTTPRecorder(t, folder, update, &SocketPassthroughTransport{})
	if err != nil {
		return nil, nil, err
	}

	return &http.Client{Transport: &retryableHTTPTransport{transport: r}}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
