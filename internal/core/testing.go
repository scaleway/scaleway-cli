package core

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

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

func getTestClient(t *testing.T, e2eClient bool) *scw.Client {
	clientOpts := []scw.ClientOption{
		scw.WithDefaultRegion(scw.RegionFrPar),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithEnv(),
		scw.WithUserAgent("cli-e2e-test"),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
	}

	client, err := scw.NewClient(clientOpts...)
	panicIfErr(err)

	if !e2eClient {
		return client
	}

	res, err := test.NewAPI(client).Register(&test.RegisterRequest{Username: "sidi"})
	panicIfErr(err)

	client, err = scw.NewClient(append(clientOpts, scw.WithAuth(res.AccessKey, res.SecretKey))...)
	panicIfErr(err)

	return client
}

// Run a CLI integration test. See TestConfig for configuration option
func Test(config *TestConfig) func(t *testing.T) {
	return func(t *testing.T) {

		stdout := &bytes.Buffer{}
		stderr := &bytes.Buffer{}
		client := getTestClient(t, config.UseE2EClient)

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
		testName := t.Name()
		testName = strings.Replace(testName, "/", "-", -1)
		testName = strcase.ToBashArg(testName)
		testGolden(t, testName+".stderr", result.Stderr)
	}
}

// TestCheckStdoutGolden assert stdout using golden
func TestCheckStdoutGolden() TestCheck {
	return func(t *testing.T, result *TestResult) {
		testName := t.Name()
		testName = strings.Replace(testName, "/", "-", -1)
		testName = strcase.ToBashArg(testName)
		testGolden(t, testName+".stdout", result.Stdout)
	}
}

// TestCheckGolden assert stderr and stdout using golden
func TestCheckGolden() TestCheck {
	return TestCheckCombine(
		TestCheckStdoutGolden(),
		TestCheckStderrGolden(),
	)
}

func testGolden(t *testing.T, goldenName string, actual []byte) {
	goldenPath := filepath.Join(".", "testdata", goldenName+".golden")
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
