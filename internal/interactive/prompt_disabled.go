//go:build wasm

package interactive

import (
	"context"
	"fmt"
)

type PromptPasswordConfig struct {
	Ctx    context.Context
	Prompt string
}

func PromptPasswordWithConfig(config *PromptPasswordConfig) (string, error) {
	return "", fmt.Errorf("prompt is disabled for this build")
}

type PromptBoolConfig struct {
	Ctx          context.Context
	Prompt       string
	DefaultValue bool
}

func PromptBoolWithConfig(config *PromptBoolConfig) (bool, error) {
	return config.DefaultValue, nil
}

type PromptStringConfig struct {
	Ctx             context.Context
	Prompt          string
	DefaultValue    string
	DefaultValueDoc string
	ValidateFunc    ValidateFunc
}

func PromptStringWithConfig(config *PromptStringConfig) (string, error) {
	return config.DefaultValue, nil
}

type ReadlineConfig struct {
	Ctx          context.Context
	Prompt       string
	PromptFunc   func(string) string
	Password     bool
	ValidateFunc ValidateFunc
	DefaultValue string
}

func Readline(config *ReadlineConfig) (string, error) {
	return "", fmt.Errorf("prompt is disabled for this build")
}
