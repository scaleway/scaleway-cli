package editor

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/internal/config"
)

var SkipEditor = false
var marshalMode = MarshalModeYAML

type GetResourceFunc func(interface{}) (interface{}, error)
type Config struct {
	// PutRequest means that the request replace all fields
	// If false, fields that were not edited will not be sent
	// If true, all fields will be sent
	PutRequest bool

	MarshalMode MarshalMode

	// If not empty, this will replace edited text as if it was edited in the terminal
	// Should be paired with global SkipEditor as true, useful for tests
	editedResource string
}

func editorPathAndArgs(fileName string) (string, []string) {
	defaultEditor := config.GetDefaultEditor()
	editorAndArguments := strings.Fields(defaultEditor)
	args := []string{fileName}

	if len(editorAndArguments) > 1 {
		args = append(editorAndArguments[1:], args...)
	}

	return editorAndArguments[0], args
}

// edit create a temporary file with given content, start a text editor then return edited content
// temporary file will be deleted on complete
// temporary file is not deleted if edit fails
func edit(content []byte) ([]byte, error) {
	if SkipEditor {
		return content, nil
	}

	tmpFileName, err := createTemporaryFile(content, marshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tmpFileName)

	editorPath, args := editorPathAndArgs(tmpFileName)
	cmd := exec.Command(editorPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	err = cmd.Run()
	if err != nil {
		return nil, fmt.Errorf("failed to edit temporary file %q: %w", tmpFileName, err)
	}

	editedContent, err := os.ReadFile(tmpFileName)
	if err != nil {
		return nil, fmt.Errorf("failed to read temporary file %q: %w", tmpFileName, err)
	}

	return editedContent, nil
}

// updateResourceEditor takes a complete resource and a partial updateRequest
// will return a copy of updateRequest that has been edited
func updateResourceEditor(resource interface{}, updateRequest interface{}, cfg *Config) (interface{}, error) {
	// Create a copy of updateRequest completed with resource content
	completeUpdateRequest := copyAndCompleteUpdateRequest(updateRequest, resource)

	// TODO: fields present in updateRequest should be removed from marshal
	// ex: namespace_id, region, zone
	// Currently not an issue as fields that should be removed are mostly path parameter /{zone}/namespace/{namespace_id}
	// Path parameter have "-" as json tag and are not marshaled

	updateRequestMarshaled, err := marshal(completeUpdateRequest, cfg.MarshalMode)
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
	if cfg.editedResource != "" {
		updateRequestMarshaled = []byte(cfg.editedResource)
	}

	// Create a new updateRequest as destination for edited yaml/json
	// Must be a new one to avoid merge of maps content
	updateRequestEdited := newRequest(updateRequest)

	// TODO: if !putRequest
	// TODO: fill updateRequestEdited with only edited fields and fields present in updateRequest
	// TODO: fields should be compared with completeUpdateRequest to find edited ones

	err = unmarshal(updateRequestMarshaled, updateRequestEdited, cfg.MarshalMode)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal edited data: %w", err)
	}

	return updateRequestEdited, nil
}

// UpdateResourceEditor takes a complete resource and a partial updateRequest
// will return a copy of updateRequest that has been edited
// Only edited fields will be present in returned updateRequest
// If putRequest is true, all fields will be present, edited or not
func UpdateResourceEditor(resource interface{}, updateRequest interface{}, cfg *Config) (interface{}, error) {
	return updateResourceEditor(resource, updateRequest, cfg)
}
