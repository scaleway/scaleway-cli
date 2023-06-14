//go:build wasm

package wasm

import (
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

type ConfigureOutputConfig struct {
	Width int `js:"width"`
}

type ConfigureOutputResponse struct {
}

func ConfigureOutput(cfg *ConfigureOutputConfig) (*ConfigureOutputResponse, error) {
	terminal.Width = cfg.Width

	return &ConfigureOutputResponse{}, nil
}
