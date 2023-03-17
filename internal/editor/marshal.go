package editor

import (
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v3"
)

type MarshalMode string

const (
	MarshalModeYaml = MarshalMode("yaml")
	MarshalModeJson = MarshalMode("json")
)

func Marshal(i interface{}, mode MarshalMode) ([]byte, error) {
	switch mode {
	case MarshalModeYaml:
		return yaml.Marshal(i)
	case MarshalModeJson:
		return json.MarshalIndent(i, "", "  ")
	}

	return nil, fmt.Errorf("unknown marshal mode %q", mode)
}

func Unmarshal(data []byte, i interface{}, mode MarshalMode) error {
	switch mode {
	case MarshalModeYaml:
		return yaml.Unmarshal(data, i)
	case MarshalModeJson:
		return json.Unmarshal(data, i)
	}

	return fmt.Errorf("unknown marshal mode %q", mode)
}
