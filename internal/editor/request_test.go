package editor_test

import (
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/editor"
)

func Test_createGetResourceRequest(t *testing.T) {
	type GetRequest struct {
		ID string
	}
	updateRequest := struct {
		ID   string
		Name string
	}{"uuid", ""}
	expectedRequest := &GetRequest{"uuid"}

	actualRequest := editor.CreateGetRequest(updateRequest, reflect.TypeOf(GetRequest{}))
	assert.Equal(t, expectedRequest, actualRequest)
}
