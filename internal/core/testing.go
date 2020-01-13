package core

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var updateGolden = os.Getenv("UPDATE_GOLDEN") == "true"

// TestResult contain the result of a command execution
type TestResult struct {
	// Exit code return by the CLI
	ExitCode int

	// Content print on stdout
	Stdout []byte

	// Content print on stderr
	Stderr []byte
}

// TestCheck is a function that perform assertion on a TestResult
type TestCheck func(*testing.T, *TestResult)

type BeforeFuncCtx struct {
	Client     *scw.Client
	ExecuteCmd func(cmd string) interface{}
	Meta       map[string]interface{}
}

type AfterFuncCtx struct {
	Client     *scw.Client
	ExecuteCmd func(cmd string) interface{}
	Meta       map[string]interface{}
}

// TestConfig contain configuration that can be used with the Test function
type TestConfig struct {

	// Array of command to load (see main.go)
	Commands *Commands

	// If set to true the client will be initialize to use a e2e token.
	UseE2EClient bool

	// Hook that will be called before test is run. You can use this function to bootstrap resources.
	BeforeFunc func(ctx *BeforeFuncCtx) error

	// The command line you want to test
	Cmd string

	//  Hook that will be called after test is run. You can use this function to teardown resources.
	AfterFunc func(ctx *AfterFuncCtx) error

	// A list of check function that will be run on result
	Check TestCheck
}

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, suffix string) string {
	fileName := t.Name()
	fileName = strings.Replace(fileName, "/", "-", -1)
	fileName = strcase.ToBashArg(fileName)
	return filepath.Join(".", "testdata", fileName+suffix)
}

func getTestClient(t *testing.T, e2eClient bool) (client *scw.Client, cleanup func()) {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithEnv(),
		scw.WithUserAgent("cli-e2e-test"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
	}

	if !e2eClient {
		httpClient, cleanup, err := getHttpRecoder(t, updateGolden)
		require.NoError(t, err)
		clientOpts = append(clientOpts, scw.WithHTTPClient(httpClient))
		config, err := scw.LoadConfig()
		if err == nil {
			p, err := config.GetActiveProfile()
			require.NoError(t, err)
			clientOpts = append(clientOpts, scw.WithProfile(p))
		}

		client, err := scw.NewClient(clientOpts...)
		require.NoError(t, err)
		return client, cleanup
	}

	client, err := scw.NewClient(clientOpts...)
	require.NoError(t, err)
	res, err := test.NewAPI(client).Register(&test.RegisterRequest{Username: "sidi"})
	require.NoError(t, err)

	client, err = scw.NewClient(append(clientOpts, scw.WithAuth(res.AccessKey, res.SecretKey))...)
	require.NoError(t, err)

	return client, func() {}
}

// Run a CLI integration test. See TestConfig for configuration option
func Test(config *TestConfig) func(t *testing.T) {
	return func(t *testing.T) {

		// Because human marshal of date is relative (e.g 3 minutes ago) we must make sure it stay consistent for golden to works.
		// Here we return a constant string. We may need to find a better place to put this.
		human.RegisterMarshalerFunc(time.Time{}, func(i interface{}, opt *human.MarshalOpt) (string, error) {
			return "few seconds ago", nil
		})

		client, cleanup := getTestClient(t, config.UseE2EClient)
		defer cleanup()

		meta := map[string]interface{}{}

		cmdTemplate := func(cmd string) string {
			cmdBuf := &bytes.Buffer{}
			require.NoError(t, template.Must(template.New("cmd").Parse(cmd)).Execute(cmdBuf, meta))
			return cmdBuf.String()
		}

		executeCmd := func(cmd string) interface{} {
			stdoutBuffer := &bytes.Buffer{}
			stderrBuffer := &bytes.Buffer{}

			exitCode := Bootstrap(&BootstrapConfig{
				Args:      append(strings.Split(cmdTemplate(cmd), " "), "-o", "json"),
				Commands:  config.Commands,
				BuildInfo: &BuildInfo{},
				Stdout:    stdoutBuffer,
				Stderr:    stderrBuffer,
				Client:    client,
			})
			require.Equal(t, 0, exitCode, "stdout: %s\nstderr: %s", stdoutBuffer.String(), stderrBuffer.String())

			result := map[string]interface{}{}
			require.NoError(t, json.Unmarshal(stdoutBuffer.Bytes(), &result))

			return result
		}

		// Run config.BeforeFunc

		if config.BeforeFunc != nil {
			require.NoError(t, config.BeforeFunc(&BeforeFuncCtx{
				Client:     client,
				ExecuteCmd: executeCmd,
				Meta:       meta,
			}))
		}

		// Run config.Cmd

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		exitCode := Bootstrap(&BootstrapConfig{
			Args:      strings.Split(cmdTemplate(config.Cmd), " "),
			Commands:  config.Commands,
			BuildInfo: &BuildInfo{},
			Stdout:    stdout,
			Stderr:    stderr,
			Client:    client,
		})

		result := &TestResult{
			ExitCode: exitCode,
			Stdout:   stdout.Bytes(),
			Stderr:   stderr.Bytes(),
		}

		config.Check(t, result)

		// Run config.AfterFunc

		if config.AfterFunc != nil {
			require.NoError(t, config.AfterFunc(&AfterFuncCtx{
				Client:     client,
				ExecuteCmd: executeCmd,
				Meta:       meta,
			}))
		}
	}
}

// TestCheckCombine Combine multiple check function into one
func TestCheckCombine(checks ...TestCheck) TestCheck {
	return func(t *testing.T, result *TestResult) {
		for _, check := range checks {
			check(t, result)
		}
	}
}

// TestCheckExitCode assert exitCode
func TestCheckExitCode(expectedCode int) TestCheck {
	return func(t *testing.T, result *TestResult) {
		assert.Equal(t, expectedCode, result.ExitCode, "Invalid exit code")
	}
}

// TestCheckStderrGolden assert stderr using golden
func TestCheckStderrGolden() TestCheck {
	return func(t *testing.T, result *TestResult) {
		testGolden(t, getTestFilePath(t, ".stderr.golden"), result.Stderr)
	}
}

// TestCheckStdoutGolden assert stdout using golden
func TestCheckStdoutGolden() TestCheck {
	return func(t *testing.T, result *TestResult) {
		testGolden(t, getTestFilePath(t, ".stdout.golden"), result.Stdout)
	}
}

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return TestCheckCombine(
		TestCheckStdoutGolden(),
		TestCheckStderrGolden(),
	)
}

func testGolden(t *testing.T, goldenPath string, actual []byte) {
	actualIsEmpty := len(actual) == 0

	if updateGolden {

		if actualIsEmpty {
			_ = os.Remove(goldenPath)
		} else {
			require.NoError(t, os.MkdirAll(path.Dir(goldenPath), 0755))
			require.NoError(t, ioutil.WriteFile(goldenPath, actual, 0644))
		}
	}

	expected, err := ioutil.ReadFile(goldenPath)
	if actualIsEmpty {
		assert.NotNil(t, err)
	} else {
		require.NoError(t, err)
		assert.Equal(t, string(actual), string(expected))
	}

}

// getHttpRecoder creates a new httpClient that records all HTTP requests in a cassette.
// This cassette is then replayed whenever tests are executed again. This means that once the
// requests are recorded in the cassette, no more real HTTP request must be made to run the tests.
//
// It is important to call add a `defer cleanup()` so the given cassette files are correctly
// closed and saved after the requests.
func getHttpRecoder(t *testing.T, update bool) (client *http.Client, cleanup func(), err error) {

	recorderMode := recorder.ModeReplaying
	if update {
		recorderMode = recorder.ModeRecording
	}

	// Setup recorder and scw client
	r, err := recorder.NewAsMode(getTestFilePath(t, ".cassette"), recorderMode, nil)
	if err != nil {
		return nil, nil, err
	}

	// Add a filter which removes Authorization headers from all requests:
	r.AddFilter(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "x-auth-token")
		delete(i.Request.Headers, "X-Auth-Token")
		return nil
	})

	return &http.Client{Transport: r}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
