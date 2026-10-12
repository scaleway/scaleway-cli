package tui

import (
	"sync/atomic"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

// locality is the session's zone and explicit region. When region is
// empty it is derived from the zone at fetch time.
type locality struct {
	zone   scw.Zone
	region scw.Region
}

// Locality is a thread-safe holder for the session zone/region. It is
// shared between the App (which mutates it when the user picks a new
// zone or region) and the resource fetch closures (which read it in
// fetch goroutines). An atomic pointer swap keeps those reads
// race-free.
type Locality struct {
	ptr atomic.Pointer[locality]
}

// NewLocality creates a Locality with the given zone and no explicit
// region (derived from the zone).
func NewLocality(zone scw.Zone) *Locality {
	l := &Locality{}
	l.ptr.Store(&locality{zone: zone})

	return l
}

// Get returns the current zone and explicit region. A nil receiver
// yields empty values.
func (l *Locality) Get() (scw.Zone, scw.Region) {
	if l == nil {
		return "", ""
	}
	cur := l.ptr.Load()
	if cur == nil {
		return "", ""
	}

	return cur.zone, cur.region
}

// GetZone returns the current zone.
func (l *Locality) GetZone() scw.Zone {
	z, _ := l.Get()

	return z
}

// GetRegion returns the current explicit region ("" when derived).
func (l *Locality) GetRegion() scw.Region {
	_, r := l.Get()

	return r
}

// SetZone changes the session zone and clears any explicit region, so
// the region becomes derived from the new zone.
func (l *Locality) SetZone(zone scw.Zone) {
	l.ptr.Store(&locality{zone: zone})
}

// SetRegion changes the session region, keeping the current zone.
func (l *Locality) SetRegion(region scw.Region) {
	z, _ := l.Get()
	l.ptr.Store(&locality{zone: z, region: region})
}

// setLocalityMsg asks the app to change the session zone or region.
type setLocalityMsg struct {
	isZone bool
	zone   scw.Zone   // used when isZone
	region scw.Region // used when !isZone
}
