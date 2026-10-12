package connect_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
)

func Test_Description(t *testing.T) {
	tests := []struct {
		name       string
		product    string
		location   string
		resourceID string
		username   string
		want       string
	}{
		{
			name:       "rdb",
			product:    "rdb",
			location:   "fr-par",
			resourceID: "11111111-1111-1111-1111-111111111111",
			username:   "foobar",
			want:       "Secret used to enable cli-connect. Connect directly using: scw rdb instance connect 11111111-1111-1111-1111-111111111111 username=foobar --region fr-par",
		},
		{
			name:       "mongodb",
			product:    "mongodb",
			location:   "fr-par",
			resourceID: "22222222-2222-2222-2222-222222222222",
			username:   "admin",
			want:       "Secret used to enable cli-connect. Connect directly using: scw mongodb instance connect 22222222-2222-2222-2222-222222222222 username=admin --region fr-par",
		},
		{
			name:       "rdb without username",
			product:    "rdb",
			location:   "nl-ams",
			resourceID: "33333333-3333-3333-3333-333333333333",
			username:   "",
			want:       "Secret used to enable cli-connect. Connect directly using: scw rdb instance connect 33333333-3333-3333-3333-333333333333 --region nl-ams",
		},
		{
			name:       "redis",
			product:    "redis",
			location:   "fr-par-1",
			resourceID: "44444444-4444-4444-4444-444444444444",
			username:   "redis",
			want:       "Secret used to enable cli-connect. Connect directly using: scw redis cluster connect 44444444-4444-4444-4444-444444444444 --zone fr-par-1",
		},
		{
			name:       "datawarehouse",
			product:    "datawarehouse",
			location:   "fr-par",
			resourceID: "55555555-5555-5555-5555-555555555555",
			username:   "scwadmin",
			want:       "Secret used to enable cli-connect. Connect directly using: scw datawarehouse deployment connect 55555555-5555-5555-5555-555555555555 --region fr-par",
		},
		{
			name:       "unknown product falls back to the generic description",
			product:    "mnq",
			location:   "fr-par",
			resourceID: "66666666-6666-6666-6666-666666666666",
			username:   "user",
			want:       "secret used by the cli to ease up connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := connect.Description(tt.product, tt.location, tt.resourceID, tt.username)
			if got != tt.want {
				t.Fatalf("Description() = %q, want %q", got, tt.want)
			}
		})
	}
}
