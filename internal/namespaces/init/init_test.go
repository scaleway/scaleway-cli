package init

import (
	"fmt"
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

const dummyUUID = "11111111-1111-1111-1111-111111111111"

func checkConfig(f func(t *testing.T, config *scw.Config)) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		homeDir := ctx.OverrideEnv["HOME"]
		config, err := scw.LoadConfigFromPath(path.Join(homeDir, ".config", "scw", "config.yaml"))
		require.NoError(t, err)
		f(t, config)
	}
}

func cmdFromSettings(prefix string, settings map[string]string) string {
	res := prefix
	for k, v := range settings {
		res += fmt.Sprintf(" %s=%s", k, v)
	}
	return res
}

func TestInit(t *testing.T) {
	secretKey := dummyUUID
	organizationID := dummyUUID
	// if you are recording, you must place a valid token in the environment variable SCW_TEST_SECRET_KEY
	if os.Getenv("SCW_TEST_SECRET_KEY") != "" {
		secretKey = os.Getenv("SCW_TEST_SECRET_KEY")
	}
	defaultSettings := map[string]string{
		"secret-key":           secretKey,
		"organization-id":      organizationID,
		"send-telemetry":       "true",
		"install-autocomplete": "false",
		"remove-v1-config":     "false",
		"with-ssh-key":         "false",
	}

	t.Run("Simple", func(t *testing.T) {
		core.Test(&core.TestConfig{
			Commands:            GetCommands(),
			PromptResponseMocks: []string{},
			TmpHomeDir:          true,
			Cmd:                 cmdFromSettings("scw init", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				checkConfig(func(t *testing.T, config *scw.Config) {
					assert.Equal(t, secretKey, *config.SecretKey)
				}),
			),
		})(t)
	})

	t.Run("Configuration Path", func(t *testing.T) {
		fileName := "new_config_path.yml"
		core.Test(&core.TestConfig{
			Commands:            GetCommands(),
			PromptResponseMocks: []string{},
			TmpHomeDir:          true,
			Cmd:                 cmdFromSettings("scw -c {{ .HOME }}/"+fileName+" init", defaultSettings),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					homeDir := ctx.OverrideEnv["HOME"]
					config, err := scw.LoadConfigFromPath(path.Join(homeDir, fileName))
					require.NoError(t, err)
					if config == nil {
						t.FailNow()
					}
					assert.Equal(t, secretKey, *config.SecretKey)
				},
			),
		})(t)
	})

	t.Run("Profile", func(t *testing.T) {
		t.Run("Named", func(t *testing.T) {
			profileName := "foobar"
			core.Test(&core.TestConfig{
				Commands:            GetCommands(),
				PromptResponseMocks: []string{},
				Cmd:                 cmdFromSettings("scw -p "+profileName+" init", defaultSettings),
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(0),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, secretKey, *config.Profiles[profileName].SecretKey)
					}),
				),
				TmpHomeDir: true,
			})(t)
		})
	})

	t.Run("CLIv2Config", func(t *testing.T) {
		dummySecretKey := "22222222-2222-2222-2222-222222222222"
		dummyAccessKey := "SCW22222222222222222"
		dummyConfig := &scw.Config{
			Profile: scw.Profile{
				AccessKey: &dummyAccessKey,
				SecretKey: &dummySecretKey,
			},
		}

		t.Run("NoOverwrite", func(t *testing.T) {
			core.Test(&core.TestConfig{
				Commands: GetCommands(),
				BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
					// Persist the dummy Config in the temp directory
					err := dummyConfig.SaveTo(path.Join(ctx.OverrideEnv["HOME"], ".config", "scw", "config.yaml"))
					if err != nil {
						t.FailNow()
					}
					return nil
				},
				Cmd: cmdFromSettings("scw init", defaultSettings),
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(1),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, dummySecretKey, *config.SecretKey)
						assert.Equal(t, dummyAccessKey, *config.AccessKey)
					}),
				),
				TmpHomeDir: true,
				PromptResponseMocks: []string{
					// Do you want to override the current config?
					"no",
				},
			})(t)
		})

		t.Run("Overwrite", func(t *testing.T) {
			core.Test(&core.TestConfig{
				Commands: GetCommands(),
				BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
					// Persist the dummy Config in the temp directory
					err := dummyConfig.Save()
					if err != nil {
						t.FailNow()
					}
					return nil
				},
				Cmd: cmdFromSettings("scw init", defaultSettings),
				Check: core.TestCheckCombine(
					core.TestCheckExitCode(0),
					core.TestCheckGolden(),
					checkConfig(func(t *testing.T, config *scw.Config) {
						assert.Equal(t, secretKey, *config.SecretKey)
						assert.Equal(t, organizationID, *config.DefaultOrganizationID)
					}),
				),
				TmpHomeDir: true,
				PromptResponseMocks: []string{
					// Do you want to override the current config?
					"yes",
				},
			})(t)
		})
	})
}
