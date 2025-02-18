package editor

import (
	"fmt"
	"os"
)

func temporaryFileNamePattern(marshalMode MarshalMode) string {
	pattern := "scw-updateResourceEditor" // #nosec G101
	switch marshalMode {
	case MarshalModeYAML:
		pattern += "*.yml"
	case MarshalModeJSON:
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
