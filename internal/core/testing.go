package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
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
	ExecuteCmd func(cmd string)
}

type AfterFuncCtx struct {
	Client     *scw.Client
	ExecuteCmd func(cmd string)
}

// TestConfig contain configuration that can be used with the Test function
type TestConfig struct {

	// Array of command to load (see main.go)
	Commands *Commands

	// The command line you want to test
	Cmd string

	// If set to true the client will be initialize to use a e2e token.
	UseE2EClient bool

	// Hook that will be called before and after test is run. You can use this function to bootstrap and teardown resources.
	BeforeFunc func(ctx *BeforeFuncCtx) error
	AfterFunc  func(ctx *AfterFuncCtx) error

	// A list of check function that will be run on result
	Check TestCheck
}

func panicIfErr(err error) {
	if err != nil {
		panic(err)
	}
}

// Return a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
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
		panicIfErr(err)
		clientOpts = append(clientOpts, scw.WithHTTPClient(httpClient))
		client, err := scw.NewClient(clientOpts...)
		panicIfErr(err)
		return client, cleanup
	}

	client, err := scw.NewClient(clientOpts...)
	panicIfErr(err)
	res, err := test.NewAPI(client).Register(&test.RegisterRequest{Username: "sidi"})
	panicIfErr(err)

	client, err = scw.NewClient(append(clientOpts, scw.WithAuth(res.AccessKey, res.SecretKey))...)
	panicIfErr(err)

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

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		client, cleanup := getTestClient(t, config.UseE2EClient)
		defer cleanup()

		executeCmd := func(cmd string) {
			exitCode := Bootstrap(&BootstrapConfig{
				Args:      strings.Split(cmd, " "),
				Commands:  config.Commands,
				BuildInfo: &BuildInfo{},
				Stdout:    &bytes.Buffer{},
				Stderr:    &bytes.Buffer{},
				Client:    client,
			})
			if exitCode != 0 {
				panic(fmt.Errorf("invalid exit code %d for command %s", exitCode, cmd))
			}
		}

		if config.BeforeFunc != nil {
			err := config.BeforeFunc(&BeforeFuncCtx{
				Client:     client,
				ExecuteCmd: executeCmd,
			})
			panicIfErr(err)
		}

		exitCode := Bootstrap(&BootstrapConfig{
			Args:      strings.Split(config.Cmd, " "),
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

		if config.AfterFunc != nil {
			err := config.AfterFunc(&AfterFuncCtx{
				Client:     client,
				ExecuteCmd: executeCmd,
			})
			panicIfErr(err)
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
			err := os.MkdirAll(path.Dir(goldenPath), 0755)
			panicIfErr(err)
			err = ioutil.WriteFile(goldenPath, actual, 0644)
			panicIfErr(err)
		}
	}

	expected, err := ioutil.ReadFile(goldenPath)
	if actualIsEmpty {
		assert.NotNil(t, err)
	} else {
		panicIfErr(err)
		assert.Equal(t, string(actual), string(expected))
	}

}

// CreateRecordedScwClient creates a new httpClient that records all HTTP requests in a cassette.
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
