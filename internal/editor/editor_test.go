package editor_test

import (
	"testing"

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

	_, err := editor.UpdateResourceEditor(resource, updateRequest, &editor.Config{})
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
