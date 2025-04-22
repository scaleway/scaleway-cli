//go:build wasm

package wasm

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
)

type AutoCompleteConfig struct {
	JWT          string   `js:"jwt"`
	LeftWords    []string `js:"leftWords"`
	SelectedWord string   `js:"selectedWord"`
	RightWords   []string `js:"rightWords"`
}

func Autocomplete(buildInfo *core.BuildInfo, cfg *AutoCompleteConfig) ([]string, error) {
	indexToComplete := int64(len(cfg.LeftWords) + 2)

	words := append(cfg.LeftWords, cfg.SelectedWord)
	words = append(words, cfg.RightWords...)

	completeCommand := []string{
		"autocomplete",
		"complete",
		"zsh",
		strconv.FormatInt(indexToComplete, 10),
		"scw",
	}

	completeCommand = append(completeCommand, words...)

	resp, err := Run(buildInfo, &RunConfig{
		JWT: cfg.JWT,
	}, completeCommand)
	if err != nil {
		return nil, fmt.Errorf("error running complete command: %w", err)
	}

	if resp.ExitCode != 0 {
		return nil, fmt.Errorf("invalid exit code %d: %s", resp.ExitCode, resp.Stderr)
	}

	suggestions := strings.Fields(resp.Stdout)

	return suggestions, nil
}

func AutocompleteWithBuildInfo(
	buildInfo *core.BuildInfo,
) func(cfg *AutoCompleteConfig) ([]string, error) {
	return func(cfg *AutoCompleteConfig) ([]string, error) {
		return Autocomplete(buildInfo, cfg)
	}
}
