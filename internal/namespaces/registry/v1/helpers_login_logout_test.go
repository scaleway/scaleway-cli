package registry

import (
	"testing"

	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestGetRegistryEndpoint(t *testing.T) {
	tests := []struct {
		region   scw.Region
		expected string
	}{
		{region: scw.RegionFrPar, expected: "rg.fr-par.scw.cloud"},
		{region: scw.RegionNlAms, expected: "rg.nl-ams.scw.cloud"},
		{region: scw.RegionPlWaw, expected: "rg.pl-waw.scw.cloud"},
		{region: scw.RegionItMil, expected: "rg.it-mil.scw.eu"},
	}

	for _, test := range tests {
		t.Run(test.region.String(), func(t *testing.T) {
			assert.Equal(t, test.expected, getRegistryEndpoint(test.region))
		})
	}
}
