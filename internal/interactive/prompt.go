package interactive

import (
	"strings"

	"github.com/chzyer/readline"
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-cli/internal/terminal"
)

type PromptPasswordConfig struct {
	Prompt string
}

func PromptPasswordWithConfig(config *PromptPasswordConfig) (string, error) {
	return Readline(&ReadlineConfig{
		Prompt:   config.Prompt + ": ",
		Password: true,
	})
}

func PromptPassword(prompt string) (string, error) {
	return PromptPasswordWithConfig(&PromptPasswordConfig{Prompt: prompt})
}

type PromptBoolConfig struct {
	Prompt       string
	DefaultValue bool
}

func PromptBoolWithConfig(config *PromptBoolConfig) (bool, error) {
	for {
		prompt := terminal.Style(config.Prompt, color.Bold)
		if config.DefaultValue {
			prompt = prompt + " (Y/n): "
		} else {
			prompt = prompt + " (y/N): "
		}

		str, err := Readline(&ReadlineConfig{Prompt: prompt})
		if err != nil {
			return false, err
		}

		switch strings.ToLower(str) {
		case "":
			return config.DefaultValue, nil
		case "y", "yes":
			return true, nil
		case "n", "no":
			return false, nil
		}
		// update the config.Readline to say "invalid value" or something
	}
}

// TODO: cleanup
func PromptBool(prompt string) (bool, error) {
	return PromptBoolWithConfig(&PromptBoolConfig{
		Prompt: prompt,
	})
}

type PromptStringConfig struct {
	Prompt          string
	DefaultValue    string
	DefaultValueDoc string
	ValidateFunc    ValidateFunc
}

func PromptStringWithConfig(config *PromptStringConfig) (string, error) {
	prompt := terminal.Style(config.Prompt, color.Bold)
	if config.DefaultValue != "" {
		prompt += terminal.Style(" (default: "+config.DefaultValueDoc+")", color.Italic)
	}
	prompt += ": "

	v, err := Readline(&ReadlineConfig{
		Prompt:       prompt,
		ValidateFunc: config.ValidateFunc,
		DefaultValue: config.DefaultValue,
	})
	if err != nil {
		return v, err
	}
	if v == "" {
		v = config.DefaultValue
	}
	return v, err
}

func PromptString(prompt string) (string, error) {
	return PromptStringWithConfig(&PromptStringConfig{Prompt: prompt})
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

func Readline(config *ReadlineConfig) (string, error) {
	promptFunc := func(s string) string {
		return config.Prompt
	}
	if config.PromptFunc != nil {
		promptFunc = config.PromptFunc
	}
	validateFunc := defaultValidateFunc
	if config.ValidateFunc != nil {
		validateFunc = config.ValidateFunc
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
		Listener: readline.FuncListener(func(line []rune, pos int, key rune) (newLine []rune, newPos int, ok bool) {
			value := string(line)
			promptHandler.SetPrompt(promptFunc(value))
			promptHandler.rl.Refresh()
			return nil, 0, false
		}),
	})

	if err != nil {
		return "", err
	}

	promptHandler = &ReadlineHandler{rl: rl}
	s := ""
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
			Println(s)
			continue
		}

		// If there was no validation error return the result
		return s, nil
	}

}
