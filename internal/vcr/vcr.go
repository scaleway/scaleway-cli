package vcr

import (
	"net/http"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-sdk-go/strcase"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/cassette"
	"gopkg.in/dnaeon/go-vcr.v4/pkg/recorder"
)

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, pkgFolder string, suffix string) string {
	t.Helper()

	specialChars := regexp.MustCompile(`[\\?%*:|"'<>. ]`)

	// Replace nested tests separators.
	fileName := strings.ReplaceAll(t.Name(), "/", "-")

	fileName = strcase.ToBashArg(fileName)

	// Replace special characters.
	fileName = specialChars.ReplaceAllLiteralString(fileName, "") + suffix

	// Remove prefix to simplify
	fileName = strings.TrimPrefix(fileName, "test-acc-")

	return filepath.Join(pkgFolder, "testdata", fileName)
}

func cassetteRequestFilter(i *cassette.Interaction) error {
	// Remove authorization headers
	delete(i.Request.Headers, "x-auth-token")
	delete(i.Request.Headers, "X-Auth-Token")
	delete(i.Request.Headers, "Authorization")
	delete(i.Request.Headers, "X-Registry-Auth")

	// URL substitutions
	i.Request.URL = regexp.MustCompile(`(.+)organization_id=[0-9a-f-]{36}(.+)`).
		ReplaceAllString(i.Request.URL, "${1}organization_id=11111111-1111-1111-1111-111111111111${2}")
	i.Request.URL = regexp.MustCompile(`^https://api\.scaleway\.com/account/v1/tokens/[0-9a-f-]{36}$`).
		ReplaceAllString(i.Request.URL, "api.scaleway.com/account/v1/tokens/11111111-1111-1111-1111-111111111111")
	i.Request.URL = regexp.MustCompile(`(.+)?SCW[0-9A-Z]{17}(.+)?`).
		ReplaceAllString(i.Request.URL, "${1}SCWXXXXXXXXXXXXXXXXX${2}")

	// Buildpacks
	i.Request.URL = regexp.MustCompile(`pack\.local%2Fbuilder%2F[0-9a-f]{20}`).
		ReplaceAllString(i.Request.URL, "pack.local%2Fbuilder%2F11111111111111111111")
	i.Request.URL = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Request.URL, "pack.local/builder/11111111111111111111")
	i.Request.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Request.Body, "pack.local/builder/11111111111111111111")

	for _, vs := range i.Request.Form {
		for idx, v := range vs {
			vs[idx] = regexp.MustCompile(`(.+)organization_id=[0-9a-f-]{36}(.+)`).
				ReplaceAllString(v, "${1}organization_id=11111111-1111-1111-1111-111111111111${2}")
			vs[idx] = regexp.MustCompile(`(.+)?SCW[0-9A-Z]{17}(.+)?`).
				ReplaceAllString(v, "${1}SCWXXXXXXXXXXXXXXXXX${2}")
		}
	}

	return nil
}

func cassetteResponseFilter(i *cassette.Interaction) error {
	delete(i.Response.Headers, "Content-Security-Policy")
	delete(i.Response.Headers, "Strict-Transport-Security")
	delete(i.Response.Headers, "X-Content-Type-Options")
	delete(i.Response.Headers, "X-Frame-Options")

	i.Response.Body = regexp.MustCompile(`"secret_key":"[0-9a-f-]{36}"`).
		ReplaceAllString(i.Response.Body, `"secret_key":"11111111-1111-1111-1111-111111111111"`)

	// Buildpacks
	i.Response.Body = regexp.MustCompile(`pack\.local/builder/[0-9a-f]{20}`).
		ReplaceAllString(i.Response.Body, "pack.local/builder/11111111111111111111")

	// Some interactions like the instance snapshot download have API keys in their responses.
	i.Response.Body = regexp.MustCompile(`(.+)?SCW[0-9A-Z]{17}(.+)?`).
		ReplaceAllString(i.Response.Body, "${1}SCWXXXXXXXXXXXXXXXXX${2}")

	return nil
}

type AdditionalHook struct {
	HookFunc recorder.HookFunc
	Kind     recorder.HookKind
}

// NewHTTPRecorder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP requests must be made to run the tests.
//
// It is important to add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func NewHTTPRecorder(
	t *testing.T,
	pkgFolder string,
	update bool,
	realTransport http.RoundTripper,
	additionalHooks ...AdditionalHook,
) (*recorder.Recorder, error) {
	t.Helper()

	recorderOptions := []recorder.Option{
		// Ignore requests durations when replaying cassettes
		recorder.WithSkipRequestLatency(true),
		// Starting with v3, go-vcr now calls net/url.Parse to build the interaction which results in an error for escaped
		// Docker URLs on Test_Deploy (container), so we need to unescape these paths on this step.
		recorder.WithHook(unescapeDockerURL, recorder.AfterCaptureHook),
		// Remove information that is either sensitive or that would disrupt matching from the request
		recorder.WithHook(cassetteRequestFilter, recorder.BeforeSaveHook),
		// Remove secrets and unnecessary information from the response
		recorder.WithHook(cassetteResponseFilter, recorder.BeforeSaveHook),
		// Add custom matcher for requests and cassettes
		recorder.WithMatcher(CassetteMatcher),
	}

	// Set up recorder mode
	recorderMode := recorder.ModeReplayOnly
	if update {
		recorderMode = recorder.ModeRecordOnly
	}
	recorderOptions = append(recorderOptions, recorder.WithMode(recorderMode))

	// Set up custom transport (for the CLI only)
	if realTransport != nil {
		recorderOptions = append(recorderOptions, recorder.WithRealTransport(realTransport))
	}

	// Set up specific hooks for each DevTool
	for _, hook := range additionalHooks {
		recorderOptions = append(recorderOptions, recorder.WithHook(hook.HookFunc, hook.Kind))
	}

	// Setup recorder and scw client
	cassetteFilePath := getTestFilePath(t, pkgFolder, ".cassette")
	r, err := recorder.New(cassetteFilePath, recorderOptions...)
	if err != nil {
		return nil, err
	}
	defer func(r *recorder.Recorder) {
		_ = r.Stop()
	}(r)

	return r, nil
}
