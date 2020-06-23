package init

import (
	"fmt"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func checkConfig(f func(t *testing.T, ctx *core.CheckFuncCtx, config *scw.Config)) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		homeDir := ctx.OverrideEnv["HOME"]
		config, err := scw.LoadConfigFromPath(path.Join(homeDir, ".config", "scw", "config.yaml"))
		require.NoError(t, err)
		f(t, ctx, config)
	}
}

func appendArgs(prefix string, settings map[string]string) string {
	res := prefix
	for k, v := range settings {
		res += fmt.Sprintf(" %s=%s", k, v)
	}
	return res
}

func beforeFuncSaveConfig(config *scw.Config) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		// Persist the dummy Config in the temp directory
		return config.SaveTo(path.Join(ctx.OverrideEnv["HOME"], ".config", "scw", "config.yaml"))
	}
}

func TestInit(t *testing.T) {
	defaultArgs := map[string]string{
		"secret-key":           "{{ .SecretKey }}",
		"organization-id":      "{{ .OrganizationID }}",
		"send-telemetry":       "true",
		"install-autocomplete": "false",
		"remove-v1-config":     "false",
		"with-ssh-key":         "false",
	}

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:            GetCommands(),
		BeforeFunc:          baseBeforeFunc(),
		PromptResponseMocks: []string{},
		TmpHomeDir:          true,
		Cmd:                 appendArgs("scw init", defaultArgs),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, ctx *core.CheckFuncCtx, config *scw.Config) {
				secretKey, _ := ctx.Client.GetSecretKey()
				assert.Equal(t, secretKey, *config.SecretKey)
			}),
		),
	}))

	t.Run("Configuration Path", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			baseBeforeFunc(),
			func(ctx *core.BeforeFuncCtx) error {
				ctx.Meta["CONFIG_PATH"] = path.Join(ctx.Meta["HOME"].(string), "new_config_path.yml")
				return nil
			},
		),
		PromptResponseMocks: []string{},
		TmpHomeDir:          true,
		Cmd:                 appendArgs("scw -c {{ .CONFIG_PATH }} init", defaultArgs),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				config, err := scw.LoadConfigFromPath(ctx.Meta["CONFIG_PATH"].(string))
				require.NoError(t, err)
				secretKey, _ := ctx.Client.GetSecretKey()
				assert.Equal(t, secretKey, *config.SecretKey)
			},
		),
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:            GetCommands(),
		BeforeFunc:          baseBeforeFunc(),
		PromptResponseMocks: []string{},
		Cmd:                 appendArgs("scw -p foobar init", defaultArgs),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, ctx *core.CheckFuncCtx, config *scw.Config) {
				secretKey, _ := ctx.Client.GetSecretKey()
				assert.Equal(t, secretKey, *config.Profiles["foobar"].SecretKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("CLIv2Config", func(t *testing.T) {
		dummySecretKey := "22222222-2222-2222-2222-222222222222"
		dummyAccessKey := "SCW22222222222222222"
		dummyConfig := &scw.Config{
			Profile: scw.Profile{
				AccessKey: &dummyAccessKey,
				SecretKey: &dummySecretKey,
			},
			Profiles: map[string]*scw.Profile{
				"test": {
					AccessKey: &dummyAccessKey,
					SecretKey: &dummySecretKey,
				},
			},
		}

		t.Run("NoOverwrite", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				baseBeforeFunc(),
				beforeFuncSaveConfig(dummyConfig),
			),
			Cmd: appendArgs("scw init", defaultArgs),
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				checkConfig(func(t *testing.T, ctx *core.CheckFuncCtx, config *scw.Config) {
					assert.Equal(t, dummyConfig.String(), config.String())
				}),
			),
			TmpHomeDir: true,
			PromptResponseMocks: []string{
				// Do you want to override the current config?
				"no",
			},
		}))

		t.Run("Overwrite", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			BeforeFunc: core.BeforeFuncCombine(
				baseBeforeFunc(),
				beforeFuncSaveConfig(dummyConfig),
			),
			Cmd: appendArgs("scw init", defaultArgs),
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				checkConfig(func(t *testing.T, ctx *core.CheckFuncCtx, config *scw.Config) {
					secretKey, _ := ctx.Client.GetSecretKey()
					organizationID, _ := ctx.Client.GetDefaultOrganizationID()
					assert.Equal(t, secretKey, *config.SecretKey)
					assert.Equal(t, organizationID, *config.DefaultOrganizationID)
				}),
			),
			TmpHomeDir: true,
			PromptResponseMocks: []string{
				// Do you want to override the current config?
				"yes",
			},
		}))
	})
}
