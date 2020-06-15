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
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config get access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config get access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config get send-telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p test config get access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigSetCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.Profiles["p1"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set send-telemetry=true",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, true, *config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p test config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.Profiles["test"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigUnsetCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset access-key",
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
		Cmd:        "scw -p p1 config unset access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.Profiles["p1"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset send-telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				assert.Nil(t, config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p test config unset access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
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

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config profile delete test",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
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
			AccessKey:             scw.StringPtr("SCWXXXXXXXXXXXXXXXXX"),
			SecretKey:             scw.StringPtr("11111111-1111-1111-1111-111111111111"),
			APIURL:                scw.StringPtr("https://mock-api-url.com"),
			Insecure:              scw.BoolPtr(true),
			DefaultOrganizationID: scw.StringPtr("11111111-1111-1111-1111-111111111111"),
			DefaultRegion:         scw.StringPtr("fr-par"),
			DefaultZone:           scw.StringPtr("fr-par-1"),
			SendTelemetry:         scw.BoolPtr(true),
		},
		Profiles: map[string]*scw.Profile{
			"p1": {
				AccessKey:             scw.StringPtr("SCWP1XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr("11111111-1111-1111-1111-111111111111"),
				APIURL:                scw.StringPtr("https://p1-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("11111111-1111-1111-1111-111111111111"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
			"p2": {
				AccessKey:             scw.StringPtr("SCWP2XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr("11111111-1111-1111-1111-111111111111"),
				APIURL:                scw.StringPtr("https://p2-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("11111111-1111-1111-1111-111111111111"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
		},
	})
}
