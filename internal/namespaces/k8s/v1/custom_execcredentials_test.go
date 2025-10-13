package k8s_test

import (
	"encoding/json"
	"os"
	"path"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

const (
	p1Secret  = "00000000-0000-0000-0000-111111111111"
	p2Secret  = "00000000-0000-0000-0000-222222222222"
	p3Secret  = "00000000-0000-0000-0000-333333333333"
	envSecret = "66666666-6666-6666-6666-666666666666"
)

func Test_ExecCredential(t *testing.T) {
	////
	// Simple expect to return current secret_key
	////
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Ignore keys in test env
			scw.ScwSecretKeyEnv: "", // Ignore keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			assertTokenInResponse(p1Secret),
		),
	}))

	////
	// expect to return 66666666-6666-6666-6666-666666666666
	////
	t.Run("with scw_secret_key env", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwSecretKeyEnv: envSecret,
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			assertTokenInResponse(envSecret),
		),
	}))

	////
	// expect to return p2 secret_key
	////
	t.Run("with profile env", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "p2",
			scw.ScwAccessKeyEnv:     "", // Ignore keys in test env
			scw.ScwSecretKeyEnv:     "", // Ignore keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			assertTokenInResponse(p2Secret),
		),
	}))

	////
	// expect to return p3 secret_key
	////
	t.Run("with profile flag", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw --profile p3 k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwAccessKeyEnv: "", // Ignore keys in test env
			scw.ScwSecretKeyEnv: "", // Ignore keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			assertTokenInResponse(p3Secret),
		),
	}))

	////
	// expect to return p3 secret_key
	////
	t.Run("with profile env and flag", core.Test(&core.TestConfig{
		Commands:   k8s.GetCommands(),
		TmpHomeDir: true,
		BeforeFunc: beforeFuncCreateFullConfig(),
		Cmd:        "scw --profile p3 k8s exec-credential",
		OverrideEnv: map[string]string{
			scw.ScwActiveProfileEnv: "p2",
			scw.ScwAccessKeyEnv:     "", // Ignore keys in test env
			scw.ScwSecretKeyEnv:     "", // Ignore keys in test env
		},
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			core.TestCheckGolden(),
			assertTokenInResponse(p3Secret),
		),
	}))
}

func beforeFuncCreateConfigFile(c *scw.Config) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		homeDir := ctx.OverrideEnv["HOME"]
		scwDir := path.Join(homeDir, ".config", "scw")
		if err := os.MkdirAll(scwDir, 0o0755); err != nil {
			return err
		}

		scwPath := path.Join(scwDir, "config.yaml")
		if err := c.SaveTo(scwPath); err != nil {
			return err
		}

		return nil
	}
}

func beforeFuncCreateFullConfig() core.BeforeFunc {
	return beforeFuncCreateConfigFile(&scw.Config{
		Profile: scw.Profile{
			AccessKey:             scw.StringPtr("SCWXXXXXXXXXXXXXXXXX"),
			SecretKey:             scw.StringPtr(p1Secret),
			APIURL:                scw.StringPtr("https://mock-api-url.com"),
			Insecure:              scw.BoolPtr(true),
			DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
			DefaultProjectID:      scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
			DefaultRegion:         scw.StringPtr("fr-par"),
			DefaultZone:           scw.StringPtr("fr-par-1"),
			SendTelemetry:         scw.BoolPtr(true),
		},
		Profiles: map[string]*scw.Profile{
			"p2": {
				AccessKey:             scw.StringPtr("SCWP2XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr(p2Secret),
				APIURL:                scw.StringPtr("https://p2-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultProjectID:      scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
			"p3": {
				AccessKey:             scw.StringPtr("SCWP3XXXXXXXXXXXXXXX"),
				SecretKey:             scw.StringPtr(p3Secret),
				APIURL:                scw.StringPtr("https://p3-mock-api-url.com"),
				Insecure:              scw.BoolPtr(true),
				DefaultOrganizationID: scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultProjectID:      scw.StringPtr("deadbeef-dead-dead-dead-deaddeafbeef"),
				DefaultRegion:         scw.StringPtr("fr-par"),
				DefaultZone:           scw.StringPtr("fr-par-1"),
				SendTelemetry:         scw.BoolPtr(true),
			},
		},
	})
}

func assertTokenInResponse(expectedToken string) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		t.Helper()
		res := ctx.Result.(string)
		creds := k8s.ExecCredential{}
		err := json.Unmarshal([]byte(res), &creds)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, expectedToken, creds.Status.Token)
	}
}
