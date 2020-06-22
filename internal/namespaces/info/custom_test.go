package info

import (
	"testing"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

func Test_Info(t *testing.T) {
	client, err := scw.NewClient(
		scw.WithAuth(
			"SCWXXXXXXXXXXXXXXXXX",
			"11111111-1111-1111-1111-111111111111",
		),
		scw.WithDefaultOrganizationID("11111111-1111-1111-1111-111111111111"),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	require.NoError(t, err)

	t.Run("Simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw info",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		Client: client,
	}))

	t.Run("Show Secret", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw info show-secret=true",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		OverrideEnv: map[string]string{
			"SCW_SECRET_KEY":         "22222222-2222-2222-2222-222222222222",
			"SCW_DEFAULT_PROJECT_ID": "22222222-2222-2222-2222-222222222222",
			"SCW_ACCESS_KEY":         "SCWYYYYYYYYYYYYYYYYY",
		},
		Client: client,
	}))
}
