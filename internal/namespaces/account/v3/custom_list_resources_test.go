package account_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
	function "github.com/scaleway/scaleway-cli/v2/internal/namespaces/function/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipam/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
)

func TestBuildQueryKey(t *testing.T) {
	t.Run("zone-based product returns product:zone", func(t *testing.T) {
		fetcher := fetch.WrapFetcher(&instance.FetchServers{})
		result := account.BuildQueryKey(fetcher, scw.ZoneFrPar1, "instances")
		assert.Equal(t, "instances:fr-par-1", result)
	})

	t.Run("region-based product returns product:region", func(t *testing.T) {
		fetcher := fetch.WrapFetcher(&rdb.FetchInstances{})
		result := account.BuildQueryKey(fetcher, scw.ZoneFrPar1, "rdb")
		assert.Equal(t, "rdb:fr-par", result)
	})

	t.Run("region-based product with zone fr-par-3", func(t *testing.T) {
		fetcher := fetch.WrapFetcher(&lb.FetchLoadBalancers{})
		result := account.BuildQueryKey(fetcher, scw.ZoneFrPar3, "lb")
		assert.Equal(t, "lb:fr-par", result)
	})
}

func TestFetcherLocalityTypes(t *testing.T) {
	tests := []struct {
		product  string
		fetcher  fetch.FetcherAny
		expected fetch.LocalityType
	}{
		{"rdb", fetch.WrapFetcher(&rdb.FetchInstances{}), fetch.LocalityTypeRegion},
		{"lb", fetch.WrapFetcher(&lb.FetchLoadBalancers{}), fetch.LocalityTypeRegion},
		{"k8s", fetch.WrapFetcher(&k8s.FetchClusters{}), fetch.LocalityTypeRegion},
		{"containers", fetch.WrapFetcher(&container.FetchNamespaces{}), fetch.LocalityTypeRegion},
		{"functions", fetch.WrapFetcher(&function.FetchNamespaces{}), fetch.LocalityTypeRegion},
		{"ipam", fetch.WrapFetcher(&ipam.FetchIPs{}), fetch.LocalityTypeRegion},
		{"buckets", fetch.WrapFetcher(&object.FetchBuckets{}), fetch.LocalityTypeRegion},
		{"secrets", fetch.WrapFetcher(&secret.FetchSecrets{}), fetch.LocalityTypeRegion},
		{"instances", fetch.WrapFetcher(&instance.FetchServers{}), fetch.LocalityTypeZone},
		{"instance-ips", fetch.WrapFetcher(&instance.FetchIPs{}), fetch.LocalityTypeZone},
		{"redis", fetch.WrapFetcher(&redis.FetchClusters{}), fetch.LocalityTypeZone},
		{"flexibleip", fetch.WrapFetcher(&flexibleip.FetchFlexibleIPs{}), fetch.LocalityTypeZone},
		{"block-volumes", fetch.WrapFetcher(&block.FetchVolumes{}), fetch.LocalityTypeZone},
		{"block-snapshots", fetch.WrapFetcher(&block.FetchSnapshots{}), fetch.LocalityTypeZone},
	}

	for _, tt := range tests {
		t.Run(tt.product, func(t *testing.T) {
			result := tt.fetcher.LocalityType()
			assert.Equal(
				t,
				tt.expected,
				result,
				"product %s should have locality type %v",
				tt.product,
				tt.expected,
			)
		})
	}
}

func TestResolveZones(t *testing.T) {
	t.Run("with requested zones", func(t *testing.T) {
		requested := []string{"fr-par-1", "nl-ams-1"}
		result := account.ResolveZones(requested)
		expected := []scw.Zone{scw.ZoneFrPar1, scw.ZoneNlAms1}
		assert.Equal(t, expected, result)
	})

	t.Run("with empty requested uses all zones", func(t *testing.T) {
		result := account.ResolveZones(nil)
		assert.NotEmpty(t, result)
		assert.Len(t, result, len(scw.AllZones))
	})
}

func TestResolveProducts(t *testing.T) {
	t.Run("with requested products", func(t *testing.T) {
		requested := []string{"instances", "rdb"}
		result := account.ResolveProducts(requested)
		assert.Equal(t, requested, result)
	})

	t.Run("with empty requested uses all products", func(t *testing.T) {
		result := account.ResolveProducts(nil)
		assert.NotEmpty(t, result)
		assert.Len(t, result, len(account.ProductFetchers))
	})
}

func TestSortResults(t *testing.T) {
	input := []fetch.ResourceResult{
		{Locality: "nl-ams", Product: "rdb", ID: "2"},
		{Locality: "fr-par", Product: "rdb", ID: "1"},
		{Locality: "fr-par", Product: "instances", ID: "1"},
		{Locality: "fr-par", Product: "rdb", ID: "2"},
	}

	account.SortResults(input)

	expected := []fetch.ResourceResult{
		{Locality: "fr-par", Product: "instances", ID: "1"},
		{Locality: "fr-par", Product: "rdb", ID: "1"},
		{Locality: "fr-par", Product: "rdb", ID: "2"},
		{Locality: "nl-ams", Product: "rdb", ID: "2"},
	}

	assert.Equal(t, expected, input)
}
