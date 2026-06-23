package secret

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
			raw:          "11111111-1111-1111-1111-111111111111",
			wantID:       "11111111-1111-1111-1111-111111111111",
			wantRevision: "latest",
		},
		{
			raw:          "11111111-1111-1111-1111-111111111111@2",
			wantID:       "11111111-1111-1111-1111-111111111111",
			wantRevision: "2",
		},
		{
			raw:          "11111111-1111-1111-1111-111111111111:api-key",
			wantID:       "11111111-1111-1111-1111-111111111111",
			wantRevision: "latest",
			wantField:    "api-key",
		},
		{
			raw:          "11111111-1111-1111-1111-111111111111@latest:api-key",
			wantID:       "11111111-1111-1111-1111-111111111111",
			wantRevision: "latest",
			wantField:    "api-key",
		},
		{
			raw:          "my-secret",
			wantName:     "my-secret",
			wantPath:     "/",
			wantRevision: "latest",
		},
		{
			raw:          "db/my-secret",
			wantName:     "my-secret",
			wantPath:     "/db",
			wantRevision: "latest",
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
			raw:     "my-secret:",
			wantErr: true,
		},
		{
			raw:     "my-secret@",
			wantErr: true,
		},
		{
			raw:     "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.raw, func(t *testing.T) {
			ref, err := parseSecretRef(tt.raw)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.wantID, ref.secretID)
			assert.Equal(t, tt.wantName, ref.secretName)
			assert.Equal(t, tt.wantPath, ref.secretPath)
			assert.Equal(t, tt.wantRevision, ref.revision)
			assert.Equal(t, tt.wantField, ref.field)
		})
	}
}

func Test_RenderTemplate_NoRefs(t *testing.T) {
	rendered, err := renderTemplate("hello world\nno secrets here", nil, "")
	require.NoError(t, err)
	assert.Equal(t, "hello world\nno secrets here", rendered)
}

func Test_IsSecretUUID(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"11111111-1111-1111-1111-111111111111", true},
		{"abcdef01-abcd-abcd-abcd-abcdef012345", true},
		{"not-a-uuid", false},
		{"11111111-1111-1111-1111-11111111111G", false}, // uppercase G
		{"11111111-1111-1111-1111-111111111111X", false}, // too long
		{"", false},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			assert.Equal(t, tt.want, isSecretUUID(tt.input))
		})
	}
}
