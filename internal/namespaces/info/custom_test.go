package info

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
)

func Test_Info(t *testing.T) {
	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw info",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideEnv: map[string]string{
			"SCW_SECRET_KEY":              "22222222-2222-2222-2222-222222222222",
			"SCW_DEFAULT_ORGANIZATION_ID": "22222222-2222-2222-2222-222222222222",
			"SCW_DEFAULT_PROJECT_ID":      "22222222-2222-2222-2222-222222222222",
			"SCW_ACCESS_KEY":              "SCWYYYYYYYYYYYYYYYYY",
			"SCW_CONFIG_PATH":             "/tmp/.config/scw/config.yaml",
			"SCW_DEFAULT_REGION":          "fr-par",
			"SCW_DEFAULT_ZONE":            "fr-par-1",
		},
	}))

	t.Run("Show Secret", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw info show-secret=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideEnv: map[string]string{
			"SCW_SECRET_KEY":              "22222222-2222-2222-2222-222222222222",
			"SCW_DEFAULT_ORGANIZATION_ID": "22222222-2222-2222-2222-222222222222",
			"SCW_DEFAULT_PROJECT_ID":      "22222222-2222-2222-2222-222222222222",
			"SCW_ACCESS_KEY":              "SCWYYYYYYYYYYYYYYYYY",
			"SCW_CONFIG_PATH":             "/tmp/.config/scw/config.yaml",
			"SCW_DEFAULT_REGION":          "fr-par",
			"SCW_DEFAULT_ZONE":            "fr-par-1",
		},
	}))
}
