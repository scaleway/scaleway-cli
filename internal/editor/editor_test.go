package editor_test

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
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
	assert.Nil(t, err)
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

	editedUpdateRequestI, err := editor.UpdateResourceEditor(resource, updateRequest, &editor.Config{})
	assert.Nil(t, err)
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

	editedUpdateRequestI, err := editor.UpdateResourceEditor(resource, updateRequest, &editor.Config{
		EditedResource: `
id: uuid
env: {}
`,
	})
	assert.Nil(t, err)
	editedUpdateRequest := editedUpdateRequestI.(*UpdateRequest)
	assert.NotNil(t, editedUpdateRequest.Env)
	assert.True(t, len(*editedUpdateRequest.Env) == 0)
}
