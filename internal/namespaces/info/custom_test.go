package info

import (
	"testing"

	"github.com/hashicorp/go-version"
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

	t.Run("Default", func(t *testing.T) {
		v, err := version.NewVersion("v2.0.0")
		require.NoError(t, err)
		t.Run("Simple", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw info",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			BuildInfo: core.BuildInfo{
				Version:   v,
				BuildDate: "a",
				GoVersion: "b",
				GitBranch: "c",
				GitCommit: "d",
				GoArch:    "e",
				GoOS:      "f",
			},
			Client: client,
		}))
	})
}
