package editor

import (
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
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

	actualRequest := createGetRequest(updateRequest, reflect.TypeOf(GetRequest{}))
	assert.Equal(t, expectedRequest, actualRequest)
}
