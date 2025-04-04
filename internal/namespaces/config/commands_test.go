package config_test

import (
	"os"
	"path"
	"regexp"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/config"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ConfigGetCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config get access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config get access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config get send-telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
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
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.Profiles["p1"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config set send-telemetry=true",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.True(t, *config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p test config set access-key=SCWNEWXXXXXXXXXXXXXX",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Equal(t, "SCWNEWXXXXXXXXXXXXXX", *config.Profiles["test"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigUnsetCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Nil(t, config.AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config unset access-key",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Nil(t, config.Profiles["p1"].AccessKey)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Telemetry", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config unset send-telemetry",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Nil(t, config.SendTelemetry)
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
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
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config profile delete p2",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			checkConfig(func(t *testing.T, config *scw.Config) {
				t.Helper()
				assert.Nil(t, config.Profiles["p2"])
			}),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
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
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config dump",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigDestroyCommand(t *testing.T) {
	path := "/tmp/test_config_destroy/"

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config destroy",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))

	t.Run("Check Config File", core.Test(&core.TestConfig{
		Commands: config.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			func(_ *core.BeforeFuncCtx) error {
				err := os.MkdirAll(path, os.ModePerm)
				if err != nil {
					t.Fatalf("MkdirAll %q: %s", path, err)
				}

				return nil
			},
			beforeFuncCreateFullConfig(),
			core.ExecStoreBeforeCmd(
				"Destroy",
				"scw config destroy",
			),
		),
		Cmd: "scw config dump",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		OverrideEnv: map[string]string{
			"HOME": path,
		},
		AfterFunc: func(_ *core.AfterFuncCtx) error {
			_ = os.RemoveAll(path)

			return nil
		},
	}))
}

func Test_ConfigInfoCommand(t *testing.T) {
	// replace ConfigPath lines with "/tmp/scw/.config/scw/config.yaml"
	configPathReplacements := []core.GoldenReplacement{
		{
			Pattern:     regexp.MustCompile(`(ConfigPath\s+).*`),
			Replacement: "$1/tmp/scw/.config/scw/config.yaml",
		},
		{
			Pattern:     regexp.MustCompile(`(?m)^(\s*"ConfigPath":\s*").*(",)`),
			Replacement: "$1/tmp/scw/.config/scw/config.yaml$2",
		},
	}

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config info",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGoldenAndReplacePatterns(configPathReplacements...),
		),
		TmpHomeDir: true,
	}))

	t.Run("Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p p1 config info",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGoldenAndReplacePatterns(configPathReplacements...),
		),
		TmpHomeDir: true,
	}))

	t.Run("Unknown Profile", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw -p test config info",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func Test_ConfigImportCommand(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		tmpFile, err := createTempConfigFile()
		if err != nil {
			t.Fatal(err)
		}
		defer os.Remove(tmpFile.Name())

		core.Test(&core.TestConfig{
			Commands:   config.GetCommands(),
			BeforeFunc: beforeFuncCreateFullConfig(),
			Cmd:        "scw config import " + tmpFile.Name(),
			Check: core.TestCheckCombine(
				core.TestCheckExitCode(0),
				core.TestCheckGolden(),
				checkConfig(func(t *testing.T, config *scw.Config) {
					t.Helper()
					// config
					assert.Equal(t, "22222222-2222-2222-2222-222222222222", *config.SecretKey)
					assert.Equal(t, "nl-ams", *config.DefaultRegion)
					// modified p1
					assert.Equal(
						t,
						"99999999-9999-9999-9999-999999999999",
						*config.Profiles["p1"].SecretKey,
					)
					assert.Equal(t, "nl-ams", *config.Profiles["p1"].DefaultRegion)
					// new p3
					assert.Equal(
						t,
						"33333333-3333-3333-3333-333333333333",
						*config.Profiles["p3"].SecretKey,
					)
					assert.Equal(t, "fr-par", *config.Profiles["p3"].DefaultRegion)
				}),
			),
			TmpHomeDir: true,
		})(t)
	})
}

func Test_ConfigValidateCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw config validate",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
	t.Run("Invalid default access key", core.Test(&core.TestConfig{
		Commands:   config.GetCommands(),
		BeforeFunc: beforeFuncCreateInvalidConfig(),
		Cmd:        "scw config validate",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
	t.Run("Invalid profile p1 secret key", core.Test(&core.TestConfig{
		Commands: config.GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			beforeFuncCreateInvalidConfig(),
			core.ExecBeforeCmd("scw config set access-key=SCWNEWXXXXXXXXXXXXXX"),
		),
		Cmd: "scw config validate",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
		TmpHomeDir: true,
	}))
}

func checkConfig(f func(t *testing.T, config *scw.Config)) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		t.Helper()
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
		err := os.MkdirAll(scwDir, 0o755)
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

func beforeFuncCreateInvalidConfig() core.BeforeFunc {
	return beforeFuncCreateConfigFile(&scw.Config{
		Profile: scw.Profile{
			AccessKey:             scw.StringPtr("invalidAccessKey"),
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
				SecretKey:             scw.StringPtr("invalidSecretKey"),
				APIURL:                scw.StringPtr("https://p1-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("11111111-1111-1111-1111-111111111111"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
		},
	})
}

func createTempConfigFile() (*os.File, error) {
	tmpFile, err := os.CreateTemp("", "tmp.yaml")
	if err != nil {
		return nil, err
	}

	configContent := `
access_key: SCWXXXXXXXXXXXXXXXXX
secret_key: 22222222-2222-2222-2222-222222222222
api_url: https://mock-api-url.com
insecure: true
default_organization_id: 22222222-2222-2222-2222-222222222222
default_region: nl-ams
default_zone: nl-ams-1
send_telemetry: true
profiles:
  p1:
    access_key: SCWP1XXXXXXXXXXXXXXX
    secret_key: 99999999-9999-9999-9999-999999999999
    api_url: https://p1-mock-api-url.com
    insecure: true
    default_organization_id: 99999999-9999-9999-9999-999999999999
    default_region: nl-ams
    default_zone: nl-ams-1
    send_telemetry: true
  p3:
    access_key: SCWP3XXXXXXXXXXXXXXX
    secret_key: 33333333-3333-3333-3333-333333333333
    api_url: https://p3-mock-api-url.com
    insecure: true
    default_organization_id: 33333333-3333-3333-3333-333333333333
    default_region: fr-par
    default_zone: fr-par-1
    send_telemetry: true
`

	if _, err := tmpFile.WriteString(configContent); err != nil {
		return nil, err
	}
	if err := tmpFile.Close(); err != nil {
		return nil, err
	}

	return tmpFile, nil
}
