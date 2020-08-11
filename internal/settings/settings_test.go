package settings

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

const emptyFile = ""

// TestLoadConfig tests config getters return correct values
func Test_Load(t *testing.T) {
	// create home dir
	dir := initEnv(t)
	tests := []struct {
		name  string
		env   map[string]string
		files map[string]string

		expectedError  string
		expectedDebug  *bool
		expectedOutput string
	}{
		{
			name: "Valid YAML but empty settings",
			files: map[string]string{
				path.Join(".config", "scw", "cli.yaml"): emptyFile,
			},
			env: map[string]string{
				"HOME": "{HOME}",
			},
			expectedError: "cli settings file does not exist",
		},
		{
			name: "Invalid YAML",
			files: map[string]string{
				path.Join(".config", "scw", "cli.yaml"): `foo;bar`,
			},
			env: map[string]string{
				"HOME": "{HOME}",
			},
			expectedError: fmt.Sprintf("content of cli settings file %s is invalid", path.Join(dir, ".config", "scw", "cli.yaml")),
		},
		{
			name: "Valid YAML with extra keys",
			files: map[string]string{
				path.Join(".config", "scw", "cli.yaml"): `output: json
foo: bar
`,
			},
			env: map[string]string{
				"HOME": "{HOME}",
			},
			expectedOutput: "json",
		},
		{
			name: "Complete config",
			files: map[string]string{
				path.Join(".config", "scw", "cli.yaml"): `output: json=pretty`,
			},
			env: map[string]string{
				"HOME": "{HOME}",
			},
			expectedOutput: "json=pretty",
		},
	}

	// delete home dir and reset env variables
	defer resetEnv(t, os.Environ(), dir)
	logger.EnableDebugMode()
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// set up env and config file(s)
			setEnv(t, test.env, test.files, dir)
			test.expectedError = strings.Replace(test.expectedError, "{HOME}", dir, -1)

			// remove config file(s)
			defer cleanEnv(t, test.files, dir)

			config, err := Load()
			if test.expectedError == "" {
				assert.NoError(t, err)

				// assert getters
				assert.Equal(t, test.expectedOutput, *config.Output)
			} else {
				if err != nil {
					assert.Equal(t, test.expectedError, err.Error())
				}
			}
		})
	}
}

func initEnv(t *testing.T) string {
	dir, err := ioutil.TempDir("", "home")
	if err != nil {
		t.Fatal(err)
	}
	return dir
}

func cleanEnv(t *testing.T, files map[string]string, homeDir string) {
	for path := range files {
		assert.NoError(t, os.RemoveAll(filepath.Join(homeDir, path)))
	}
}

func setEnv(t *testing.T, env, files map[string]string, homeDir string) {
	os.Clearenv()
	for key, value := range env {
		value = strings.Replace(value, "{HOME}", homeDir, -1)
		assert.NoError(t, os.Setenv(key, value))
	}

	for path, content := range files {
		targetPath := filepath.Join(homeDir, path)
		assert.NoError(t, os.MkdirAll(filepath.Dir(targetPath), 0700))
		assert.NoError(t, ioutil.WriteFile(targetPath, []byte(content), defaultPermission))
	}
}

// function taken from https://golang.org/src/os/env_test.go
func resetEnv(t *testing.T, origEnv []string, homeDir string) {
	assert.NoError(t, os.RemoveAll(homeDir))
	for _, pair := range origEnv {
		// Environment variables on Windows can begin with =
		// https://blogs.msdn.com/b/oldnewthing/archive/2010/05/06/10008132.aspx
		i := strings.Index(pair[1:], "=") + 1
		if err := os.Setenv(pair[:i], pair[i+1:]); err != nil {
			t.Errorf("Setenv(%q, %q) failed during reset: %v", pair[:i], pair[i+1:], err)
		}
	}
}

func s(value string) *string {
	return &value
}

func TestConfig_ConfigFile(t *testing.T) {
	type testCase struct {
		config *Settings
		result string
	}

	run := func(c *testCase) func(t *testing.T) {
		return func(t *testing.T) {
			config, err := c.config.HumanSettings()
			assert.NoError(t, err)
			assert.Equal(t, c.result, config)

			loaded, err2 := unmarshalSettings([]byte(config))
			assert.NoError(t, err2)
			assert.Equal(t, c.config, loaded)
		}
	}

	t.Run("empty", run(&testCase{
		config: &Settings{},
		result: `# Scaleway CLI settings file
# This settings file can be used only with Scaleway CLI (>2.0.0) (https://github.com/scaleway/scaleway-cli)
# Output sets the output format for all commands you run
# output: human
`,
	}))

	t.Run("partial", run(&testCase{
		config: &Settings{
			Output: s("json")},
		result: `# Scaleway CLI settings file
# This settings file can be used only with Scaleway CLI (>2.0.0) (https://github.com/scaleway/scaleway-cli)
# Output sets the output format for all commands you run
output: json
`,
	}))

	t.Run("full", run(&testCase{
		config: &Settings{
			Output: s("json:pretty"),
		},
		result: `# Scaleway CLI settings file
# This settings file can be used only with Scaleway CLI (>2.0.0) (https://github.com/scaleway/scaleway-cli)
# Output sets the output format for all commands you run
output: json:pretty
`,
	}))
}

func TestEmptyConfig(t *testing.T) {
	assert.True(t, (&Settings{}).IsEmpty())
	assert.False(t, (&Settings{Output: s("json")}).IsEmpty())
}
