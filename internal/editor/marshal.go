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

func marshal(i interface{}, mode MarshalMode) ([]byte, error) {
	if mode == "" {
		mode = MarshalModeYAML
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

func unmarshal(data []byte, i interface{}, mode MarshalMode) error {
	if mode == "" {
		mode = MarshalModeYAML
	}

	switch mode {
	case MarshalModeYAML:
		return yaml.Unmarshal(data, i)
	case MarshalModeJSON:
		return json.Unmarshal(data, i)
	}

	return fmt.Errorf("unknown marshal mode %q", mode)
}
