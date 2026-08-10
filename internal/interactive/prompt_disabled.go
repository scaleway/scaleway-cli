//go:build wasm

package interactive

import (
	"context"
	"fmt"
)

type PromptPasswordConfig struct {
	Prompt string
}

func PromptPasswordWithConfig(ctx context.Context, config *PromptPasswordConfig) (string, error) {
	return "", fmt.Errorf("prompt is disabled for this build")
}

func PromptPassword(ctx context.Context, prompt string) (string, error) {
	return "", fmt.Errorf("prompt is disabled for this build")
}

type PromptBoolConfig struct {
	Prompt       string
	DefaultValue bool
}

func PromptBoolWithConfig(ctx context.Context, config *PromptBoolConfig) (bool, error) {
	return config.DefaultValue, nil
}

func PromptBool(ctx context.Context, prompt string, defaultValue bool) (bool, error) {
	return defaultValue, nil
}

type PromptStringConfig struct {
	Prompt          string
	DefaultValue    string
	DefaultValueDoc string
	ValidateFunc    ValidateFunc
}

func PromptStringWithConfig(ctx context.Context, config *PromptStringConfig) (string, error) {
	return config.DefaultValue, nil
}

func PromptString(
	ctx context.Context,
	prompt string,
	defaultValue string,
	defaultValueDoc string,
	validateFunc ValidateFunc,
) (string, error) {
	return defaultValue, nil
}

type ReadlineHandler struct{}

func (h *ReadlineHandler) SetPrompt(prompt string) {}

type ReadlineConfig struct {
	Prompt       string
	PromptFunc   func(string) string
	Password     bool
	ValidateFunc ValidateFunc
	DefaultValue string
}

func Readline(ctx context.Context, config *ReadlineConfig) (string, error) {
	return "", fmt.Errorf("prompt is disabled for this build")
}
