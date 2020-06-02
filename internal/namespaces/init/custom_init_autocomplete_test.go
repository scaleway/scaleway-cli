package init

import (
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/stretchr/testify/require"
)

func Test_InitAutocomplete(t *testing.T) {
	secretKey := dummyUUID
	organizationID := dummyUUID
	// if you are recording, you must place a valid token in the environment variable SCW_TEST_SECRET_KEY
	if os.Getenv("SCW_TEST_SECRET_KEY") != "" {
		secretKey = os.Getenv("SCW_TEST_SECRET_KEY")
	}
	defaultSettings := map[string]string{
		"secret-key":       secretKey,
		"organization-id":  organizationID,
		"send-telemetry":   "false",
		"remove-v1-config": "false",
		"with-ssh-key":     "false",
	}

	t.Run("Without", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      cmdFromSettings("scw init install-autocomplete=false", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
			),
			TmpHomeDir: true,
		})(t)
	})

	t.Run("Zsh", func(t *testing.T) {
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=zsh)"
`
		core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      cmdFromSettings("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".zshrc")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						t.FailNow()
					}
					require.Equal(t, evalLine, string(fileContent))
				},
			),
			TmpHomeDir: true,
			OverrideEnv: map[string]string{
				"SHELL": "/usr/local/bin/zsh",
			},
			PromptResponseMocks: []string{
				// What type of shell are you using
				"zsh",
				// Do you want to proceed with these changes? (Y/n):
				"yes",
			},
		})(t)
	})

	t.Run("bash (mac)", func(t *testing.T) {
		if runtime.GOOS != "darwin" {
			t.SkipNow()
		}
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=bash)"
`
		core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      cmdFromSettings("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".bash_profile")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						t.FailNow()
					}
					require.Equal(t, evalLine, string(fileContent))
				},
			),
			TmpHomeDir: true,
			OverrideEnv: map[string]string{
				"SHELL": "/usr/local/bin/bash",
			},
			PromptResponseMocks: []string{
				// What type of shell are you using
				"bash",
				// Do you want to proceed with these changes? (Y/n):
				"yes",
			},
		})(t)
	})

	t.Run("bash (linux)", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.SkipNow()
		}
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=bash)"
`
		core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      cmdFromSettings("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".bashrc")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						t.FailNow()
					}
					require.Equal(t, evalLine, string(fileContent))
				},
			),
			TmpHomeDir: true,
			OverrideEnv: map[string]string{
				"SHELL": "/usr/local/bin/bash",
			},
			PromptResponseMocks: []string{
				// What type of shell are you using
				"bash",
				// Do you want to proceed with these changes? (Y/n):
				"yes",
			},
		})(t)
	})

	t.Run("fish", func(t *testing.T) {
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval (scw autocomplete script shell=fish)
`
		core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
				homeDir := ctx.OverrideEnv["HOME"]
				configPath := path.Join(homeDir, ".config", "fish", "config.fish")

				// Ensure the subfolders for the configuration files are all created
				err := os.MkdirAll(filepath.Dir(configPath), 0755)
				if err != nil {
					return err
				}

				// Write the configuration file
				err = ioutil.WriteFile(configPath, []byte(``), 0600)
				if err != nil {
					return err
				}
				return nil
			},
			Cmd: cmdFromSettings("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".config", "fish", "config.fish")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						t.FailNow()
					}
					require.Equal(t, evalLine, string(fileContent))
				},
			),
			TmpHomeDir: true,
			OverrideEnv: map[string]string{
				"SHELL": "/usr/local/bin/fish",
			},
			PromptResponseMocks: []string{
				// What type of shell are you using
				"fish",
				// Do you want to proceed with these changes? (Y/n):
				"yes",
			},
		})(t)
	})
}
