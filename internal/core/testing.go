package core

import (
	"bytes"
	"context"
	"flag"
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
	"github.com/hashicorp/go-version"
	args "github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/interactive"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test flags
// You can create a binary of each test using "go test -c -o myBinary"
var (
	// UpdateGoldens will update all the golden files of a given test
	UpdateGoldens = flag.Bool("goldens", os.Getenv("CLI_UPDATE_GOLDENS") == "true", "Record goldens")

	// UpdateCassettes will update all cassettes of a given test
	UpdateCassettes = flag.Bool("cassettes", os.Getenv("CLI_UPDATE_CASSETTES") == "true", "Record Cassettes")

	// Debug set the log level to LogLevelDebug
	Debug = flag.Bool("debug", os.Getenv("SCW_DEBUG") == "true", "Enable Debug Mode")
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
	Meta testMetadata

	// Scaleway client
	Client *scw.Client

	// OverrideEnv passed in the TestConfig
	OverrideEnv map[string]string

	Logger *Logger

	// The content logged by the command
	LogBuffer string
}

// testMetadata contains arbitrary data that can be passed along a test lifecycle.
type testMetadata map[string]interface{}

// render renders a go template using where content of meta can be used
func (meta testMetadata) render(strTpl string) string {
	t := meta["t"].(*testing.T)
	buf := &bytes.Buffer{}
	require.NoError(t, template.Must(template.New("tpl").Parse(strTpl)).Execute(buf, meta))
	return buf.String()
}

func BeforeFuncStoreInMeta(key string, value interface{}) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.Meta[key] = value
		return nil
	}
}

// TestCheck is a function that perform assertion on a CheckFuncCtx
type TestCheck func(t *testing.T, ctx *CheckFuncCtx)

type BeforeFunc func(ctx *BeforeFuncCtx) error

type AfterFunc func(ctx *AfterFuncCtx) error

type ExecFuncCtx struct {
	T      *testing.T
	Meta   testMetadata
	Client *scw.Client
}

type OverrideExecTestFunc func(ctx *ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error)

type BeforeFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) interface{}
	Meta        testMetadata
	OverrideEnv map[string]string
	Logger      *Logger
}

type AfterFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) interface{}
	Meta        testMetadata
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
	// The arguments in this command MUST have only one space between each others to be split successfully
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
	BuildInfo *BuildInfo

	// If set, it will create a temporary home directory during the tests.
	// Get this folder with ExtractUserHomeDir()
	// This will also use this temporary directory as a cache directory.
	// Get this folder with ExtractCacheDir()
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
		scw.WithDefaultProjectID("11111111-1111-1111-1111-111111111111"),
		scw.WithUserAgent("cli-e2e-test"),
		scw.WithHTTPClient(&http.Client{
			Transport: &retryableHTTPTransport{transport: http.DefaultTransport},
		}),
	}

	// If client is NOT an E2E client we init http recorder and load configuration.
	if !testConfig.UseE2EClient {
		clientOpts = append(clientOpts, scw.WithHTTPClient(httpClient))

		if *UpdateCassettes {
			clientOpts = append(clientOpts, scw.WithEnv())
			config, err := scw.LoadConfig()
			if err == nil {
				activeProfile, err := config.GetActiveProfile()
				require.NoError(t, err)
				envProfile := scw.LoadEnvProfile()
				profile := scw.MergeProfiles(activeProfile, envProfile)
				clientOpts = append(clientOpts, scw.WithProfile(profile))
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

		testLogger := &Logger{
			writer: os.Stderr,
			level:  logger.LogLevelInfo,
		}

		if *Debug {
			testLogger.level = logger.LogLevelDebug
		}

		// We need to set up this variable to ensure that relative date parsing stay consistent
		args.TestForceNow = scw.TimePtr(time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC))

		// Because human marshal of date is relative (e.g 3 minutes ago) we must make sure it stay consistent for golden to works.
		// Here we return a constant string. We may need to find a better place to put this.
		human.RegisterMarshalerFunc(time.Time{}, func(i interface{}, opt *human.MarshalOpt) (string, error) {
			return "few seconds ago", nil
		})

		if !*UpdateCassettes {
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

		httpClient, cleanup, err := getHTTPRecoder(t, *UpdateCassettes)
		require.NoError(t, err)
		defer cleanup()

		// We try to use the client provided in the config
		// if no client is provided in the config we create a test client
		client := config.Client
		if client == nil {
			client = createTestClient(t, config, httpClient)
		}

		meta := testMetadata{
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
			overrideEnv[scw.ScwCacheDirEnv] = dir
			meta["HOME"] = dir
			meta[scw.ScwCacheDirEnv] = dir
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

		buildInfo := config.BuildInfo
		if buildInfo == nil {
			buildInfo = &BuildInfo{
				Version:   version.Must(version.NewSemver("v0.0.0+test")),
				BuildDate: "unknown",
				GoVersion: "runtime.Version()",
				GitBranch: "unknown",
				GitCommit: "unknown",
				GoArch:    "runtime.GOARCH",
				GoOS:      "runtime.GOOS",
			}
		}

		executeCmd := func(args []string) interface{} {
			stdoutBuffer := &bytes.Buffer{}
			stderrBuffer := &bytes.Buffer{}
			_, result, err := Bootstrap(&BootstrapConfig{
				Args:             args,
				Commands:         config.Commands,
				BuildInfo:        buildInfo,
				Stdout:           stdoutBuffer,
				Stderr:           stderrBuffer,
				Client:           client,
				DisableTelemetry: true,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           testLogger,
				HTTPClient:       httpClient,
			})
			require.NoError(t, err, "error executing cmd (%s)\nstdout: %s\nstderr: %s", args, stdoutBuffer.String(), stderrBuffer.String())

			return result
		}

		// Run config.BeforeFunc
		if config.BeforeFunc != nil {
			testLogger.Debug("Start BeforeFunc")
			require.NoError(t, config.BeforeFunc(&BeforeFuncCtx{
				T:           t,
				Client:      client,
				ExecuteCmd:  executeCmd,
				Meta:        meta,
				OverrideEnv: overrideEnv,
				Logger:      testLogger,
			}))
			testLogger.Debug("End BeforeFunc")
		}

		// Run config.Cmd
		var result interface{}
		var exitCode int
		renderedArgs := []string(nil)
		rawArgs := config.Args
		if config.Cmd != "" {
			renderedArgs = cmdToArgs(meta, config.Cmd)
		} else {
			// We render raw arguments from meta
			for _, arg := range rawArgs {
				renderedArgs = append(renderedArgs, meta.render(arg))
			}
		}

		// We create a separate logger for the command we want to test.
		// This separate logger allow check function to test content log by a command
		// without content log by the test-engine (Before/After func, ...).
		cmdLoggerBuffer := &bytes.Buffer{}
		cmdLogger := &Logger{
			writer: io.MultiWriter(cmdLoggerBuffer, os.Stderr),
			level:  testLogger.level,
		}
		if len(renderedArgs) > 0 {
			stdout := &bytes.Buffer{}
			stderr := &bytes.Buffer{}
			exitCode, result, err = Bootstrap(&BootstrapConfig{
				Args:             renderedArgs,
				Commands:         config.Commands,
				BuildInfo:        buildInfo,
				Stdout:           stdout,
				Stderr:           stderr,
				Stdin:            stdin,
				Client:           client,
				DisableTelemetry: true,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           cmdLogger,
				HTTPClient:       httpClient,
			})

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
				Logger:      testLogger,
				LogBuffer:   cmdLoggerBuffer.String(),
			})
		}

		// Run config.AfterFunc
		if config.AfterFunc != nil {
			testLogger.Debug("Start AfterFunc")
			require.NoError(t, config.AfterFunc(&AfterFuncCtx{
				T:           t,
				Client:      client,
				ExecuteCmd:  executeCmd,
				Meta:        meta,
				CmdResult:   result,
				OverrideEnv: overrideEnv,
				Logger:      testLogger,
			}))
			testLogger.Debug("End AfterFunc")
		}
	}
}

func cmdToArgs(meta testMetadata, s string) []string {
	return strings.Split(meta.render(s), " ")
}

// BeforeFuncCombine combines multiple before functions into one.
func BeforeFuncCombine(beforeFuncs ...BeforeFunc) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
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
		if *UpdateCassettes {
			return beforeFunc(ctx)
		}
		return nil
	}
}

// AfterFuncCombine combines multiple after functions into one.
func AfterFuncCombine(afterFuncs ...AfterFunc) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
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
		ctx.Logger.Debugf("ExecStoreBeforeCmd: metaKey=%s args=%s\n", metaKey, args)
		ctx.Meta[metaKey] = ctx.ExecuteCmd(args)
		return nil
	}
}

func BeforeFuncOsExec(cmd string, args ...string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.Logger.Debugf("BeforeFuncOsExec: cmd=%s args=%s\n", cmd, args)
		return exec.Command(cmd, args...).Run()
	}
}

// ExecBeforeCmd executes the given before command.
func ExecBeforeCmd(cmd string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecBeforeCmd: args=%s\n", args)
		ctx.ExecuteCmd(args)
		return nil
	}
}

// ExecAfterCmd executes the given before command.
func ExecAfterCmd(cmd string) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecAfterCmd: args=%s\n", args)
		ctx.ExecuteCmd(args)
		return nil
	}
}

// TestCheckCombine combines multiple check functions into one.
func TestCheckCombine(checks ...TestCheck) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
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

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		actual := marshalGolden(t, ctx)

		goldenPath := getTestFilePath(t, ".golden")
		// In order to avoid diff in goldens we set all timestamp to the same date
		if *UpdateGoldens {
			require.NoError(t, os.MkdirAll(path.Dir(goldenPath), 0755))
			require.NoError(t, ioutil.WriteFile(goldenPath, []byte(actual), 0644)) //nolint:gosec
		}

		expected, err := ioutil.ReadFile(goldenPath)
		require.NoError(t, err, "expected to find golden file %s", goldenPath)
		assert.Equal(t, string(expected), actual)
	}
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
		assert.Equal(ctx.T, ctx.Meta.render(cmdStr), strings.Join(cmd.Args, " "))
		return exitCode, nil
	}
}

var regTimestamp = regexp.MustCompile(`(\d+-\d+-\d+T\d+:\d+:\d+\.\d+Z)`)

// uniformTimestamps replaces all timestamp to the date "1970-01-01T00:00:00.0Z"
func uniformTimestamps(input string) string {
	return regTimestamp.ReplaceAllString(input, "1970-01-01T00:00:00.0Z")
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
		i.Request.URL = regexp.MustCompile(`account\.scaleway\.com/tokens/[0-9a-f-]{36}`).ReplaceAllString(i.Request.URL, "account.scaleway.com/tokens/11111111-1111-1111-1111-111111111111")
		return nil
	})

	return &http.Client{Transport: &retryableHTTPTransport{transport: r}}, func() {
		assert.NoError(t, r.Stop()) // Make sure recorder is stopped once done with it
	}, nil
}

func marshalGolden(t *testing.T, ctx *CheckFuncCtx) string {
	jsonStderr := &bytes.Buffer{}
	jsonStdout := &bytes.Buffer{}

	jsonPrinter, err := NewPrinter(&PrinterConfig{
		OutputFlag: "json=pretty",
		Stdout:     jsonStdout,
		Stderr:     jsonStderr,
	})
	require.NoError(t, err)

	if ctx.Err != nil {
		err = jsonPrinter.Print(ctx.Err, nil)
		require.NoError(t, err)
	}
	if ctx.Result != nil {
		err = jsonPrinter.Print(ctx.Result, nil)
		require.NoError(t, err)
	}

	buffer := bytes.Buffer{}
	buffer.WriteString(fmt.Sprintf("\U0001F3B2\U0001F3B2\U0001F3B2 EXIT CODE: %d \U0001F3B2\U0001F3B2\U0001F3B2\n", ctx.ExitCode))

	if len(ctx.Stdout) > 0 {
		buffer.WriteString("\U0001F7E9\U0001F7E9\U0001F7E9 STDOUT️ \U0001F7E9\U0001F7E9\U0001F7E9️\n")
		buffer.Write(ctx.Stdout)
	}

	if len(ctx.Stderr) > 0 {
		buffer.WriteString("\U0001F7E5\U0001F7E5\U0001F7E5 STDERR️️ \U0001F7E5\U0001F7E5\U0001F7E5️\n")
		buffer.Write(ctx.Stderr)
	}

	if jsonStdout.Len() > 0 {
		buffer.WriteString("\U0001F7E9\U0001F7E9\U0001F7E9 JSON STDOUT \U0001F7E9\U0001F7E9\U0001F7E9\n")
		buffer.Write(jsonStdout.Bytes())
	}

	if jsonStderr.Len() > 0 {
		buffer.WriteString("\U0001F7E5\U0001F7E5\U0001F7E5 JSON STDERR \U0001F7E5\U0001F7E5\U0001F7E5\n")
		buffer.Write(jsonStderr.Bytes())
	}

	str := buffer.String()
	// In order to avoid diff in goldens we set all timestamp to the same date
	str = uniformTimestamps(str)
	// Replace Windows return carriage.
	str = strings.ReplaceAll(str, "\r", "")
	return str
}
