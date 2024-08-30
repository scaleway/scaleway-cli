package init_test

import (
	"os"
	"path"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	initCLI "github.com/scaleway/scaleway-cli/v2/internal/namespaces/init" // alias required to not collide with go init func
	"github.com/stretchr/testify/require"
)

func baseBeforeFunc() core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		ctx.Meta["AccessKey"], _ = ctx.Client.GetAccessKey()
		ctx.Meta["SecretKey"], _ = ctx.Client.GetSecretKey()
		ctx.Meta["ProjectID"], _ = ctx.Client.GetDefaultProjectID()
		ctx.Meta["OrganizationID"], _ = ctx.Client.GetDefaultOrganizationID()
		return nil
	}
}

const (
	linux   = "linux"
	darwin  = "darwin"
	windows = "windows"
)

func Test_InitAutocomplete(t *testing.T) {
	defaultSettings := map[string]string{
		"access-key":      "{{ .AccessKey }}",
		"secret-key":      "{{ .SecretKey }}",
		"send-telemetry":  "false",
		"with-ssh-key":    "false",
		"organization-id": "{{ .OrganizationID }}",
		"project-id":      "{{ .ProjectID }}",
	}

	runAllShells := func(t *testing.T) {
		t.Helper()
		t.Run("Without", core.Test(&core.TestConfig{
			Commands:   initCLI.GetCommands(),
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
				Commands:   initCLI.GetCommands(),
				BeforeFunc: baseBeforeFunc(),
				Cmd:        appendArgs("scw init install-autocomplete=true", defaultSettings),
				Check: core.TestCheckCombine(
					core.TestCheckGolden(),
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						t.Helper()
						if runtime.GOOS == windows {
							// autocomplete installation is not yet supported on windows
							return
						}
						homeDir := ctx.OverrideEnv["HOME"]
						filePath := path.Join(homeDir, ".zshrc")
						fileContent, err := os.ReadFile(filePath)
						require.NoError(t, err)
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

		t.Run("fish", func(t *testing.T) {
			evalLine := `
# Scaleway CLI autocomplete initialization.
eval (scw autocomplete script shell=fish)
`
			core.Test(&core.TestConfig{
				Commands: initCLI.GetCommands(),
				BeforeFunc: core.BeforeFuncCombine(
					func(ctx *core.BeforeFuncCtx) error {
						homeDir := ctx.OverrideEnv["HOME"]
						configPath := path.Join(homeDir, ".config", "fish", "config.fish")

						// Ensure the subfolders for the configuration files are all created
						err := os.MkdirAll(filepath.Dir(configPath), 0o755)
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
						t.Helper()
						if runtime.GOOS == windows {
							// autocomplete installation is not yet supported on windows
							return
						}
						homeDir := ctx.OverrideEnv["HOME"]
						filePath := path.Join(homeDir, ".config", "fish", "config.fish")
						fileContent, err := os.ReadFile(filePath)
						require.NoError(t, err)
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

		t.Run("bash", func(t *testing.T) {
			evalLine := `
# Scaleway CLI autocomplete initialization.
eval "$(scw autocomplete script shell=bash)"
`
			core.Test(&core.TestConfig{
				Commands:   initCLI.GetCommands(),
				BeforeFunc: baseBeforeFunc(),
				Cmd:        appendArgs("scw init install-autocomplete=true", defaultSettings),
				Check: core.TestCheckCombine(
					core.TestCheckGolden(),
					func(t *testing.T, ctx *core.CheckFuncCtx) {
						t.Helper()
						homeDir := ctx.OverrideEnv["HOME"]
						filePath := ""
						switch runtime.GOOS {
						case linux:
							filePath = path.Join(homeDir, ".bashrc")
						case darwin:
							filePath = path.Join(homeDir, ".bash_profile")
						case windows:
							// autocomplete installation is not yet supported on windows
							return
						default:
							t.Fatalf("unsupported OS")
						}
						fileContent, err := os.ReadFile(filePath)
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
	}

	t.Run(darwin, func(t *testing.T) {
		if runtime.GOOS != darwin {
			t.SkipNow()
		}
		runAllShells(t)
	})

	t.Run(linux, func(t *testing.T) {
		if runtime.GOOS != linux {
			t.SkipNow()
		}
		runAllShells(t)
	})

	t.Run(windows, func(t *testing.T) {
		if runtime.GOOS != windows {
			t.SkipNow()
		}
		runAllShells(t)
	})
}
