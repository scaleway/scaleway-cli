package instance_test

import (
	"context"
	"testing"

	instance "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instancesdk "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWarningServerTypeDeprecated(t *testing.T) {
	// Point the client at an unreachable endpoint so the catalog and
	// compatible-types lookups fail and the function returns early.
	client, err := scw.NewClient(
		scw.WithAuth("SCWXXXXXXXXXXXXXXXXX", "11111111-1111-1111-1111-111111111111"),
		scw.WithAPIURL("http://127.0.0.1:1"),
		scw.WithDefaultZone(scw.ZoneFrPar1),
		scw.WithDefaultRegion(scw.RegionFrPar),
	)
	require.NoError(t, err)

	server := &instancesdk.Server{CommercialType: "DEV1-S", Zone: scw.ZoneFrPar1}
	warnings := instance.WarningServerTypeDeprecated(context.Background(), client, server)

	require.NotEmpty(t, warnings)
	assert.Contains(t, warnings[0], "EndOfService")
}
