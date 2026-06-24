package secret_test

import (
	"testing"

	secret "github.com/scaleway/scaleway-cli/v2/internal/namespaces/secret/v1beta1"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testUUID     = "11111111-1111-1111-1111-111111111111"
	testLatest   = "latest"
	testMySecret = "my-secret"
)

func Test_ParseSecretRef(t *testing.T) {
	tests := []struct {
		raw          string
		wantID       string
		wantName     string
		wantPath     string
		wantRevision string
		wantField    string
		wantErr      bool
	}{
		{
			raw:          testUUID,
			wantID:       testUUID,
			wantRevision: testLatest,
		},
		{
			raw:          testUUID + "@2",
			wantID:       testUUID,
			wantRevision: "2",
		},
		{
			raw:          testUUID + ":api-key",
			wantID:       testUUID,
			wantRevision: testLatest,
			wantField:    "api-key",
		},
		{
			raw:          testUUID + "@latest:api-key",
			wantID:       testUUID,
			wantRevision: testLatest,
			wantField:    "api-key",
		},
		{
			raw:          testMySecret,
			wantName:     testMySecret,
			wantPath:     "/",
			wantRevision: testLatest,
		},
		{
			raw:          "db/" + testMySecret,
			wantName:     testMySecret,
			wantPath:     "/db",
			wantRevision: testLatest,
		},
		{
			raw:          "my-app/db/password@2:key",
			wantName:     "password",
			wantPath:     "/my-app/db",
			wantRevision: "2",
			wantField:    "key",
		},
		{
			raw:     "@",
			wantErr: true,
		},
		{
			raw:     testMySecret + ":",
			wantErr: true,
		},
		{
			raw:     testMySecret + "@",
			wantErr: true,
		},
		{
			raw:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			ref, err := secret.ParseSecretRef(tt.raw)
			if tt.wantErr {
				require.Error(t, err)

				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantID, ref.SecretID)
			assert.Equal(t, tt.wantName, ref.SecretName)
			assert.Equal(t, tt.wantPath, ref.SecretPath)
			assert.Equal(t, tt.wantRevision, ref.Revision)
			assert.Equal(t, tt.wantField, ref.Field)
		})
	}
}

func Test_RenderTemplate_NoRefs(t *testing.T) {
	rendered, err := secret.RenderTemplate("hello world\nno secrets here", nil, "")
	require.NoError(t, err)
	assert.Equal(t, "hello world\nno secrets here", rendered)
}

func Test_IsSecretUUID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{testUUID, true},
		{"abcdef01-abcd-abcd-abcd-abcdef012345", true},
		{"not-a-uuid", false},
		{"11111111-1111-1111-1111-11111111111G", false},
		{"11111111-1111-1111-1111-111111111111X", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, secret.IsSecretUUID(tt.input))
		})
	}
}
