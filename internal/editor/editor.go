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

// editor takes a complete resource and a partial resourceUpdate
func editor(resource interface{}, updateResourceRequest interface{}, editedJson ...string) (interface{}, error) {
	resourceV := reflect.ValueOf(resource)
	updateResourceRequestV := reflect.ValueOf(updateResourceRequest)

	// Create a new updateResourceRequest that will be edited
	// It will allow user to edit it, then we will extract diff to perform update
	updateResourceRequestToEditV := reflect.New(updateResourceRequestV.Type().Elem())
	valueMapper(updateResourceRequestToEditV, updateResourceRequestV)
	valueMapper(updateResourceRequestToEditV, resourceV)

	updateResourceRequestJson, err := Marshal(updateResourceRequestToEditV.Interface(), marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal update request: %w", err)
	}

	// Start text editor to edit json
	updateResourceRequestJson, err = edit(updateResourceRequestJson)
	if err != nil {
		return nil, fmt.Errorf("failed to edit marshalled data: %w", err)
	}

	// If editedJson is present, override edited json
	// This is useful for testing purpose
	if len(editedJson) == 1 {
		updateResourceRequestJson = []byte(editedJson[0])
	}

	// Create a new updateResourceRequest as destination for edited one
	updateResourceRequestEdited := reflect.New(updateResourceRequestV.Type().Elem())

	err = Unmarshal(updateResourceRequestJson, updateResourceRequestEdited.Interface(), marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal edited data: %w", err)
	}

	return updateResourceRequestEdited.Interface(), nil
}

// Editor takes a partial UpdateResourceRequest, the type of the GetResourceRequest and a function that return the resource using GetResourceRequest
func Editor(updateResourceRequest interface{}, getResourceRequestType reflect.Type, getResource GetResourceFunc) (interface{}, error) {
	updateResourceRequestV := reflect.ValueOf(updateResourceRequest)

	// Create GetResourceRequest to be able to fetch the actual resource
	getResourceRequest := reflect.New(getResourceRequestType).Interface()
	getResourceRequestV := reflect.ValueOf(getResourceRequest)

	// Fill GetResourceRequest args using Update arg content
	// This should copy important argument like ResourceID
	valueMapper(getResourceRequestV, updateResourceRequestV)

	resource, err := getResource(getResourceRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	// Start edition of UpdateResourceRequest
	editedArgs, err := editor(resource, updateResourceRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to edit update arguments: %w", err)
	}

	return editedArgs, nil
}
