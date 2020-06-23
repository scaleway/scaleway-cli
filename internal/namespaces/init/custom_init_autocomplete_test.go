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

func baseBeforeFunc() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta["SecretKey"], _ = ctx.Client.GetSecretKey()
		ctx.Meta["OrganizationID"], _ = ctx.Client.GetDefaultOrganizationID()
		return nil
	}
}

func Test_InitAutocomplete(t *testing.T) {
	defaultSettings := map[string]string{
		"secret-key":       "{{ .SecretKey }}",
		"organization-id":  "{{ .OrganizationID }}",
		"send-telemetry":   "false",
		"remove-v1-config": "false",
		"with-ssh-key":     "false",
	}

	t.Run("Without", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: baseBeforeFunc(),
		Cmd:        appendArgs("scw init install-autocomplete=false", defaultSettings),
		Check:      core.TestCheckGolden(),
		TmpHomeDir: true,
	}))

	t.Run("Zsh", func(t *testing.T) {
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=zsh)"
`
		core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: baseBeforeFunc(),
			Cmd:        appendArgs("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".zshrc")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						require.NoError(t, err)
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

	t.Run("bash-mac", func(t *testing.T) {
		if runtime.GOOS != "darwin" {
			t.SkipNow()
		}
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=bash)"
`
		core.Test(&core.TestConfig{
			Commands:   GetCommands(),
			BeforeFunc: baseBeforeFunc(),
			Cmd:        appendArgs("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					filePath := path.Join(homeDir, ".bash_profile")
					fileContent, err := ioutil.ReadFile(filePath)
					if err != nil {
						require.NoError(t, err)
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

	t.Run("bash-linux", func(t *testing.T) {
		if runtime.GOOS != "linux" {
			t.SkipNow()
		}
		evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=bash)"
`
		core.Test(&core.TestConfig{
			BeforeFunc: baseBeforeFunc(),
			Commands:   GetCommands(),
			Cmd:        appendArgs("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
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
			BeforeFunc: core.BeforeFuncCombine(
				func(ctx *core.BeforeFuncCtx) error {
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
				baseBeforeFunc(),
			),
			Cmd: appendArgs("scw init install-autocomplete=true", defaultSettings),
			Check: core.TestCheckCombine(
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
