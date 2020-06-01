package config

import (
	"os"
	"path"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_ConfigGetCommand(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: beforeFuncCreateConfigFile(&scw.Config{
			Profile: scw.Profile{
				AccessKey: scw.StringPtr("mock-access-key"),
			},
		}),
		Cmd: "scw config get access_key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: beforeFuncCreateConfigFile(&scw.Config{
			Profile: scw.Profile{},
			Profiles: map[string]*scw.Profile{
				"test": {
					AccessKey: scw.StringPtr("mock-access-key"),
				},
			},
		}),
		Cmd: "scw -p test config get access_key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: beforeFuncCreateConfigFile(&scw.Config{
			SendTelemetry: scw.BoolPtr(true),
		}),
		Cmd: "scw -p test config get send_telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigSetCommand(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set access_key mock-access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, "mock-access-key", *config.AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config set access_key mock-access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, "mock-access-key", *config.Profiles["test"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set send_telemetry true",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, true, *config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigUnsetCommand(t *testing.T) {

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset access_key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config unset access_key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.Profiles["test"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset send_telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigDeleteProfileCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config profile delete p2",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.Profiles["p2"])
			}),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigDumpCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config dump",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func checkConfig(f func(t *testing.T, config *scw.Config)) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		homeDir := ctx.OverrideEnv["HOME"]
		config, err := scw.LoadConfigFromPath(path.Join(homeDir, ".config", "scw", "config.yaml"))
		require.NoError(t, err)
		f(t, config)
	}
}

func beforeFuncCreateConfigFile(c *scw.Config) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		homeDir := ctx.OverrideEnv["HOME"]
		scwDir := path.Join(homeDir, ".config", "scw")
		err := os.MkdirAll(scwDir, 0755)
		if err != nil {
			return err
		}

		return c.SaveTo(path.Join(scwDir, "config.yaml"))
	}
}

func beforeFuncCreateFullConfig() core.BeforeFunc {
	return beforeFuncCreateConfigFile(&scw.Config{
		Profile: scw.Profile{
			AccessKey:             scw.StringPtr("mock-access-key"),
			SecretKey:             scw.StringPtr("mock-secret-key"),
			APIURL:                scw.StringPtr("mock-api-url"),
			Insecure:              scw.BoolPtr(true),
			DefaultOrganizationID: scw.StringPtr("mock-orgaid"),
			DefaultRegion:         scw.StringPtr("fr-par"),
			DefaultZone:           scw.StringPtr("fr-par-1"),
		},
		Profiles: map[string]*scw.Profile{
			"p1": {
				AccessKey:             scw.StringPtr("p1-mock-access-key"),
				SecretKey:             scw.StringPtr("p1-mock-secret-key"),
				APIURL:                scw.StringPtr("p1-mock-api-url"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("p1-mock-orgaid"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
			},
			"p2": {
				AccessKey:             scw.StringPtr("p2-mock-access-key"),
				SecretKey:             scw.StringPtr("p2-mock-secret-key"),
				APIURL:                scw.StringPtr("p2-mock-api-url"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("p2-mock-orgaid"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
			},
		},
	})
}
