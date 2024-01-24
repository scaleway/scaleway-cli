package k8s

import (
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func Test_ExecCredential(t *testing.T) {
	// expect to return default secret_key
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	// expect to return 66666666-6666-6666-6666-666666666666
	t.Run("with scw_secret_key env", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwSecretKeyEnv: "66666666-6666-6666-6666-666666666666",
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	// expect to return p2 secret_key
	t.Run("with profile env", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "p2",
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	// expect to return p3 secret_key
	t.Run("with profile flag", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw --profile p3 k8s exec-credential",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))

	// expect to return p3 secret_key
	t.Run("with profile env and flag", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw --profile p3 k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "p2",
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
		),
	}))
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
			SecretKey:             scw.StringPtr("00000000-0000-0000-0000-111111111111"),
			APIURL:                scw.StringPtr("https://mock-api-url.com"),
			Insecure:              scw.BoolPtr(true),
			DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
			DefaultRegion:         scw.StringPtr("fr-par"),
			DefaultZone:           scw.StringPtr("fr-par-1"),
			SendTelemetry:         scw.BoolPtr(true),
		},
		Profiles: map[string]*scw.Profile{
			"p2": {
				AccessKey:             scw.StringPtr("SCWP2XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr("00000000-0000-0000-0000-222222222222"),
				APIURL:                scw.StringPtr("https://p2-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
			"p3": {
				AccessKey:             scw.StringPtr("SCWP3XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr("00000000-0000-0000-0000-333333333333"),
				APIURL:                scw.StringPtr("https://p3-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
		},
	})
}
