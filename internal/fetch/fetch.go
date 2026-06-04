package fetch

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// LocalityType indicates whether a fetcher operates at zone or region level.
type LocalityType string

const (
	LocalityTypeZone   LocalityType = "zone"
	LocalityTypeRegion LocalityType = "region"
)

// Locality is a constraint interface for locality types.
// Only scw.Zone and scw.Region are valid implementations.
type Locality interface {
	scw.Zone | scw.Region
}

// Fetcher is a generic interface for fetching resources.
// The type parameter L must be either scw.Zone or scw.Region.
// Zone-based fetchers operate directly on zones.
// Region-based fetchers operate on regions.
// Fetchers should use core.ExtractClient(ctx) to obtain the client.
type Fetcher[L Locality] interface {
	Fetch(ctx context.Context, locality L, projectID string) ([]ResourceResult, error)
	Product() string
	Resource() string
	LocalityType() LocalityType
}

// FetcherAny is a non-generic interface used for storing heterogeneous
// fetchers in a common map or slice.
type FetcherAny interface {
	FetchAny(ctx context.Context, zone scw.Zone, projectID string) ([]ResourceResult, error)
	Product() string
	Resource() string
	LocalityType() LocalityType
}

// ZoneFetcher is a type alias for zone-based fetchers.
type ZoneFetcher = Fetcher[scw.Zone]

// RegionFetcher is a type alias for region-based fetchers.
type RegionFetcher = Fetcher[scw.Region]

// WrapFetcher wraps a generic Fetcher[L] into a FetcherAny.
// For region-based fetchers, it converts the zone to a region before calling Fetch.
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
	// Use the LocalityType to determine how to call Fetch
	switch w.fetcher.LocalityType() {
	case LocalityTypeZone:
		// Type assertion to zone fetcher
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
		// Type assertion to region fetcher
		regionFetcher, ok := any(w.fetcher).(Fetcher[scw.Region])
		if !ok {
			return nil, errors.New("failed to assert region fetcher")
		}

		return regionFetcher.Fetch(ctx, region, projectID)
	default:
		return nil, errors.New("unknown locality type")
	}
}

func (w *fetcherWrapper[L]) Product() string {
	return w.fetcher.Product()
}

func (w *fetcherWrapper[L]) Resource() string {
	return w.fetcher.Resource()
}

func (w *fetcherWrapper[L]) LocalityType() LocalityType {
	return w.fetcher.LocalityType()
}

// ResourceResult represents a single resource in the output.
// This is the non-generic version used by Fetcher.
type ResourceResult struct {
	Locality string `human:"locality" json:"locality"`
	Product  string `human:"product"  json:"product"`
	Resource string `human:"resource" json:"resource"`
	ID       string `human:"id"       json:"id"`
	Name     string `human:"name"     json:"name,omitempty"`
}

// ShouldIgnoreError checks if the error should be silently ignored.
// This includes:
// - Zone/region not available errors
// - HTTP 501 Not Implemented errors (service not available in certain zones)
func ShouldIgnoreError(err error) bool {
	if err == nil {
		return false
	}

	// Check for HTTP 501 errors
	if isHTTP501Error(err) {
		return true
	}

	// Check for zone/region unavailable errors
	errStr := err.Error()

	return strings.Contains(errStr, "not found") ||
		strings.Contains(errStr, "unavailable") ||
		strings.Contains(errStr, "not available") ||
		strings.Contains(errStr, "zone not found") ||
		strings.Contains(errStr, "region not found")
}

// isHTTP501Error checks if the error is an HTTP 501 Not Implemented error.
func isHTTP501Error(err error) bool {
	if responseErr, ok := errors.AsType[*scw.ResponseError](err); ok {
		return responseErr.StatusCode == http.StatusNotImplemented
	}

	return false
}
