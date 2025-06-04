package registry_test

import (
	"bytes"
	"os"
	"path"
	"runtime"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRegistryInstallDockerHelperCommand(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("windows is not supported")
	}

	t.Run("Simple", core.Test(&core.TestConfig{
		BeforeFunc: nil,
		Commands:   registry.GetCommands(),
		Cmd:        "scw registry install-docker-helper path={{ .HOME }}",
		Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
			t.Helper()
			scriptPath := path.Join(ctx.Meta["HOME"].(string), "docker-credential-scw")
			scriptContent, err := os.ReadFile(scriptPath)
			require.NoError(t, err)
			assert.Equal(t, "#!/bin/sh\nscw registry docker-helper \"$@\"\n", string(scriptContent))
			stats, err := os.Stat(scriptPath)
			require.NoError(t, err)
			assert.Equal(t, os.FileMode(0o755), stats.Mode())

			dockerConfigPath := path.Join(ctx.Meta["HOME"].(string), ".docker", "config.json")
			dockerConfigContent, err := os.ReadFile(dockerConfigPath)
			require.NoError(t, err)
			assert.JSONEq(
				t,
				"{\n  \"credHelpers\": {\n    \"rg.fr-par.scw.cloud\": \"scw\",\n    \"rg.nl-ams.scw.cloud\": \"scw\",\n    \"rg.pl-waw.scw.cloud\": \"scw\"\n  }\n}\n",
				string(dockerConfigContent),
			)
		},
		AfterFunc:   nil,
		TmpHomeDir:  true,
		OverrideEnv: nil,
		PromptResponseMocks: []string{
			"yes",
		},
	}))

	t.Run("With profile", core.Test(&core.TestConfig{
		BeforeFunc: nil,
		Commands:   registry.GetCommands(),
		Cmd:        "scw -p profile01 registry install-docker-helper path={{ .HOME }}",
		Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
			t.Helper()
			scriptPath := path.Join(ctx.Meta["HOME"].(string), "docker-credential-scw")
			scriptContent, err := os.ReadFile(scriptPath)
			require.NoError(t, err)
			assert.Equal(
				t,
				"#!/bin/sh\nPROFILE_NAME=\"profile01\"\nif [[ ! -z \"$SCW_PROFILE\" ]]\nthen \n\tPROFILE_NAME=\"$SCW_PROFILE\"\nfi\nscw --profile $PROFILE_NAME registry docker-helper \"$@\"\n",
				string(scriptContent),
			)
		},
		AfterFunc:   nil,
		TmpHomeDir:  true,
		OverrideEnv: nil,
		PromptResponseMocks: []string{
			"yes",
		},
	}))
}

func TestRegistryDockerHelperGetCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		Cmd:      "scw registry docker-helper get",
		Stdin:    bytes.NewBufferString("rg.fr-par.scw.cloud\n"),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: nil,
	}))
}

func TestRegistryDockerHelperListCommand(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: registry.GetCommands(),
		Cmd:      "scw registry docker-helper list",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: nil,
	}))
}
