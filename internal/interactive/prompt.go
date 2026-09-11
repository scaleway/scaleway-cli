//go:build !wasm

package interactive

import (
	"context"
	"os"
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

type PromptPasswordConfig struct {
	Prompt string
}

func PromptPasswordWithConfig(ctx context.Context, config *PromptPasswordConfig) (string, error) {
	return PromptPassword(ctx, config.Prompt)
}

func PromptPassword(ctx context.Context, prompt string) (string, error) {
	return Readline(ctx, &ReadlineConfig{
		Prompt:   prompt + ": ",
		Password: true,
	})
}

type PromptBoolConfig struct {
	Prompt       string
	DefaultValue bool
}

func PromptBoolWithConfig(ctx context.Context, config *PromptBoolConfig) (bool, error) {
	return PromptBool(ctx, config.Prompt, config.DefaultValue)
}

func PromptBool(ctx context.Context, prompt string, defaultValue bool) (bool, error) {
	for {
		promptText := terminal.Style(prompt, color.Bold)
		if defaultValue {
			promptText += " (Y/n): "
		} else {
			promptText += " (y/N): "
		}

		str, err := Readline(ctx, &ReadlineConfig{
			Prompt: promptText,
		})
		if err != nil {
			return false, err
		}

		switch strings.ToLower(str) {
		case "":
			return defaultValue, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		// update the config.Readline to say "invalid value" or something
	}
}

type PromptStringConfig struct {
	Prompt          string
	DefaultValue    string
	DefaultValueDoc string
	ValidateFunc    ValidateFunc
}

func PromptStringWithConfig(ctx context.Context, config *PromptStringConfig) (string, error) {
	return PromptString(
		ctx,
		config.Prompt,
		config.DefaultValue,
		config.DefaultValueDoc,
		config.ValidateFunc,
	)
}

func PromptString(
	ctx context.Context,
	prompt string,
	defaultValue string,
	defaultValueDoc string,
	validateFunc ValidateFunc,
) (string, error) {
	promptText := terminal.Style(prompt, color.Bold)
	if defaultValue != "" {
		promptText += terminal.Style(" (default: "+defaultValueDoc+")", color.Italic)
	}
	promptText += ": "

	v, err := Readline(ctx, &ReadlineConfig{
		Prompt:       promptText,
		ValidateFunc: validateFunc,
		DefaultValue: defaultValue,
	})
	if err != nil {
		return v, err
	}
	if v == "" {
		v = defaultValue
	}

	return v, err
}

type ReadlineHandler struct {
	rl *readline.Instance
}

func (h *ReadlineHandler) SetPrompt(prompt string) {
	h.rl.SetPrompt(prompt)
}

type ReadlineConfig struct {
	Prompt       string
	PromptFunc   func(string) string
	Password     bool
	ValidateFunc ValidateFunc
	DefaultValue string
}

func Readline(ctx context.Context, config *ReadlineConfig) (string, error) {
	promptFunc := func(string) string {
		return config.Prompt
	}
	if config.PromptFunc != nil {
		promptFunc = config.PromptFunc
	}
	validateFunc := defaultValidateFunc
	if config.ValidateFunc != nil {
		validateFunc = config.ValidateFunc
	}

	// Extract mock responses from context for testing
	var mockResponses *[]string
	if contextValue := ctx.Value(contextKey); contextValue != nil {
		if mockValues, ok := contextValue.(*[]string); ok && mockValues != nil {
			mockResponses = mockValues
		}
	}

	var promptHandler *ReadlineHandler
	rl, err := readline.NewEx(&readline.Config{
		Prompt:                 promptFunc(""),
		DisableAutoSaveHistory: true,
		InterruptPrompt:        "^C",
		EnableMask:             config.Password,
		FuncIsTerminal: func() bool {
			return IsInteractive
		},
		Stdin: &mockResponseReader{
			mockResponses: mockResponses,
			defaultReader: os.Stdin,
		},
		Listener: readline.FuncListener(
			func(line []rune, _ int, _ rune) (newLine []rune, newPos int, ok bool) {
				value := string(line)
				promptHandler.SetPrompt(promptFunc(value))
				promptHandler.rl.Refresh()

				return nil, 0, false
			},
		),
	})
	if err != nil {
		return "", err
	}

	promptHandler = &ReadlineHandler{rl: rl}
	var s string
	for {
		s, err = rl.Readline()
		// If readline returns an error we return it
		if err != nil {
			switch err.Error() {
			case "Interrupt":
				return "", &InterruptError{}
			case "EOF":
				return "", nil
			default:
				return "", err
			}
		}

		if s == "" {
			s = config.DefaultValue
		}

		// Handle user input validation
		err = validateFunc(s)
		// If ValidateFunc returns an error we print it and Readline again
		if err != nil {
			s, err := human.Marshal(err, nil)
			if err != nil {
				return "", err
			}
			_, err = Println(s)
			if err != nil {
				return "", err
			}

			continue
		}

		// If there was no validation error return the result
		return s, nil
	}
}
