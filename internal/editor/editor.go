package editor

import (
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/config"
)

var SkipEditor = false
var marshalMode = MarshalModeYaml

type GetResourceFunc func(interface{}) (interface{}, error)

// edit create a temporary file with given content, start a text editor then return edited content
// temporary file will be deleted on complete
// temporary file is not deleted if edit fails
func edit(content []byte) ([]byte, error) {
	if SkipEditor == true {
		return content, nil
	}

	tmpFileName, err := createTemporaryFile(content, marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}

	defaultEditor := config.GetDefaultEditor()
	editorAndArguments := strings.Fields(defaultEditor)
	args := []string{tmpFileName}
	if len(editorAndArguments) > 1 {
		args = append(editorAndArguments[1:], args...)
	}
	cmd := exec.Command(editorAndArguments[0], args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to edit temporary file %q: %w", tmpFileName, err)
	}

	editedContent, err := readAndDeleteFile(tmpFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read and delete temporary file: %w", err)
	}

	return editedContent, nil
}

// updateResourceEditor takes a complete resource and a partial updateRequest
// will return a copy of updateRequest that has been edited
func updateResourceEditor(resource interface{}, updateRequest interface{}, editedResource ...string) (interface{}, error) {
	// Create a copy of updateRequest completed with resource content
	completeUpdateRequest := copyAndCompleteUpdateRequest(updateRequest, resource)

	updateRequestMarshaled, err := Marshal(completeUpdateRequest, marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update request: %w", err)
	}

	// Start text editor to edit marshaled request
	updateRequestMarshaled, err = edit(updateRequestMarshaled)
	if err != nil {
		return nil, fmt.Errorf("failed to edit marshalled data: %w", err)
	}

	// If editedResource is present, override edited resource
	// This is useful for testing purpose
	if len(editedResource) == 1 {
		updateRequestMarshaled = []byte(editedResource[0])
	}

	// Create a new updateRequest as destination for edited yaml/json
	// Must be a new one to avoid merge of maps content
	updateRequestEdited := newRequest(updateRequest)

	// TODO: fill updateRequestEdited with only edited fields and fields present in updateRequest
	// TODO: fields should be compared with completeUpdateRequest to find edited ones

	err = Unmarshal(updateRequestMarshaled, updateRequestEdited, marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal edited data: %w", err)
	}

	return updateRequestEdited, nil
}

// UpdateResourceEditor takes:
// - a partial UpdateRequest
// - the type of the GetRequest
// - a function that return the resource using GetRequest
func UpdateResourceEditor(updateRequest interface{}, getRequestType reflect.Type, getResource GetResourceFunc) (interface{}, error) {
	getRequest := createGetRequest(updateRequest, getRequestType)

	resource, err := getResource(getRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	// Start edition of UpdateRequest
	editedUpdateRequest, err := updateResourceEditor(resource, updateRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to edit update arguments: %w", err)
	}

	return editedUpdateRequest, nil
}
