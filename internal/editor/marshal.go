package editor

import (
	"encoding/json"
	"fmt"

	"github.com/ghodss/yaml"
)

type MarshalMode string

const (
	MarshalModeYAML = MarshalMode("yaml")
	MarshalModeJSON = MarshalMode("json")
)

func Marshal(i interface{}, mode MarshalMode) ([]byte, error) {
	switch mode {
	case MarshalModeYAML:
		return yaml.Marshal(i)
	case MarshalModeJSON:
		return json.MarshalIndent(i, "", "  ")
	}

	return nil, fmt.Errorf("unknown marshal mode %q", mode)
}

func Unmarshal(data []byte, i interface{}, mode MarshalMode) error {
	switch mode {
	case MarshalModeYAML:
		return yaml.Unmarshal(data, i)
	case MarshalModeJSON:
		return json.Unmarshal(data, i)
	}

	return fmt.Errorf("unknown marshal mode %q", mode)
}
