package core

import (
	"encoding/json"
	"strings"
	"unicode"

	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
)

// CliError is an all-in-one error structure that can be used in commands to return useful errors to the user.
// CliError implements JSON and human marshaler for a smooth experience.
type CliError struct {
	// The original error that triggers this CLI error.
	// The Err.String() will be print in red to the user.
	Err error

	// Message allow to override the red message shown to the use.
	// By default, we will use Err.String() but in same case you may want to keep Err
	// to avoid losing detail in json output.
	Message string

	Details string
	Hint    string

	// Code allows to return a specific error code from the main binary.
	Code int

	// Empty tells the marshaler to not print any message for the error
	Empty bool
}

func (s *CliError) Error() string {
	return s.Err.Error()
}

// FormatHuman lets humanv2 render a CliError as a leaf: a value receiver so
// the pointer is dereferenced before isFormattable sees the value type.
func (s CliError) FormatHuman() (string, error) {
	if s.Empty {
		return "", nil
	}
	sections := []string(nil)
	if s.Err != nil {
		message := s.Err.Error()
		if s.Message != "" {
			message = s.Message
		}
		sections = append(sections, terminal.Style(capitalize(message), color.FgRed))
	}

	if s.Details != "" {
		str, err := humanv2.Render(capitalize(s.Details), humanv2.Layout{Title: "Details"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	if s.Hint != "" {
		str, err := humanv2.Render(capitalize(s.Hint), humanv2.Layout{Title: "Hint"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	return strings.Join(sections, "\n\n"), nil
}

func (s *CliError) MarshalJSON() ([]byte, error) {
	if s.Empty {
		type emptyRes struct{}

		return json.Marshal(&emptyRes{})
	}

	message := s.Message
	if message == "" && s.Err != nil {
		message = s.Err.Error()
	}

	type tmpRes struct {
		Message string `json:"message,omitempty"`
		Error   error  `json:"error,omitempty"`
		Details string `json:"details,omitempty"`
		Hint    string `json:"hint,omitempty"`
	}

	return json.Marshal(&tmpRes{
		Message: message,
		Error:   s.Err,
		Details: s.Details,
		Hint:    s.Hint,
	})
}

// capitalize returns the given string with a first character in uppercase.
func capitalize(s string) string {
	words := strings.Fields(s)
	if len(words) > 0 && strings.Contains(words[0], "-") {
		// First word is probably an id, don't capitalize
		return s
	}

	for i, c := range s {
		return string(unicode.ToUpper(c)) + s[i+1:]
	}

	return ""
}
