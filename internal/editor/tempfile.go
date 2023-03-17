package editor

import (
	"fmt"
	"os"
)

func temporaryFileNamePattern(marshalMode MarshalMode) string {
	pattern := "scw-editor"
	switch marshalMode {
	case MarshalModeYaml:
		pattern += "*.yml"
	case MarshalModeJson:
		pattern += "*.json"
	}
	return pattern
}

func createTemporaryFile(content []byte, marshalMode MarshalMode) (string, error) {
	tmpFile, err := os.CreateTemp(os.TempDir(), temporaryFileNamePattern(marshalMode))
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	_, err = tmpFile.Write(content)
	if err != nil {
		return "", fmt.Errorf("failed to write to file %q: %w", tmpFile.Name(), err)
	}
	err = tmpFile.Close()
	if err != nil {
		return "", fmt.Errorf("failed to close file %q: %w", tmpFile.Name(), err)
	}

	return tmpFile.Name(), nil
}

func readAndDeleteFile(name string) ([]byte, error) {
	content, err := os.ReadFile(name)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %q: %w", name, err)
	}

	err = os.Remove(name)
	if err != nil {
		return nil, fmt.Errorf("failed to delete file %q: %w", name, err)
	}

	return content, nil
}
