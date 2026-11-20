package mongodb_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/mongodb/v1"
	mongodbSDK "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PickAdminUser(t *testing.T) {
	t.Parallel()

	dbAdmin := &mongodbSDK.UserRole{Role: mongodbSDK.UserRoleRoleDbAdmin}
	readWrite := &mongodbSDK.UserRole{Role: mongodbSDK.UserRoleRoleReadWrite}

	tests := []struct {
		name    string
		users   []*mongodbSDK.User
		want    string
		wantErr string
	}{
		{
			name:    "empty",
			wantErr: "no user found on instance",
		},
		{
			name: "prefers db_admin over earlier non-admin",
			users: []*mongodbSDK.User{
				{Name: "app", Roles: []*mongodbSDK.UserRole{readWrite}},
				{Name: "admin", Roles: []*mongodbSDK.UserRole{dbAdmin}},
			},
			want: "admin",
		},
		{
			name: "first db_admin wins among several",
			users: []*mongodbSDK.User{
				{Name: "admin-a", Roles: []*mongodbSDK.UserRole{dbAdmin}},
				{Name: "admin-b", Roles: []*mongodbSDK.UserRole{dbAdmin}},
			},
			want: "admin-a",
		},
		{
			name: "falls back to first user when no db_admin",
			users: []*mongodbSDK.User{
				{Name: "app", Roles: []*mongodbSDK.UserRole{readWrite}},
				{Name: "other", Roles: []*mongodbSDK.UserRole{readWrite}},
			},
			want: "app",
		},
		{
			name:  "user with no roles still usable as fallback",
			users: []*mongodbSDK.User{{Name: "lonely"}},
			want:  "lonely",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := mongodb.PickAdminUser(tt.users)
			if tt.wantErr != "" {
				require.EqualError(t, err, tt.wantErr)

				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetEndpoints(t *testing.T) {
	t.Parallel()

	public := &mongodbSDK.Endpoint{
		DNSRecord:     "public.example",
		PublicNetwork: &mongodbSDK.EndpointPublicNetworkDetails{},
	}
	private := &mongodbSDK.Endpoint{
		DNSRecord:      "private.example",
		PrivateNetwork: &mongodbSDK.EndpointPrivateNetworkDetails{},
	}

	t.Run("public found", func(t *testing.T) {
		t.Parallel()
		got, err := mongodb.GetPublicEndpoint([]*mongodbSDK.Endpoint{private, public})
		require.NoError(t, err)
		assert.Equal(t, "public.example", got.DNSRecord)
	})

	t.Run("public missing", func(t *testing.T) {
		t.Parallel()
		_, err := mongodb.GetPublicEndpoint([]*mongodbSDK.Endpoint{private})
		require.EqualError(t, err, "public endpoint not found")
	})

	t.Run("private found", func(t *testing.T) {
		t.Parallel()
		got, err := mongodb.GetPrivateEndpoint([]*mongodbSDK.Endpoint{public, private})
		require.NoError(t, err)
		assert.Equal(t, "private.example", got.DNSRecord)
	})

	t.Run("private missing", func(t *testing.T) {
		t.Parallel()
		_, err := mongodb.GetPrivateEndpoint([]*mongodbSDK.Endpoint{public})
		require.EqualError(t, err, "private endpoint not found")
	})

	t.Run("empty list", func(t *testing.T) {
		t.Parallel()
		_, err := mongodb.GetPublicEndpoint(nil)
		require.EqualError(t, err, "public endpoint not found")
	})
}
