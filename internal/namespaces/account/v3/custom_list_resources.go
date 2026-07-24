package account

import (
	"context"
	"reflect"
	"sort"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/fetch"
	applesilicon "github.com/scaleway/scaleway-cli/v2/internal/namespaces/applesilicon/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/baremetal/v1"
	block "github.com/scaleway/scaleway-cli/v2/internal/namespaces/block/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/cockpit/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/container/v1"
	file "github.com/scaleway/scaleway-cli/v2/internal/namespaces/file/v1alpha1"
	flexibleip "github.com/scaleway/scaleway-cli/v2/internal/namespaces/flexibleip/v1alpha1"
	function "github.com/scaleway/scaleway-cli/v2/internal/namespaces/function/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/inference/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/ipam/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/k8s/v1"
	key_manager "github.com/scaleway/scaleway-cli/v2/internal/namespaces/key_manager/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/lb/v1"
	mongodb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1alpha1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/object/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/rdb/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/redis/v1"
	registry "github.com/scaleway/scaleway-cli/v2/internal/namespaces/registry/v1"
	s2s_vpn "github.com/scaleway/scaleway-cli/v2/internal/namespaces/s2s_vpn/v1alpha1"
	searchdb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/searchdb/v1alpha1"
	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpcgw/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/webhosting/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// ListRequest defines the input parameters for listing all resources.
type ListRequest struct {
	ProjectID *string  `arg:"project-id" json:"project_id"`
	Zones     []string `arg:"zones"      json:"zones"`
	Products  []string `arg:"products"   json:"products"`
}

// ProductFetchers maps product names to their fetcher instances.
var ProductFetchers = map[string]fetch.FetcherAny{
	"baremetal-servers":        fetch.WrapFetcher(&baremetal.FetchServers{}),
	"applesilicon-servers":     fetch.WrapFetcher(&applesilicon.FetchServers{}),
	"instance-servers":         fetch.WrapFetcher(&instance.FetchServers{}),
	"instance-ips":             fetch.WrapFetcher(&instance.FetchIPs{}),
	"instance-volumes":         fetch.WrapFetcher(&instance.FetchVolumes{}),
	"instance-snapshots":       fetch.WrapFetcher(&instance.FetchSnapshots{}),
	"instance-security-groups": fetch.WrapFetcher(&instance.FetchSecurityGroups{}),
	"ipam-ip":                  fetch.WrapFetcher(&ipam.FetchIPs{}),
	"block-volumes":            fetch.WrapFetcher(&block.FetchVolumes{}),
	"block-snapshots":          fetch.WrapFetcher(&block.FetchSnapshots{}),
	"object-buckets":           fetch.WrapFetcher(&object.FetchBuckets{}),
	"rdb-instances":            fetch.WrapFetcher(&rdb.FetchInstances{}),
	"redis-cluster":            fetch.WrapFetcher(&redis.FetchClusters{}),
	"lb-lb":                    fetch.WrapFetcher(&lb.FetchLoadBalancers{}),
	"k8s-cluster":              fetch.WrapFetcher(&k8s.FetchClusters{}),
	"container-namespaces":     fetch.WrapFetcher(&container.FetchNamespaces{}),
	"function-namespaces":      fetch.WrapFetcher(&function.FetchNamespaces{}),
	"fip-ip":                   fetch.WrapFetcher(&flexibleip.FetchFlexibleIPs{}),
	"secret-secrets":           fetch.WrapFetcher(&secret.FetchSecrets{}),
	"vpc-vpc":                  fetch.WrapFetcher(&vpc.FetchVPCs{}),
	"file-filesystem":          fetch.WrapFetcher(&file.FetchFileSystem{}),
	"webhosting-hosting":       fetch.WrapFetcher(&webhosting.FetchHostings{}),
	"vpc-gw-vpc-gw":            fetch.WrapFetcher(&vpcgw.FetchGateways{}),
	"vpc-gw-ip":                fetch.WrapFetcher(&vpcgw.FetchIPs{}),
	"mongodb-instance":         fetch.WrapFetcher(&mongodb.FetchInstances{}),
	"mongodb-snapshot":         fetch.WrapFetcher(&mongodb.FetchSnapshot{}),
	"key-manager":              fetch.WrapFetcher(&key_manager.FetchKey{}),
	"inference-deployment":     fetch.WrapFetcher(&inference.FetchDeployment{}),
	"cockpit-token":            fetch.WrapFetcher(&cockpit.FetchToken{}),
	"cockpit-datasource":       fetch.WrapFetcher(&cockpit.FetchDataSource{}),
	"registry-namespaces":      fetch.WrapFetcher(&registry.FetchNamespaces{}),
	"searchdb-deployments":     fetch.WrapFetcher(&searchdb.FetchDeployments{}),
	"s2s-vpn-vpn-gateway":      fetch.WrapFetcher(&s2s_vpn.FetchVpnGateways{}),
}

// listResources returns the command for listing all resources.
func listResources() *core.Command {
	return &core.Command{
		Short:     `List all resources across all zones`,
		Long:      `List all resources across all zones and products. Results are grouped by locality and product. Errors are aggregated and not fail-fast.`,
		Namespace: "account",
		Resource:  "project",
		Verb:      "list-resources",
		ArgsType:  reflect.TypeOf(ListRequest{}),
		ArgSpecs: core.ArgSpecs{
			{
				Name:         "project-id",
				Short:        `Filter by project ID. If none is passed the default project ID will be used`,
				Required:     false,
				Deprecated:   false,
				Positional:   true,
				ValidateFunc: core.ValidateProjectID(),
				Default: func(ctx context.Context) (value string, doc string) {
					client := core.ExtractClient(ctx)
					projectID, _ := client.GetDefaultProjectID()

					return projectID, projectID
				},
			},
			{
				Name:       "zones",
				Short:      `Filter by zones (comma-separated, e.g. fr-par-1,nl-ams-1). If empty, all zones are queried`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
			{
				Name:       "products",
				Short:      `Filter by products (comma-separated: instances,instance-ips,instance-volumes,instance-snapshots,instance-security-groups,ipam,block-volumes,block-snapshots,buckets,rdb,redis,lb,k8s,containers,functions,flexibleip,secrets,vpc,file,webhosting,registry,searchdb,s2s-vpn-vpn-gateway). If empty, all products are queried`,
				Required:   false,
				Deprecated: false,
				Positional: false,
			},
		},
		Run: runListResources,
		View: &core.View{
			Fields: []*core.ViewField{
				{
					FieldName: "Locality",
					Label:     "Locality",
				},
				{
					FieldName: "Product",
					Label:     "Product",
				},
				{
					FieldName: "Resource",
					Label:     "Resource",
				},
				{
					FieldName: "ID",
					Label:     "ID",
				},
				{
					FieldName: "Name",
					Label:     "Name",
				},
			},
		},
	}
}

func runListResources(ctx context.Context, argsI any) (any, error) {
	request := argsI.(*ListRequest)

	// Determine zones to query
	zones := ResolveZones(request.Zones)
	if len(zones) == 0 {
		return []fetch.ResourceResult{}, nil
	}

	// Determine products to query
	products := ResolveProducts(request.Products)
	if len(products) == 0 {
		return []fetch.ResourceResult{}, nil
	}

	// Result aggregation
	var allResults []fetch.ResourceResult
	var resultsMu sync.Mutex

	// WaitGroup for parallel execution
	var wg sync.WaitGroup

	// Semaphore to limit concurrency (50 concurrent requests max)
	// Chosen to balance API rate limits with throughput
	sem := make(chan struct{}, 50)

	// Track which localities have been queried per product to avoid duplicates
	queried := make(map[string]bool)
	var mu sync.Mutex

	for _, zone := range zones {
		for _, product := range products {
			wg.Add(1)
			go func(zone scw.Zone, product string) {
				defer wg.Done()
				sem <- struct{}{}
				defer func() { <-sem }()

				fetcher, ok := ProductFetchers[product]
				if !ok {
					return
				}

				// Build the query key based on the fetcher's locality type
				queryKey := BuildQueryKey(fetcher, zone, product)

				// Check if we already queried this locality for this product
				mu.Lock()
				if queried[queryKey] {
					mu.Unlock()

					return // Already queried this locality for this product
				}
				queried[queryKey] = true
				mu.Unlock()

				// Call Fetch method on the fetcher interface with project filter
				resources, err := fetcher.FetchAny(ctx, zone, *request.ProjectID)
				if err != nil {
					// Log error in debug mode for troubleshooting
					core.ExtractLogger(ctx).
						Debugf("error fetching %s in %s: %v", product, zone, err)

					return
				}

				if len(resources) == 0 {
					return
				}

				// Populate Product and Resource fields from the fetcher
				productName := fetcher.Product()
				resourceName := fetcher.Resource()
				for i := range resources {
					resources[i].Product = productName
					resources[i].Resource = resourceName
				}

				resultsMu.Lock()
				defer resultsMu.Unlock()

				allResults = append(allResults, resources...)
			}(zone, product)
		}
	}

	wg.Wait()

	// Sort results for consistent output
	SortResults(allResults)

	return allResults, nil
}

// SortResults sorts results by locality, then product, then resource, then ID for consistent output
func SortResults(results []fetch.ResourceResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].Locality != results[j].Locality {
			return results[i].Locality < results[j].Locality
		}
		if results[i].Product != results[j].Product {
			return results[i].Product < results[j].Product
		}
		if results[i].Resource != results[j].Resource {
			return results[i].Resource < results[j].Resource
		}

		return results[i].ID < results[j].ID
	})
}

// BuildQueryKey builds a unique key for deduplication based on the fetcher's locality type.
// For zone-based fetchers: "product:zone" (e.g., "instance:fr-par-1")
// For region-based fetchers: "product:region" (e.g., "vpc:fr-par")
func BuildQueryKey(fetcher fetch.FetcherAny, zone scw.Zone, product string) string {
	switch fetcher.LocalityType() {
	case fetch.LocalityTypeRegion:
		region, err := zone.Region()
		if err != nil {
			return product + ":" + zone.String()
		}

		return product + ":" + region.String()
	default:
		return product + ":" + zone.String()
	}
}

// ResolveZones returns the list of zones to query.
func ResolveZones(requested []string) []scw.Zone {
	if len(requested) > 0 {
		zones := make([]scw.Zone, 0, len(requested))
		for _, z := range requested {
			zones = append(zones, scw.Zone(z))
		}

		return zones
	}

	// Return all known zones
	return scw.AllZones
}

// ResolveProducts returns the list of products to query.
func ResolveProducts(requested []string) []string {
	if len(requested) > 0 {
		return requested
	}

	// Return all known products
	products := make([]string, 0, len(ProductFetchers))
	for product := range ProductFetchers {
		products = append(products, product)
	}

	return products
}
