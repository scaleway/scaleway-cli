package fetch

import (
	"context"
	"errors"
	"net/http"
	"sort"
	"strings"
	"sync"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

type LocalityType string

const (
	LocalityTypeZone   LocalityType = "zone"
	LocalityTypeRegion LocalityType = "region"
)

type Locality interface {
	scw.Zone | scw.Region
}

// Fetcher is a generic interface for fetching resources.
// Fetchers implementations should use core.ExtractClient(ctx) to obtain the client.
type Fetcher[L Locality] interface {
	Fetch(ctx context.Context, locality L, projectID string) ([]ResourceResult, error)
	Namespace() string
	Resource() string
	LocalityType() LocalityType
}

// FetcherAny is a non-generic version of Fetcher.
type FetcherAny interface {
	FetchAny(ctx context.Context, zone scw.Zone, projectID string) ([]ResourceResult, error)
	Namespace() string
	Resource() string
	LocalityType() LocalityType
	ProductKey() string
}

type ZoneFetcher = Fetcher[scw.Zone]

type RegionFetcher = Fetcher[scw.Region]

// WrapFetcher wraps a generic Fetcher[L] into a FetcherAny.
func WrapFetcher[L Locality](f Fetcher[L]) FetcherAny {
	return &fetcherWrapper[L]{
		fetcher: f,
	}
}

type fetcherWrapper[L Locality] struct {
	fetcher Fetcher[L]
}

func (w *fetcherWrapper[L]) FetchAny(
	ctx context.Context,
	zone scw.Zone,
	projectID string,
) ([]ResourceResult, error) {
	switch w.fetcher.LocalityType() {
	case LocalityTypeZone:
		zoneFetcher, ok := any(w.fetcher).(Fetcher[scw.Zone])
		if !ok {
			return nil, errors.New("failed to assert zone fetcher")
		}

		return zoneFetcher.Fetch(ctx, zone, projectID)
	case LocalityTypeRegion:
		region, err := zone.Region()
		if err != nil {
			return nil, err
		}
		regionFetcher, ok := any(w.fetcher).(Fetcher[scw.Region])
		if !ok {
			return nil, errors.New("failed to assert region fetcher")
		}

		return regionFetcher.Fetch(ctx, region, projectID)
	default:
		return nil, errors.New("unknown locality type")
	}
}

func (w *fetcherWrapper[L]) Namespace() string {
	return w.fetcher.Namespace()
}

func (w *fetcherWrapper[L]) Resource() string {
	return w.fetcher.Resource()
}

func (w *fetcherWrapper[L]) LocalityType() LocalityType {
	return w.fetcher.LocalityType()
}

func (w *fetcherWrapper[L]) ProductKey() string {
	return w.fetcher.Namespace() + "-" + w.fetcher.Resource()
}

// ResourceResult represents a single resource in the output.
type ResourceResult struct {
	Locality string `human:"locality" json:"locality"`
	Product  string `human:"product"  json:"product"`
	Resource string `human:"resource" json:"resource"`
	ID       string `human:"id"       json:"id"`
	Name     string `human:"name"     json:"name,omitempty"`
}

// ShouldIgnoreError reports whether the error should be silently ignored
// (zone/region unavailable or HTTP 501).
func ShouldIgnoreError(err error) bool {
	if err == nil {
		return false
	}

	return isHTTP501Error(err) ||
		strings.Contains(err.Error(), "not found") ||
		strings.Contains(err.Error(), "unavailable") ||
		strings.Contains(err.Error(), "not available")
}

func isHTTP501Error(err error) bool {
	if responseErr, ok := errors.AsType[*scw.ResponseError](err); ok {
		return responseErr.StatusCode == http.StatusNotImplemented
	}

	return false
}

// factories maps product keys (Namespace() + "-" + Resource()) to fetcher
// factory functions, so fetchers are built lazily on demand.
var (
	fetchersMu sync.Mutex
	factories  = make(map[string]func() FetcherAny)
)

// RegisterFetcher registers a fetcher factory, keyed by Namespace() + "-" + Resource().
func RegisterFetcher[L Locality](factory func() Fetcher[L]) {
	wrapped := func() FetcherAny { return WrapFetcher(factory()) }
	RegisterFetcherAny(wrapped)
}

// RegisterFetcherAny registers a FetcherAny factory, keyed by Namespace() + "-" + Resource().
func RegisterFetcherAny(factory func() FetcherAny) {
	key := productKeyOf(factory)

	fetchersMu.Lock()
	defer fetchersMu.Unlock()
	if _, exists := factories[key]; exists {
		panic("fetch: duplicate fetcher registration for product key: " + key)
	}
	factories[key] = factory
}

// GetFetcher returns the fetcher registered for the given product key.
func GetFetcher(product string) (FetcherAny, bool) {
	fetchersMu.Lock()
	factory, ok := factories[product]
	fetchersMu.Unlock()
	if !ok {
		return nil, false
	}

	return factory(), true
}

// AllProductKeys returns the product keys of all registered fetchers.
func AllProductKeys() []string {
	fetchersMu.Lock()
	defer fetchersMu.Unlock()
	keys := make([]string, 0, len(factories))
	for k := range factories {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	return keys
}

func productKeyOf(factory func() FetcherAny) string {
	f := factory()

	return f.Namespace() + "-" + f.Resource()
}
