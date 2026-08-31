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

	return nil
}

func cassetteMatcher(r *http.Request, i cassette.Request) bool {
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
	// or path-style: https://s3.fr-par.scw.cloud/test-acc-scaleway-object-bucket-lifecycle-8445817190507446251?lifecycle=
	match, err := regexp.MatchString(`s3\.[a-z]{2}-[a-z]{3}\.scw\.cloud`, r.URL.Host)
	if err == nil && match {
		return customS3Matcher(r, i)
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

	// Extract the bucket name and remaining path, handling both
	// virtual-hosted style (bucket in host) and path style (bucket in path).
	actualBucket, actualPath := extractS3BucketAndPath(actualS3Host, actualURL.Path)
	expectedBucket, expectedPath := extractS3BucketAndPath(expectedS3Host, expectedURL.Path)

	// When no bucket is present (e.g. ListBuckets), fall back to the
	// default matcher which compares the full URL.
	if actualBucket == "" && expectedBucket == "" {
		return cassette.DefaultMatcher(r, i)
	}
	if actualBucket == "" || expectedBucket == "" {
		return false
	}

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

	return r.Method == i.Method && actualPath == expectedPath &&
		actualURL.RawQuery == expectedURL.RawQuery
}

// extractS3BucketAndPath extracts the bucket name and remaining path from
// either a virtual-hosted style URL (bucket in host, e.g.
// "bucket.s3.fr-par.scw.cloud") or a path style URL (bucket in path, e.g.
// "s3.fr-par.scw.cloud/bucket").
func extractS3BucketAndPath(hostParts []string, path string) (bucket, remainingPath string) {
	if len(hostParts) > 4 {
		// Virtual-hosted style: bucket is the first label of the host
		return hostParts[0], path
	}
	// Path style: bucket is the first segment of the path
	trimmed := strings.TrimPrefix(path, "/")
	if idx := strings.Index(trimmed, "/"); idx >= 0 {
		return trimmed[:idx], trimmed[idx:]
	}

	return trimmed, ""
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
