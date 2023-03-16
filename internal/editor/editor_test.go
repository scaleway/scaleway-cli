package editor

import (
	"log"
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_editor(t *testing.T) {
	SkipEditor = true

	resource := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"name",
	}
	resourceUpdate := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"",
	}

	_, err := editor(resource, resourceUpdate)
	assert.Nil(t, err)
}

func Test_editor_pointers(t *testing.T) {
	SkipEditor = true

	type ResourceUpdate struct {
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

	resourceUpdate := &ResourceUpdate{
		"uuid",
		nil,
	}

	editedResourceUpdateI, err := editor(resource, resourceUpdate)
	assert.Nil(t, err)
	editedResourceUpdate := editedResourceUpdateI.(*ResourceUpdate)
	assert.NotNil(t, editedResourceUpdate.Name)
	assert.Equal(t, resource.Name, *editedResourceUpdate.Name)
}

func Test_editor_map(t *testing.T) {
	SkipEditor = true

	type ResourceUpdate struct {
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

	resourceUpdate := &ResourceUpdate{
		"uuid",
		nil,
	}

	SkipEditor = true
	editedResourceUpdateI, err := editor(resource, resourceUpdate, `
	{"id":"uuid", "env": {}}
`)
	assert.Nil(t, err)
	editedResourceUpdate := editedResourceUpdateI.(*ResourceUpdate)
	assert.NotNil(t, editedResourceUpdate.Env)
	assert.True(t, len(*editedResourceUpdate.Env) == 0)
}

func TestEditor(t *testing.T) {
	SkipEditor = true

	resource := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"name",
	}
	resourceUpdate := &struct {
		ID   string
		Name string
	}{
		"uuid",
		"",
	}

	editedResource, err := Editor(resourceUpdate, reflect.TypeOf(*resource), func(i interface{}) (interface{}, error) {
		return resource, nil
	})
	assert.Nil(t, err)
	log.Println(editedResource)
}
