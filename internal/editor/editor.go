package editor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/config"
)

var SkipEditor = false

func edit(content []byte) ([]byte, error) {
	if SkipEditor == true {
		return content, nil
	}
	editionBuffer := bytes.NewBuffer(nil)

	defaultEditor := config.GetDefaultEditor()
	cmd := exec.Command(defaultEditor)
	cmd.Stdin = bytes.NewBuffer(content)
	cmd.Stdout = editionBuffer

	err := cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to edit temporary file: %w", err)
	}

	return editionBuffer.Bytes(), nil
}

// editor takes a complete resource and a partial resourceUpdate
func editor(resource interface{}, resourceUpdate interface{}, editedJson ...string) (interface{}, error) {
	resourceV := reflect.ValueOf(resource)
	resourceUpdateV := reflect.ValueOf(resourceUpdate)

	// Create a new resourceUpdate that will be edited
	// It will allow user to edit it, then we will extract diff to perform update
	resourceUpdateToEdit := reflect.New(resourceUpdateV.Type().Elem())
	valueMapper(resourceUpdateV, resourceUpdateToEdit)
	valueMapper(resourceV, resourceUpdateToEdit)

	tmpArgsJson, err := json.MarshalIndent(resourceUpdateToEdit.Interface(), "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to convert update request to json: %w", err)
	}

	tmpArgsJson, err = edit(tmpArgsJson)
	if err != nil {
		return nil, fmt.Errorf("failed to edit json: %w", err)
	}

	if len(editedJson) == 1 {
		tmpArgsJson = []byte(editedJson[0])
	}

	// Create a new resourceUpdate as destination for edited one
	resourceUpdateEdited := reflect.New(resourceUpdateV.Type().Elem())

	err = json.Unmarshal(tmpArgsJson, resourceUpdateEdited.Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to marshal edited data: %w", err)
	}

	return resourceUpdateEdited.Interface(), nil
}

// Editor takes a partial UpdateResourceRequest, the type of the GetResourceRequest and a function that return the resource using GetResourceRequest
func Editor(resourceUpdate interface{}, resourceGetType reflect.Type, resourceGetter func(interface{}) (interface{}, error)) (interface{}, error) {
	resourceUpdateV := reflect.ValueOf(resourceUpdate)

	resourceGetterArg := reflect.New(resourceGetType).Interface()
	resourceGetterArgV := reflect.ValueOf(resourceGetterArg)

	// Fill Getter args using Update arg content
	valueMapper(resourceUpdateV, resourceGetterArgV)

	resource, err := resourceGetter(resourceGetterArg)
	if err != nil {
		return nil, fmt.Errorf("failed to get resource: %w", err)
	}

	editedArgs, err := editor(resource, resourceUpdateV.Interface())
	if err != nil {
		return nil, fmt.Errorf("failed to edit update arguments: %w", err)
	}

	return editedArgs, nil
}
