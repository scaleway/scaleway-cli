//go:build wasm

package wasm

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

type ConfigureOutputConfig struct {
	Width int  `js:"width"`
	Color bool `js:"color"`
}

type ConfigureOutputResponse struct {
}

func ConfigureOutput(cfg *ConfigureOutputConfig) (*ConfigureOutputResponse, error) {
	terminal.Width = cfg.Width
	color.NoColor = !cfg.Color

	return &ConfigureOutputResponse{}, nil
}
