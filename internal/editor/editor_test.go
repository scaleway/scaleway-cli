package editor_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_updateResourceEditor(t *testing.T) {
	editor.SkipEditor = true

	resource := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"name",
	}
	updateRequest := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"",
	}

	_, err := editor.UpdateResourceEditor(t.Context(), resource, updateRequest, &editor.Config{})
	assert.NoError(t, err)
}

func Test_updateResourceEditor_pointers(t *testing.T) {
	editor.SkipEditor = true

	type UpdateRequest struct {
		ID   string
		Name *string
	}
	resource := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"name",
	}

	updateRequest := &UpdateRequest{
		"uuid",
		nil,
	}

	editedUpdateRequestI, err := editor.UpdateResourceEditor(
		t.Context(),
		resource,
		updateRequest,
		&editor.Config{},
	)
	require.NoError(t, err)
	editedUpdateRequest := editedUpdateRequestI.(*UpdateRequest)

	assert.NotNil(t, editedUpdateRequest.Name)
	assert.Equal(t, resource.Name, *editedUpdateRequest.Name)
}

func Test_updateResourceEditor_map(t *testing.T) {
	editor.SkipEditor = true

	type UpdateRequest struct {
		ID  string             `json:"id"`
		Env *map[string]string `json:"env"`
	}
	resource := &struct {
		ID  string            `json:"id"`
		Env map[string]string `json:"env"`
	}{
		"uuid",
		map[string]string{
			"foo": "bar",
		},
	}

	updateRequest := &UpdateRequest{
		"uuid",
		nil,
	}

	editedUpdateRequestI, err := editor.UpdateResourceEditor(
		t.Context(),
		resource,
		updateRequest,
		&editor.Config{
			EditedResource: `
id: uuid
env: {}
`,
		},
	)
	require.NoError(t, err)
	editedUpdateRequest := editedUpdateRequestI.(*UpdateRequest)
	assert.NotNil(t, editedUpdateRequest.Env)
	assert.Empty(t, *editedUpdateRequest.Env)
}

func TestOpenFile(t *testing.T) {
	file := filepath.Join(t.TempDir(), "my config.yaml")
	editorPath := filepath.Join(t.TempDir(), "editor with spaces.exe")
	require.NoError(t, os.WriteFile(editorPath, nil, 0o700))
	t.Setenv("VISUAL", "")
	for _, tc := range []struct {
		name    string
		setting string
		args    []string
		code    int
	}{
		{"flags", "editor --wait", []string{"editor", "--wait", file}, 0},
		{"executable path", editorPath, []string{editorPath, file}, 0},
		{"failed editor", "editor", []string{"editor", file}, 17},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Setenv("EDITOR", tc.setting)
			called := false
			ctx := core.InjectMeta(
				t.Context(),
				&core.Meta{OverrideExec: func(cmd *exec.Cmd) (int, error) {
					called = true
					require.Equal(t, tc.args, cmd.Args)

					return tc.code, nil
				}},
			)
			err := editor.OpenFile(ctx, file)
			require.True(t, called)
			if tc.code == 0 {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, "17")
			}
		})
	}
}
