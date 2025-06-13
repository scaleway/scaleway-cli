package editor

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ghodss/yaml"
)

type MarshalMode = string

const (
	MarshalModeYAML = MarshalMode("yaml")
	MarshalModeJSON = MarshalMode("json")
)

var (
	MarshalModeDefault = MarshalModeYAML
	MarshalModeEnum    = []MarshalMode{MarshalModeYAML, MarshalModeJSON}
)

func marshal(i any, mode MarshalMode) ([]byte, error) {
	if mode == "" {
		mode = MarshalModeDefault
	}

	var marshaledData []byte
	var err error

	switch mode {
	case MarshalModeYAML:
		marshaledData, err = yaml.Marshal(i)
	case MarshalModeJSON:
		marshaledData, err = json.MarshalIndent(i, "", "  ")
	}
	if err != nil {
		return marshaledData, err
	}
	if marshaledData != nil {
		return marshaledData, err
	}

	return nil, fmt.Errorf("unknown marshal mode %q", mode)
}

func unmarshal(data []byte, i any, mode MarshalMode) error {
	if mode == "" {
		mode = MarshalModeDefault
	}

	switch mode {
	case MarshalModeYAML:
		return yaml.Unmarshal(data, i)
	case MarshalModeJSON:
		return json.Unmarshal(data, i)
	}

	return fmt.Errorf("unknown marshal mode %q", mode)
}

// removeFields remove some fields from marshaled data
func removeFields(data []byte, mode MarshalMode, fields []string) ([]byte, error) {
	i := map[string]any{}
	err := unmarshal(data, &i, mode)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}
	DeleteRecursive(i, fields...)

	return marshal(i, mode)
}

func addTemplate(content []byte, template string, mode MarshalMode) []byte {
	if mode != MarshalModeYAML || len(template) == 0 {
		return content
	}
	newContent := []byte(nil)

	for _, line := range strings.Split(template, "\n") {
		newContent = append(newContent, []byte("#"+line+"\n")...)
	}

	newContent = append(newContent, content...)

	return newContent
}
