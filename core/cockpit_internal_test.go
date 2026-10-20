package core

import (
	"testing"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

func TestExtractNameAndRegion(t *testing.T) {
	type regionResource struct { // e.g. mongodb.Instance, k8s.Cluster
		Name   string
		Region scw.Region
	}
	type zoneResource struct { // e.g. redis.Cluster
		Name string
		Zone scw.Zone
	}
	type customPool struct{ ID string }
	type wrappedResource struct { // k8s cluster get response wraps *Cluster
		*regionResource
		Pools []customPool
	}

	tests := []struct {
		name       string
		in         any
		wantName   string
		wantRegion string
		wantErr    bool
	}{
		{
			name:       "region field",
			in:         &regionResource{Name: "db", Region: scw.RegionFrPar},
			wantName:   "db",
			wantRegion: "fr-par",
		},
		{
			name:       "zone to region",
			in:         &zoneResource{Name: "cache", Zone: scw.ZoneFrPar1},
			wantName:   "cache",
			wantRegion: "fr-par",
		},
		{
			name:       "embedded pointer",
			in:         wrappedResource{&regionResource{Name: "k8s", Region: scw.RegionNlAms}, nil},
			wantName:   "k8s",
			wantRegion: "nl-ams",
		},
		{
			name:    "missing name",
			in:      &regionResource{Region: scw.RegionFrPar},
			wantErr: true,
		},
		{
			name:    "nil pointer",
			in:      (*regionResource)(nil),
			wantErr: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			name, region, err := extractNameAndRegion(tc.in)
			if tc.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}

				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if name != tc.wantName {
				t.Errorf("name = %q, want %q", name, tc.wantName)
			}
			if region != tc.wantRegion {
				t.Errorf("region = %q, want %q", region, tc.wantRegion)
			}
		})
	}
}
