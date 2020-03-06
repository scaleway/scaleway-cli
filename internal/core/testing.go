package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	// Environment variable are prefixed by "CLI_" in order to avoid magic behavior with SDK variables.
	// E.g.: SDK_UPDATE_CASSETTES=false will disable retry on WaitFor* methods.
	UpdateGoldens   = os.Getenv("CLI_UPDATE_GOLDENS") == "true"
	UpdateCassettes = os.Getenv("CLI_UPDATE_CASSETTES") == "true"
)

// CheckFuncCtx contain the result of a command execution
type CheckFuncCtx struct {
	// Exit code return by the CLI
	ExitCode int

	// Content print on stdout
	Stdout []byte

	// Content print on stderr
	Stderr []byte

	// Error returned by the command
	Err error

	// Command result
	Result interface{}

	// Meta bag
	Meta map[string]interface{}
}

// TestCheck is a function that perform assertion on a CheckFuncCtx
type TestCheck func(t *testing.T, ctx *CheckFuncCtx)

type BeforeFunc func(ctx *BeforeFuncCtx) error

type AfterFunc func(ctx *AfterFuncCtx) error

type BeforeFuncCtx struct {
	Client     *scw.Client
	ExecuteCmd func(cmd string) interface{}
	Meta       map[string]interface{}
}

type AfterFuncCtx struct {
	Client     *scw.Client
	ExecuteCmd func(cmd string) interface{}
	Meta       map[string]interface{}
	CmdResult  interface{}
}

// TestConfig contain configuration that can be used with the Test function
type TestConfig struct {

	// Array of command to load (see main.go)
	Commands *Commands

	// If set to true the client will be initialize to use a e2e token.
	UseE2EClient bool

	// DefaultRegion to use with scw client (default: scw.RegionFrPar)
	DefaultRegion scw.Region

	// DefaultZone to use with scw client (default: scw.ZoneFrPar1)
	DefaultZone scw.Zone

	// BeforeFunc is a hook that will be called before test is run. You can use this function to bootstrap resources.
	BeforeFunc BeforeFunc

	// The command line you want to test
	Cmd string

	// A list of check function that will be run on result.
	Check TestCheck

	// AfterFunc is a hook that will be called after test is run. You can use this function to teardown resources.
	AfterFunc AfterFunc

	// Run tests in parallel.
	DisableParallel bool

	// Fake build info for this test.
	BuildInfo BuildInfo
}

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, suffix string) string {
	fileName := t.Name()
	fileName = strings.Replace(fileName, "/", "-", -1)
	fileName = strcase.ToBashArg(fileName)
	return filepath.Join(".", "testdata", fileName+suffix)
}

func getTestClient(t *testing.T, testConfig *TestConfig) (client *scw.Client, cleanup func()) {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithEnv(),
		scw.WithUserAgent("cli-e2e-test"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
	}

	if !testConfig.UseE2EClient {
		httpClient, cleanup, err := getHTTPRecoder(t, UpdateCassettes)
		require.NoError(t, err)
		clientOpts = append(clientOpts, scw.WithHTTPClient(httpClient))
		config, err := scw.LoadConfig()
		if err == nil {
			p, err := config.GetActiveProfile()
			require.NoError(t, err)
			clientOpts = append(clientOpts, scw.WithProfile(p))
		}
		if testConfig.DefaultRegion != "" {
			clientOpts = append(clientOpts, scw.WithDefaultRegion(testConfig.DefaultRegion))
		}

		if testConfig.DefaultZone != "" {
			clientOpts = append(clientOpts, scw.WithDefaultZone(testConfig.DefaultZone))
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
		if config.DisableParallel {
			t.Parallel()
		}

		// Because human marshal of date is relative (e.g 3 minutes ago) we must make sure it stay consistent for golden to works.
		// Here we return a constant string. We may need to find a better place to put this.
		human.RegisterMarshalerFunc(time.Time{}, func(i interface{}, opt *human.MarshalOpt) (string, error) {
			return "few seconds ago", nil
		})

		client, cleanup := getTestClient(t, config)
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
			logger.Debugf("command: %s", cmdTemplate(cmd))
			_, result, err := Bootstrap(&BootstrapConfig{
				Args:             strings.Split(cmdTemplate(cmd), " "),
				Commands:         config.Commands,
				BuildInfo:        &config.BuildInfo,
				Stdout:           stdoutBuffer,
				Stderr:           stderrBuffer,
				Client:           client,
				DisableTelemetry: true,
			})
			require.NoError(t, err, "stdout: %s\nstderr: %s", stdoutBuffer.String(), stderrBuffer.String())

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
		var result interface{}
		var exitCode int
		var err error

		if config.Cmd != "" {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			logger.Debugf("command: %s", cmdTemplate(config.Cmd))
			exitCode, result, err = Bootstrap(&BootstrapConfig{
				Args:             strings.Split(cmdTemplate(config.Cmd), " "),
				Commands:         config.Commands,
				BuildInfo:        &config.BuildInfo,
				Stdout:           stdout,
				Stderr:           stderr,
				Client:           client,
				DisableTelemetry: true,
			})

			config.Check(t, &CheckFuncCtx{
				ExitCode: exitCode,
				Stdout:   stdout.Bytes(),
				Stderr:   stderr.Bytes(),
				Meta:     meta,
				Result:   result,
				Err:      err,
			})
		}

		// Run config.AfterFunc
		if config.AfterFunc != nil {
			require.NoError(t, config.AfterFunc(&AfterFuncCtx{
				Client:     client,
				ExecuteCmd: executeCmd,
				Meta:       meta,
				CmdResult:  result,
			}))
		}
	}
}

// BeforeFuncCombine combines multiple before functions into one.
func BeforeFuncCombine(beforeFuncs ...BeforeFunc) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		if len(beforeFuncs) < 2 {
			panic(fmt.Errorf("BeforeFuncCombine must be used to combine more than one BeforeFunc"))
		}
		for _, beforeFunc := range beforeFuncs {
			err := beforeFunc(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// AfterFuncCombine combines multiple after functions into one.
func AfterFuncCombine(afterFuncs ...AfterFunc) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
		if len(afterFuncs) < 2 {
			panic(fmt.Errorf("AfterFuncCombine must be used to combine more than one AfterFunc"))
		}
		for _, afterFunc := range afterFuncs {
			err := afterFunc(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// ExecStoreBeforeCmd executes the given before command and register the result
// in the context Meta at metaKey.
func ExecStoreBeforeCmd(metaKey, cmd string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.Meta[metaKey] = ctx.ExecuteCmd(cmd)
		return nil
	}
}

// ExecBeforeCmd executes the given before command.
func ExecBeforeCmd(cmd string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.ExecuteCmd(cmd)
		return nil
	}
}

// ExecAfterCmd executes the given before command.
func ExecAfterCmd(cmd string) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
		ctx.ExecuteCmd(cmd)
		return nil
	}
}

// TestCheckCombine combines multiple check functions into one.
func TestCheckCombine(checks ...TestCheck) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, true, len(checks) > 1)
		for _, check := range checks {
			check(t, ctx)
		}
	}
}

// TestCheckExitCode assert exitCode
func TestCheckExitCode(expectedCode int) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, expectedCode, ctx.ExitCode, "Invalid exit code")
	}
}

// TestCheckStderrGolden assert stderr using golden
func TestCheckStderrGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		testGolden(t, getTestFilePath(t, ".stderr.golden"), ctx.Stderr)
	}
}

// TestCheckStdoutGolden assert stdout using golden
func TestCheckStdoutGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		testGolden(t, getTestFilePath(t, ".stdout.golden"), ctx.Stdout)
	}
}

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return TestCheckCombine(
		TestCheckStdoutGolden(),
		TestCheckStderrGolden(),
	)
}

// TestCheckError asserts error
func TestCheckError(err error) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, err, ctx.Err, "Invalid error")
	}
}

func testGolden(t *testing.T, goldenPath string, actual []byte) {
	actualIsEmpty := len(actual) == 0

	// In order to avoid diff in goldens we set all timestamp to the same date
	actual = uniformLogTimestamps(actual)
	if UpdateGoldens {
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
		assert.Equal(t, string(expected), string(actual))
	}
}

var regLogTimestamp = regexp.MustCompile(`((\d)+\/(\d)+\/(\d)+ (\d)+\:(\d)+\:(\d)+)`)

// uniformLogTimestamps replace all log timestamp to the date "2019/12/09 16:04:07"
func uniformLogTimestamps(input []byte) []byte {
	return regLogTimestamp.ReplaceAll(input, []byte("2019/12/09 16:04:07"))
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
