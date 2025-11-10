package core

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"math/rand/v2"
	"net/http"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"text/template"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	args "github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/platform/terminal"
	"github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const TestBucketNamePrefix = "cli-test-bucket"

// Test flags
// You can create a binary of each test using "go test -c -o myBinary"
var (
	// UpdateGoldens will update all the golden files of a given test
	UpdateGoldens = flag.Bool(
		"goldens",
		os.Getenv("CLI_UPDATE_GOLDENS") == "true",
		"Record goldens",
	)

	// UpdateCassettes will update all cassettes of a given test
	UpdateCassettes = flag.Bool(
		"cassettes",
		os.Getenv("CLI_UPDATE_CASSETTES") == "true",
		"Record Cassettes",
	)

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
	Result any

	// Meta bag
	Meta TestMetadata

	// Scaleway client
	Client *scw.Client

	// OverrideEnv passed in the TestConfig
	OverrideEnv map[string]string

	Logger *Logger

	// The content logged by the command
	LogBuffer string
}

var testRenderHelpers = map[string]any{
	"randint": func() string {
		return strconv.FormatUint(
			rand.Uint64(),
			10,
		)
	},
}

// TestMetadata contains arbitrary data that can be passed along a test lifecycle.
type TestMetadata map[string]any

// Render renders a go template using where content of Meta can be used
func (meta TestMetadata) Render(strTpl string) string {
	t := meta["t"].(*testing.T)
	buf := &bytes.Buffer{}
	require.NoError(
		t,
		template.Must(template.New("tpl").Funcs(testRenderHelpers).Parse(strTpl)).
			Execute(buf, meta),
	)

	return buf.String()
}

func BeforeFuncStoreInMeta(key string, value any) BeforeFunc {
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
	Meta   TestMetadata
	Client *scw.Client
}

type OverrideExecTestFunc func(ctx *ExecFuncCtx, cmd *exec.Cmd) (exitCode int, err error)

type BeforeFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) any
	Meta        TestMetadata
	OverrideEnv map[string]string
	Logger      *Logger
}

type AfterFuncCtx struct {
	T           *testing.T
	Client      *scw.Client
	ExecuteCmd  func(args []string) any
	Meta        TestMetadata
	CmdResult   any
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

	// EnabledAliases enables aliases that are disabled in tests
	EnableAliases bool
}

func (config *TestConfig) DebugString() string {
	return config.Cmd
}

// getTestFilePath returns a valid filename path based on the go test name and suffix. (Take care of non fs friendly char)
func getTestFilePath(t *testing.T, suffix string) string {
	t.Helper()
	specialChars := regexp.MustCompile(`[\\?%*:|"<>. ]`)

	// Replace nested tests separators.
	fileName := strings.ReplaceAll(t.Name(), "/", "-")

	fileName = strcase.ToBashArg(fileName)

	// Replace special characters.
	fileName = specialChars.ReplaceAllLiteralString(fileName, "") + suffix

	return filepath.Join(".", "testdata", fileName)
}

func createTestClient(
	t *testing.T,
	testConfig *TestConfig,
	httpClient *http.Client,
) (client *scw.Client) {
	t.Helper()
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
			Transport: &retryableHTTPTransport{transport: &SocketPassthroughTransport{}},
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

		client, err = scw.NewClient(
			append(clientOpts, scw.WithAuth(res.AccessKey, res.SecretKey))...)
		require.NoError(t, err)
	}

	return client
}

// DefaultRetryInterval is used across all wait functions in the CLI
// In particular it is very handy to define this RetryInterval at 0 second while running cassette in testing
// because they will be executed without waiting.
var DefaultRetryInterval *time.Duration

var foldersUsingVCRv4 = []string{
	"instance",
	"k8s",
	"marketplace",
}

func folderUsesVCRv4(fullFolderPath string) bool {
	fullPathSplit := strings.Split(fullFolderPath, string(os.PathSeparator))

	folder := fullPathSplit[len(fullPathSplit)-2]
	for _, migratedFolder := range foldersUsingVCRv4 {
		if migratedFolder == folder {
			return true
		}
	}

	return false
}

// Run a CLI integration test. See TestConfig for configuration option
func Test(config *TestConfig) func(t *testing.T) {
	return func(t *testing.T) {
		t.Helper()
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
		human.RegisterMarshalerFunc(
			time.Time{},
			func(_ any, _ *human.MarshalOpt) (string, error) {
				return "few seconds ago", nil
			},
		)

		if !*UpdateCassettes {
			tmp := 0 * time.Second
			DefaultRetryInterval = &tmp
		}

		ctx := config.Ctx
		if ctx == nil {
			ctx = t.Context()
		}
		if len(config.PromptResponseMocks) > 0 {
			ctx = interactive.InjectMockResponseToContext(ctx, config.PromptResponseMocks)
		}

		folder, err := os.Getwd()
		if err != nil {
			t.Fatalf("cannot detect working directory for testing")
		}

		// Create an HTTP client with recording capabilities
		var (
			httpClient *http.Client
			cleanup    func()
		)

		if folderUsesVCRv4(folder) {
			httpClient, cleanup, err = newHTTPRecorder(t, folder, *UpdateCassettes)
		} else {
			httpClient, cleanup, err = getHTTPRecoder(t, *UpdateCassettes)
		}

		require.NoError(t, err)
		defer cleanup()

		// We try to use the client provided in the config
		// if no client is provided in the config we create a test client
		client := config.Client
		if client == nil {
			client = createTestClient(t, config, httpClient)
		}

		meta := TestMetadata{
			"t": t,
		}

		overrideEnv := config.OverrideEnv
		if overrideEnv == nil {
			overrideEnv = map[string]string{}
		}

		if config.TmpHomeDir {
			dir := t.TempDir()
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
				Version:         version.Must(version.NewSemver("v0.0.0+test")),
				BuildDate:       "unknown",
				GoVersion:       "runtime.Version()",
				GitBranch:       "unknown",
				GitCommit:       "unknown",
				GoArch:          "runtime.GOARCH",
				GoOS:            "runtime.GOOS",
				UserAgentPrefix: "scaleway-cli",
			}
		}

		executeCmd := func(args []string) any {
			stdoutBuffer := &bytes.Buffer{}
			stderrBuffer := &bytes.Buffer{}
			_, result, err := Bootstrap(&BootstrapConfig{
				Args:             args,
				Commands:         config.Commands.Copy(), // Copy commands to ensure they are not modified
				BuildInfo:        buildInfo,
				Stdout:           stdoutBuffer,
				Stderr:           stderrBuffer,
				Client:           client,
				DisableTelemetry: true,
				DisableAliases:   !config.EnableAliases,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           testLogger,
				HTTPClient:       httpClient,
				Platform:         terminal.NewPlatform(buildInfo.GetUserAgent()),
			})
			require.NoError(
				t,
				err,
				"error executing cmd (%s)\nstdout: %s\nstderr: %s",
				args,
				stdoutBuffer.String(),
				stderrBuffer.String(),
			)

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
			}), "error executing BeforeFunc")
			testLogger.Debug("End BeforeFunc")
		}

		// Run config.Cmd
		var result any
		var exitCode int
		renderedArgs := []string(nil)
		rawArgs := config.Args
		if config.Cmd != "" {
			renderedArgs = cmdToArgs(meta, config.Cmd)
		} else {
			// We Render raw arguments from meta
			for _, arg := range rawArgs {
				renderedArgs = append(renderedArgs, meta.Render(arg))
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
				DisableAliases:   !config.EnableAliases,
				OverrideEnv:      overrideEnv,
				OverrideExec:     overrideExec,
				Ctx:              ctx,
				Logger:           cmdLogger,
				HTTPClient:       httpClient,
				Platform:         terminal.NewPlatform(buildInfo.GetUserAgent()),
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

func cmdToArgs(meta TestMetadata, s string) []string {
	return strings.Split(meta.Render(s), " ")
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

func AfterFuncWhenUpdatingCassette(afterFunc AfterFunc) AfterFunc {
	return func(ctx *AfterFuncCtx) error {
		if *UpdateCassettes {
			return afterFunc(ctx)
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

// ExecStoreBeforeCmdWithResulter executes the given before command and register the result
// in the context Meta at metaKey. The result is transformed by the resulter function.
func ExecStoreBeforeCmdWithResulter(metaKey, cmd string, resulter func(any) any) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		args := cmdToArgs(ctx.Meta, cmd)
		ctx.Logger.Debugf("ExecStoreBeforeCmd: metaKey=%s args=%s\n", metaKey, args)
		result := ctx.ExecuteCmd(args)
		ctx.Meta[metaKey] = resulter(result)

		return nil
	}
}

func BeforeFuncOsExec(cmd string, args ...string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		ctx.Logger.Debugf("BeforeFuncOsExec: cmd=%s args=%s\n", cmd, args)
		err := exec.Command(cmd, args...).Run()
		if err != nil {
			formattedCmd := strings.Join(append([]string{cmd}, args...), " ")

			return fmt.Errorf("failed to execute cmd %q: %w", formattedCmd, err)
		}

		return nil
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

// ExecBeforeCmdArgs executes the given command before command.
func ExecBeforeCmdArgs(args []string) BeforeFunc {
	return func(ctx *BeforeFuncCtx) error {
		for i := range args {
			args[i] = ctx.Meta.Render(args[i])
		}
		ctx.Logger.Debugf("ExecBeforeCmdArgs: args=%s\n", args)
		ctx.ExecuteCmd(args)

		return nil
	}
}

// ExecBeforeCmdWithResult executes the given command and returns its result.
func ExecBeforeCmdWithResult(ctx *BeforeFuncCtx, cmd string) any {
	args := cmdToArgs(ctx.Meta, cmd)
	ctx.Logger.Debugf("ExecBeforeCmd: args=%s\n", args)

	return ctx.ExecuteCmd(args)
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
		t.Helper()
		for _, check := range checks {
			check(t, ctx)
		}
	}
}

// TestCheckExitCode assert exitCode
func TestCheckExitCode(expectedCode int) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		assert.Equal(t, expectedCode, ctx.ExitCode, "Invalid exit code\n%s", string(ctx.Stderr))
	}
}

// GoldenReplacement describe patterns to be replaced in goldens
type GoldenReplacement struct {
	// Line will be matched using this regex
	Pattern *regexp.Regexp
	// Content that will replace the matched regex
	// This is the format for repl in (*regexp.Regexp).ReplaceAll
	// You can use $ to represent groups $1, $2...
	Replacement string

	// OptionalMatch allow the golden to not contain the given patterns
	// if false, the golden must contain the given pattern
	OptionalMatch bool
}

// GoldenReplacePatterns replace the list of patterns with their given replacement
func GoldenReplacePatterns(golden string, replacements ...GoldenReplacement) (string, error) {
	var matchFailed []string
	changedGolden := golden

	for _, replacement := range replacements {
		if !replacement.Pattern.MatchString(changedGolden) {
			if !replacement.OptionalMatch {
				matchFailed = append(matchFailed, replacement.Pattern.String())
			}

			continue
		}
		changedGolden = replacement.Pattern.ReplaceAllString(changedGolden, replacement.Replacement)
	}

	if len(matchFailed) > 0 {
		return changedGolden, fmt.Errorf("failed to match regex in golden: %#q", matchFailed)
	}

	return changedGolden, nil
}

// TestCheckGoldenAndReplacePatterns assert stderr and stdout using golden,
// golden are matched against given regex and edited with replacements
func TestCheckGoldenAndReplacePatterns(replacements ...GoldenReplacement) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		actual := marshalGolden(t, ctx)
		actual, actualReplaceErr := GoldenReplacePatterns(actual, replacements...)

		goldenPath := getTestFilePath(t, ".golden")
		// In order to avoid diff in goldens we set all timestamp to the same date
		if *UpdateGoldens {
			require.NoError(t, os.MkdirAll(path.Dir(goldenPath), 0o755))
			require.NoError(t, os.WriteFile(goldenPath, []byte(actual), 0o644)) //nolint:gosec
		}

		expected, err := os.ReadFile(goldenPath)
		require.NoError(t, err, "expected to find golden file %s", goldenPath)
		assert.Equal(t, string(expected), actual)
		assert.NoError(t, actualReplaceErr, "failed to match test output with regexes")
	}
}

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		actual := marshalGolden(t, ctx)

		goldenPath := getTestFilePath(t, ".golden")
		// In order to avoid diff in goldens we set all timestamp to the same date
		if *UpdateGoldens {
			require.NoError(t, os.MkdirAll(path.Dir(goldenPath), 0o755))
			require.NoError(t, os.WriteFile(goldenPath, []byte(actual), 0o644)) //nolint:gosec
		}

		expected, err := os.ReadFile(goldenPath)
		require.NoError(t, err, "expected to find golden file %s", goldenPath)
		assert.Equal(t, string(expected), actual)
	}
}

// TestCheckS3Golden assert stderr and stdout using golden, and omits the random suffix in the bucket name
func TestCheckS3Golden() TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		actual := marshalGolden(t, ctx)
		normalizedActual := removeRandomPrefixFromOutput(actual)

		goldenPath := getTestFilePath(t, ".golden")
		// In order to avoid diff in goldens we set all timestamp to the same date
		if *UpdateGoldens {
			require.NoError(t, os.MkdirAll(path.Dir(goldenPath), 0o755))
			require.NoError(
				t,
				os.WriteFile(goldenPath, []byte(normalizedActual), 0o644),
			)
		}

		expected, err := os.ReadFile(goldenPath)
		require.NoError(t, err, "expected to find golden file %s", goldenPath)
		normalizedExpected := removeRandomPrefixFromOutput(string(expected))

		assert.Equal(t, normalizedExpected, normalizedActual)
	}
}

func removeRandomPrefixFromOutput(output string) string {
	begin := strings.Index(output, TestBucketNamePrefix)
	if begin < 0 {
		return output
	}
	end := strings.IndexByte(output[begin:], '\n')
	actualBucketName := output[begin : begin+end]
	normalizedBucketName := strings.TrimRight(actualBucketName, "0123456789")

	return strings.ReplaceAll(output, actualBucketName, normalizedBucketName)
}

// TestCheckError asserts error
func TestCheckError(err error) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		assert.Equal(t, err, ctx.Err, "Invalid error")
	}
}

// TestCheckStdout asserts stdout using string
func TestCheckStdout(stdout string) TestCheck {
	return func(t *testing.T, ctx *CheckFuncCtx) {
		t.Helper()
		assert.Equal(t, stdout, string(ctx.Stdout), "Invalid stdout")
	}
}

func OverrideExecSimple(cmdStr string, exitCode int) OverrideExecTestFunc {
	return func(ctx *ExecFuncCtx, cmd *exec.Cmd) (int, error) {
		assert.Equal(ctx.T, ctx.Meta.Render(cmdStr), strings.Join(cmd.Args, " "))

		return exitCode, nil
	}
}

var regTimestamp = regexp.MustCompile(`(\d+-\d+-\d+T\d+:\d+:\d+\.\d+Z)`)

// uniformTimestamps replaces all timestamp to the date "1970-01-01T00:00:00.0Z"
func uniformTimestamps(input string) string {
	return regTimestamp.ReplaceAllString(input, "1970-01-01T00:00:00.0Z")
}

func validateJSONGolden(t *testing.T, jsonStdout, jsonStderr *bytes.Buffer) {
	t.Helper()
	var jsonInterface any
	if jsonStdout.Len() > 0 {
		err := json.Unmarshal(jsonStdout.Bytes(), &jsonInterface)
		require.NoError(t, err, "json stdout is invalid (%s)", getTestFilePath(t, ".cassette"))
	}
	if jsonStderr.Len() > 0 {
		err := json.Unmarshal(jsonStderr.Bytes(), &jsonInterface)
		require.NoError(t, err, "json stderr is invalid (%s)", getTestFilePath(t, ".cassette"))
	}
}

func marshalGolden(t *testing.T, ctx *CheckFuncCtx) string {
	t.Helper()
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

	if _, isRawResult := ctx.Result.(RawResult); !isRawResult {
		validateJSONGolden(t, jsonStdout, jsonStderr)
	}

	buffer := bytes.Buffer{}
	buffer.WriteString(
		fmt.Sprintf(
			"\U0001F3B2\U0001F3B2\U0001F3B2 EXIT CODE: %d \U0001F3B2\U0001F3B2\U0001F3B2\n",
			ctx.ExitCode,
		),
	)

	if len(ctx.Stdout) > 0 {
		buffer.WriteString(
			"\U0001F7E9\U0001F7E9\U0001F7E9 STDOUT️ \U0001F7E9\U0001F7E9\U0001F7E9️\n",
		)
		buffer.Write(ctx.Stdout)
	}

	if len(ctx.Stderr) > 0 {
		buffer.WriteString(
			"\U0001F7E5\U0001F7E5\U0001F7E5 STDERR️️ \U0001F7E5\U0001F7E5\U0001F7E5️\n",
		)
		buffer.Write(ctx.Stderr)
	}

	if jsonStdout.Len() > 0 {
		buffer.WriteString(
			"\U0001F7E9\U0001F7E9\U0001F7E9 JSON STDOUT \U0001F7E9\U0001F7E9\U0001F7E9\n",
		)
		buffer.Write(jsonStdout.Bytes())
	}

	if jsonStderr.Len() > 0 {
		buffer.WriteString(
			"\U0001F7E5\U0001F7E5\U0001F7E5 JSON STDERR \U0001F7E5\U0001F7E5\U0001F7E5\n",
		)
		buffer.Write(jsonStderr.Bytes())
	}

	str := buffer.String()
	// In order to avoid diff in goldens we set all timestamp to the same date
	str = uniformTimestamps(str)
	// Replace Windows return carriage.
	str = strings.ReplaceAll(str, "\r", "")

	return str
}
