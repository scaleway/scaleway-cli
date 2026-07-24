package account_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/commands"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	account "github.com/scaleway/scaleway-cli/v2/internal/namespaces/account/v3"
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

// ExceptionProducts lists the namespaces that do not have fetchers
// and therefore must NOT be in ProductFetchers.
//
// These are typically utility namespaces that provide CLI functionality
// but do not manage actual cloud resources.
var ExceptionProducts = map[string]struct{}{
	"account":                 {},
	"alias":                   {},
	"autocomplete":            {},
	"config":                  {},
	"datalab":                 {},
	"datawarehouse":           {},
	"feedback":                {},
	"help":                    {},
	"iam":                     {},
	"info":                    {},
	"init":                    {},
	"login":                   {},
	"marketplace":             {},
	"mcp":                     {},
	"partner":                 {},
	"product-catalog":         {},
	"search":                  {},
	"shell":                   {},
	"version":                 {},
	"billing":                 {},
	"audit-trail":             {},
	"domain":                  {},
	"environmental-footprint": {},
	"interlink":               {},
	"iot":                     {},
	"jobs":                    {},
	"kafka":                   {},
	"mnq":                     {},
	"sdb-sql":                 {},
	"tem":                     {},
	"edge-services":           {},
	// Beta namespaces
	"dedibox": {},
}

// getNamespaceToProductMap returns the mapping between CLI namespaces/resources
// and product names in ProductFetchers.
//
// This map is needed because some products have different naming between
// their CLI namespace and their ProductFetchers key (e.g., "baremetal" -> "baremetal-servers").
func getNamespaceToProductMap() map[string]string {
	return map[string]string{
		"baremetal":                "baremetal-servers",
		"apple-silicon":            "applesilicon-servers",
		"instance":                 "instance-servers",
		"instance-ips":             "instance-ips",
		"instance-volumes":         "instance-volumes",
		"instance-snapshots":       "instance-snapshots",
		"instance-security-groups": "instance-security-groups",
		"ipam":                     "ipam-ip",
		"block":                    "block-volumes",
		"block-snapshots":          "block-snapshots",
		"object":                   "object-buckets",
		"rdb":                      "rdb-instances",
		"redis":                    "redis-cluster",
		"lb":                       "lb-lb",
		"k8s":                      "k8s-cluster",
		"container":                "container-namespaces",
		"function":                 "function-namespaces",
		"fip":                      "fip-ip",
		"secret":                   "secret-secrets",
		"vpc":                      "vpc-vpc",
		"file":                     "file-filesystem",
		"webhosting":               "webhosting-hosting",
		"vpc-gw":                   "vpc-gw-vpc-gw",
		"vpc-gw-ip":                "vpc-gw-ip",
		"mongodb":                  "mongodb-instance",
		"mongodb-snapshot":         "mongodb-snapshot",
		"keymanager":               "key-manager",
		"inference":                "inference-deployment",
		"cockpit":                  "cockpit-token",
		"cockpit-datasource":       "cockpit-datasource",
		"registry":                 "registry-namespaces",
		"searchdb":                 "searchdb-deployments",
		"s2s-vpn":                  "s2s-vpn-vpn-gateway",
	}
}

// computeExpectedProductsFromCommands dynamically computes the expected products
// by iterating through all registered commands and mapping them to ProductFetchers.
//
// This eliminates the need to maintain a static list of expected products.
// The function:
// 1. Iterates through all registered commands
// 2. Skips hidden commands and exceptions
// 3. Maps namespace/resource to ProductFetchers keys
// 4. Adds extra products that share namespaces (e.g., instance-ips, block-snapshots)
func computeExpectedProductsFromCommands() map[string]struct{} {
	allCommands := commands.GetCommands()
	expectedProducts := make(map[string]struct{})
	namespaceToProduct := getNamespaceToProductMap()

	// Iterate through all commands to find available namespaces/resources
	seen := make(map[string]struct{})
	for _, cmd := range allCommands.GetAll() {
		if cmd.Namespace == "" || cmd.Hidden {
			continue
		}

		// Check if this namespace is an exception
		if _, isException := ExceptionProducts[cmd.Namespace]; isException {
			continue
		}

		// Build a namespace or namespace/resource key
		key := cmd.Namespace
		if cmd.Resource != "" {
			key = cmd.Namespace + "-" + cmd.Resource
		}

		if _, exists := seen[key]; exists {
			continue
		}
		seen[key] = struct{}{}

		// Check if this namespace/resource has a corresponding product
		if product, ok := namespaceToProduct[key]; ok {
			expectedProducts[product] = struct{}{}
		}
	}

	// Manually add products that are not separate resources
	// but have dedicated fetchers (e.g., instance-ips, block-snapshots, etc.)
	// These products share the same namespace as other resources.
	extraProducts := []string{
		"instance-ips",
		"instance-volumes",
		"instance-snapshots",
		"instance-security-groups",
		"block-snapshots",
		"vpc-gw-ip",
		"cockpit-datasource",
		"mongodb-snapshot",
	}
	for _, p := range extraProducts {
		expectedProducts[p] = struct{}{}
	}

	return expectedProducts
}

func TestProductFetchersCompleteness(t *testing.T) {
	expectedProducts := computeExpectedProductsFromCommands()

	t.Run("all_expected_products_are_in_ProductFetchers", func(t *testing.T) {
		for expectedProduct := range expectedProducts {
			fetcher, exists := account.ProductFetchers[expectedProduct]
			assert.True(
				t,
				exists,
				"product %q is expected but missing from ProductFetchers - add it to ProductFetchers in custom_list_resources.go",
				expectedProduct,
			)
			assert.NotNil(
				t,
				fetcher,
				"product %q has a nil fetcher in ProductFetchers",
				expectedProduct,
			)
		}
	})

	t.Run("no_unexpected_products_in_ProductFetchers", func(t *testing.T) {
		for product := range account.ProductFetchers {
			_, expected := expectedProducts[product]
			assert.True(
				t,
				expected,
				"product %q is in ProductFetchers but not in ExpectedProductsFromCommands - add it to getNamespaceToProductMap() or verify it should be removed",
				product,
			)
		}
	})

	t.Run("ProductFetchers_has_exactly_expected_count", func(t *testing.T) {
		expectedCount := len(expectedProducts)
		actualCount := len(account.ProductFetchers)
		assert.Equal(
			t,
			expectedCount,
			actualCount,
			"ProductFetchers count mismatch - expected %d, got %d. Check for missing or extra products.",
			expectedCount,
			actualCount,
		)
	})
}
