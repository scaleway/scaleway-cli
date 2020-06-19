package core

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/dnaeon/go-vcr/cassette"
	"github.com/dnaeon/go-vcr/recorder"
	"github.com/scaleway/scaleway-cli/internal/account"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Environment variable are prefixed by "CLI_" in order to avoid magic behavior with SDK variables.
// E.g.: SDK_UPDATE_CASSETTES=false will disable retry on WaitFor* methods.
var (
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
	Meta TestMeta

	// Scaleway client
	Client *scw.Client

	// OverrideEnv passed in the TestConfig
	OverrideEnv map[string]string

	Logger *Logger
}

// TestMeta contains arbitrary data that can be passed along a test lifecycle.
type TestMeta map[string]interface{}

// Tpl render a go template using where content of meta can be used
func (meta TestMeta) Tpl(strTpl string) string {
	t := meta["t"].(*testing.T)
	buf := &bytes.Buffer{}
	require.NoError(t, template.Must(template.New("tpl").Parse(strTpl)).Execute(buf, meta))
	return buf.String()
}

// TestCheck is a function that perform assertion on a CheckFuncCtx
type TestCheck func(t *testing.T, ctx *CheckFuncCtx)

type BeforeFunc func(ctx *BeforeFuncCtx) error

type AfterFunc func(ctx *AfterFuncCtx) error

type ExecFuncCtx struct {
	T      *testing.T
	Meta   TestMeta
	Client *scw.Client
}

type OverrideExecTestFunc func(ctx *ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error)

type BeforeFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) interface{}
	Meta        TestMeta
	OverrideEnv map[string]string
	Logger      *Logger
}

type AfterFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) interface{}
	Meta        TestMeta
	CmdResult   interface{}
	OverrideEnv map[string]string
	Logger      *Logger
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
	// Conflict with Args
	Cmd string

	// Args represents a program arguments and should be used, when you cannot Cmd because your arguments include space characters
	// Conflict with Cmd
	Args []string

	// A list of check function that will be run on result.
	Check TestCheck

	// AfterFunc is a hook that will be called after test is run. You can use this function to teardown resources.
	AfterFunc AfterFunc

	// Run tests in parallel.
	DisableParallel bool

	// Fake build info for this test.
	BuildInfo BuildInfo

	// If set, it will create a temporary home directory during the tests.
	// Get this folder with ExtractUserHomeDir()
	TmpHomeDir bool

	// OverrideEnv contains environment variables that will be overridden during the test.
	OverrideEnv map[string]string

	// see BootstrapConfig.OverrideExec
	OverrideExec OverrideExecTestFunc

	// Custom client to use for test, if none are provided will create one automatically
	Client *scw.Client

	// Context that will be forwarded to Bootstrap
	Ctx context.Context

	// If this is specified this value will be passed to interactive.InjectMockResponseToContext ans will allow
	// to mock response a user would have enter in a prompt.
	// Warning: All prompts MUST be mocked or test will hang.
	PromptResponseMocks []string

	// Allow to mock stdin
	Stdin io.Reader
}

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, suffix string) string {
	specialChars := regexp.MustCompile(`[\\?%*:|"<>. ]`)

	// Replace nested tests separators.
	fileName := strings.Replace(t.Name(), "/", "-", -1)

	fileName = strcase.ToBashArg(fileName)

	// Replace special characters.
	fileName = specialChars.ReplaceAllLiteralString(fileName, "") + suffix

	return filepath.Join(".", "testdata", fileName)
}

func createTestClient(t *testing.T, testConfig *TestConfig, httpClient *http.Client) (client *scw.Client) {
	var err error

	// Init default options
	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithUserAgent("cli-e2e-test"),
	}

	// If client is NOT an E2E client we init http recorder and load configuration.
	if !testConfig.UseE2EClient {
		clientOpts = append(clientOpts, scw.WithHTTPClient(httpClient))

		if UpdateCassettes {
			clientOpts = append(clientOpts, scw.WithEnv())
			config, err := scw.LoadConfig()
			if err == nil {
				p, err := config.GetActiveProfile()
				require.NoError(t, err)
				clientOpts = append(clientOpts, scw.WithProfile(p))
			}
		}
	}

	// We handle default zone and region configured specifically for a test
	if testConfig.DefaultRegion != "" {
		clientOpts = append(clientOpts, scw.WithDefaultRegion(testConfig.DefaultRegion))
	}
	if testConfig.DefaultZone != "" {
		clientOpts = append(clientOpts, scw.WithDefaultZone(testConfig.DefaultZone))
	}

	client, err = scw.NewClient(clientOpts...)
	require.NoError(t, err)

	// If client is an E2E client we must register and use returned credential.
	if testConfig.UseE2EClient {
		res, err := test.NewAPI(client).Register(&test.RegisterRequest{Username: "sidi"})
		require.NoError(t, err)

		client, err = scw.NewClient(append(clientOpts, scw.WithAuth(res.AccessKey, res.SecretKey))...)
		require.NoError(t, err)
	}

	return client
}

// DefaultRetryInterval is used across all wait functions in the CLI
// In particular it is very handy to define this RetryInterval at 0 second while running cassette in testing
// because they will be executed without waiting.
var DefaultRetryInterval *time.Duration

// Run a CLI integration test. See TestConfig for configuration option
func Test(config *TestConfig) func(t *testing.T) {
	return func(t *testing.T) {
		if !config.DisableParallel {
			t.Parallel()
		}

		log := &Logger{
			writer: os.Stderr,
			level:  logger.LogLevelInfo,
		}
		if os.Getenv("SCW_DEBUG") == "true" {
			log.level = logger.LogLevelDebug
		}

		// Because human marshal of date is relative (e.g 3 minutes ago) we must make sure it stay consistent for golden to works.
		// Here we return a constant string. We may need to find a better place to put this.
		human.RegisterMarshalerFunc(time.Time{}, func(i interface{}, opt *human.MarshalOpt) (string, error) {
			return "few seconds ago", nil
		})

		if !UpdateCassettes {
			tmp := 0 * time.Second
			DefaultRetryInterval = &tmp
		}

		ctx := config.Ctx
		if ctx == nil {
			ctx = context.Background()
		}
		if len(config.PromptResponseMocks) > 0 {
			ctx = interactive.InjectMockResponseToContext(ctx, config.PromptResponseMocks)
		}

		httpClient, cleanup, err := getHTTPRecoder(t, UpdateCassettes)
		require.NoError(t, err)
		defer cleanup()
		ctx = account.InjectHTTPClient(ctx, httpClient)

		// We try to use the client provided in the config
		// if no client is provided in the config we create a test client
		client := config.Client
		if client == nil {
			client = createTestClient(t, config, httpClient)
		}

		meta := TestMeta{
			"t": t,
		}

		overrideEnv := config.OverrideEnv
		if overrideEnv == nil {
			overrideEnv = map[string]string{}
		}

		if config.TmpHomeDir {
			dir, err := ioutil.TempDir(os.TempDir(), "scw")
			require.NoError(t, err)
			defer func() {
				err = os.RemoveAll(dir)
				assert.NoError(t, err)
			}()
			overrideEnv["HOME"] = dir
			meta["HOME"] = dir
		}

		overrideExec := defaultOverrideExec
		if config.OverrideExec != nil {
			overrideExec = func(cmd *exec.Cmd) (exitCode int, err error) {
				return config.OverrideExec(&ExecFuncCtx{
					T:      t,
					Meta:   meta,
					Client: client,
				}, cmd)
			}
		}

		stdin := config.Stdin
		if stdin == nil {
			stdin = os.Stdin
		}

		executeCmd := func(args []string) interface{} {
			stdoutBuffer := &bytes.Buffer{}
			stderrBuffer := &bytes.Buffer{}
			log.Debugf("command: %s", args)
			_, result, err := Bootstrap(&BootstrapConfig{
				Args:             args,
				Commands:         config.Commands,
				BuildInfo:        &config.BuildInfo,
				Stdout:           stdoutBuffer,
				Stderr:           stderrBuffer,
				Client:           client,
				DisableTelemetry: true,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           log,
			})
			require.NoError(t, err, "stdout: %s\nstderr: %s", stdoutBuffer.String(), stderrBuffer.String())

			return result
		}

		// Run config.BeforeFunc
		if config.BeforeFunc != nil {
			log.Debug("Start BeforeFunc")
			require.NoError(t, config.BeforeFunc(&BeforeFuncCtx{
				T:           t,
				Client:      client,
				ExecuteCmd:  executeCmd,
				Meta:        meta,
				OverrideEnv: overrideEnv,
				Logger:      log,
			}))
			log.Debug("End BeforeFunc")
		}

		// Run config.Cmd
		var result interface{}
		var exitCode int
		args := config.Args
		if config.Cmd != "" {
			args = cmdToArgs(meta, config.Cmd)
		}
		if len(args) > 0 {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			log.Debugf("executing command: %s\n", args)
			exitCode, result, err = Bootstrap(&BootstrapConfig{
				Args:             args,
				Commands:         config.Commands,
				BuildInfo:        &config.BuildInfo,
				Stdout:           stdout,
				Stderr:           stderr,
				Stdin:            stdin,
				Client:           client,
				DisableTelemetry: true,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           log,
			})

			log.Debugf("Store in meta (with key CmdResult): %s\n", result)
			meta["CmdResult"] = result
			config.Check(t, &CheckFuncCtx{
				ExitCode:    exitCode,
				Stdout:      stdout.Bytes(),
				Stderr:      stderr.Bytes(),
				Meta:        meta,
				Result:      result,
				Err:         err,
				Client:      client,
				OverrideEnv: overrideEnv,
				Logger:      log,
			})
		}

		// Run config.AfterFunc
		if config.AfterFunc != nil {
			log.Debug("Start AfterFunc")
			require.NoError(t, config.AfterFunc(&AfterFuncCtx{
				T:           t,
				Client:      client,
				ExecuteCmd:  executeCmd,
				Meta:        meta,
				CmdResult:   result,
				OverrideEnv: overrideEnv,
				Logger:      log,
			}))
			log.Debug("End AfterFunc")
		}
	}
}

func cmdToArgs(meta TestMeta, s string) []string {
	return strings.Split(meta.Tpl(s), " ")
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

func BeforeFuncWhenUpdatingCassette(beforeFunc BeforeFunc) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		if UpdateCassettes {
			return beforeFunc(ctx)
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
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecStoreBeforeCmd (in metaKey %s): %s\n", metaKey, args)
		ctx.Meta[metaKey] = ctx.ExecuteCmd(args)
		return nil
	}
}

func BeforeFuncOsExec(cmd string, args ...string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.Logger.Debugf("BeforeFuncOsExec: %s %s\n", cmd, args)
		return exec.Command(cmd, args...).Run()
	}
}

// ExecBeforeCmd executes the given before command.
func ExecBeforeCmd(cmd string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecBeforeCmd: %s\n", args)
		ctx.ExecuteCmd(args)
		return nil
	}
}

// ExecAfterCmd executes the given before command.
func ExecAfterCmd(cmd string) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecAfterCmd: %s\n", args)
		ctx.ExecuteCmd(args)
		return nil
	}
}

// TestCheckCombine combines multiple check functions into one.
func TestCheckCombine(checks ...TestCheck) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, true, len(checks) > 1, "TestCheckCombine must be used to combine more than one TestCheck")
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

// testCheckStderrGolden assert stderr using golden
func testCheckStderrGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		testGolden(t, getTestFilePath(t, ".stderr.golden"), ctx.Stderr)
	}
}

// testCheckStdoutGolden assert stdout using golden
func testCheckStdoutGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		testGolden(t, getTestFilePath(t, ".stdout.golden"), ctx.Stdout)
	}
}

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return TestCheckCombine(
		testCheckStdoutGolden(),
		testCheckStderrGolden(),
	)
}

// TestCheckError asserts error
func TestCheckError(err error) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, err, ctx.Err, "Invalid error")
	}
}

// TestCheckStdout asserts stdout using string
func TestCheckStdout(stdout string) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		assert.Equal(t, stdout, string(ctx.Stdout), "Invalid stdout")
	}
}

func OverrideExecSimple(cmdStr string, exitCode int) OverrideExecTestFunc {
	return func(ctx *ExecFuncCtx, cmd *exec.Cmd) (int, error) {
		assert.Equal(ctx.T, ctx.Meta.Tpl(cmdStr), strings.Join(cmd.Args, " "))
		return exitCode, nil
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
			require.NoError(t, ioutil.WriteFile(goldenPath, actual, 0644)) //nolint:gosec
		}
	}

	expected, err := ioutil.ReadFile(goldenPath)
	if actualIsEmpty {
		assert.NotNil(t, err)
	} else {
		require.NoError(t, err, "expected to find golden file with %s", string(actual))

		// Replace Windows return carriage.
		expected = bytes.ReplaceAll(expected, []byte("\r"), []byte(""))
		actual = bytes.ReplaceAll(actual, []byte("\r"), []byte(""))

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
		i.Request.URL = regexp.MustCompile("organization_id=[0-9a-f-]{36}").ReplaceAllString(i.Request.URL, "organization_id=11111111-1111-1111-1111-111111111111")
		return nil
	})

	return &http.Client{Transport: r}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}
