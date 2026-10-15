package connect_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/connect"
)

func Test_Resource_Description(t *testing.T) {
	tests := []struct {
		name string
		res  connect.Resource
		want string
	}{
		{
			name: "rdb",
			res: connect.Resource{
				Product:  "rdb",
				Location: "fr-par",
				ID:       "11111111-1111-1111-1111-111111111111",
				Username: "foobar",
			},
			want: "Secret used to enable cli-connect. Connect directly using: scw rdb instance connect 11111111-1111-1111-1111-111111111111 username=foobar --region fr-par",
		},
		{
			name: "mongodb",
			res: connect.Resource{
				Product:  "mongodb",
				Location: "fr-par",
				ID:       "22222222-2222-2222-2222-222222222222",
				Username: "admin",
			},
			want: "Secret used to enable cli-connect. Connect directly using: scw mongodb instance connect 22222222-2222-2222-2222-222222222222 username=admin --region fr-par",
		},
		{
			name: "rdb without username",
			res: connect.Resource{
				Product:  "rdb",
				Location: "nl-ams",
				ID:       "33333333-3333-3333-3333-333333333333",
			},
			want: "Secret used to enable cli-connect. Connect directly using: scw rdb instance connect 33333333-3333-3333-3333-333333333333 --region nl-ams",
		},
		{
			name: "redis",
			res: connect.Resource{
				Product:  "redis",
				Location: "fr-par-1",
				ID:       "44444444-4444-4444-4444-444444444444",
				Username: "redis",
			},
			want: "Secret used to enable cli-connect. Connect directly using: scw redis cluster connect 44444444-4444-4444-4444-444444444444 --zone fr-par-1",
		},
		{
			name: "datawarehouse",
			res: connect.Resource{
				Product:  "datawarehouse",
				Location: "fr-par",
				ID:       "55555555-5555-5555-5555-555555555555",
				Username: "scwadmin",
			},
			want: "Secret used to enable cli-connect. Connect directly using: scw datawarehouse deployment connect 55555555-5555-5555-5555-555555555555 --region fr-par",
		},
		{
			name: "unknown product falls back to the generic description",
			res: connect.Resource{
				Product:  "mnq",
				Location: "fr-par",
				ID:       "66666666-6666-6666-6666-666666666666",
				Username: "user",
			},
			want: "secret used by the cli to ease up connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.Description(); got != tt.want {
				t.Fatalf("Description() = %q, want %q", got, tt.want)
			}
		})
	}
}

func Test_Resource_Path(t *testing.T) {
	tests := []struct {
		name string
		res  connect.Resource
		want string
	}{
		{
			name: "rdb",
			res: connect.Resource{
				Product:  "rdb",
				Location: "fr-par",
				ID:       "11111111-1111-1111-1111-111111111111",
				Username: "foobar",
			},
			want: "/rdb/fr-par/11111111-1111-1111-1111-111111111111/foobar",
		},
		{
			name: "redis",
			res: connect.Resource{
				Product:  "redis",
				Location: "fr-par-1",
				ID:       "44444444-4444-4444-4444-444444444444",
				Username: "redis",
			},
			want: "/redis/fr-par-1/44444444-4444-4444-4444-444444444444/redis",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.res.Path(); got != tt.want {
				t.Fatalf("Path() = %q, want %q", got, tt.want)
			}
		})
	}
}
