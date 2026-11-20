package core

import (
	"encoding/json/jsontext"
	"encoding/json/v2"
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core/human"
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

func (s *CliError) MarshalHuman() (string, error) {
	if s.Empty {
		return "", nil
	}
	sections := []string(nil)
	if s.Err != nil {
		humanError := s.Err
		if s.Message != "" {
			humanError = fmt.Errorf("%s", s.Message)
		}
		str, err := human.Marshal(humanError, nil)
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	if s.Details != "" {
		str, err := human.Marshal(human.Capitalize(s.Details), &human.MarshalOpt{Title: "Details"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	if s.Hint != "" {
		str, err := human.Marshal(human.Capitalize(s.Hint), &human.MarshalOpt{Title: "Hint"})
		if err != nil {
			return "", err
		}
		sections = append(sections, str)
	}

	return strings.Join(sections, "\n\n"), nil
}

// jsonError returns the JSON representation of e, reproducing the encoding/json v1
// behavior: an error is marshaled through its exported fields or MarshalJSON, and an
// error with no exported fields (e.g. errors.errorString) marshals to {}.
func jsonError(e error) string {
	b, err := json.Marshal(e, jsontext.EscapeForHTML(true))
	if err != nil {
		return "{}"
	}

	return string(b)
}

// jsonQuoted returns s as a JSON string literal.
func jsonQuoted(s string) (string, error) {
	b, err := json.Marshal(s, jsontext.EscapeForHTML(true))

	return string(b), err
}

func (s *CliError) MarshalJSON() ([]byte, error) {
	if s.Empty {
		return []byte("{}"), nil
	}

	message := s.Message
	if message == "" && s.Err != nil {
		message = s.Err.Error()
	}

	// encoding/json v2 refuses to marshal errors with no exported fields and its
	// omitempty drops them, so the fields are assembled by hand to keep the v1 wire format.
	parts := make([]string, 0, 4)
	if message != "" {
		q, err := jsonQuoted(message)
		if err != nil {
			return nil, err
		}

		parts = append(parts, `"message":`+q)
	}
	if s.Err != nil {
		parts = append(parts, `"error":`+jsonError(s.Err))
	}
	if s.Details != "" {
		q, err := jsonQuoted(s.Details)
		if err != nil {
			return nil, err
		}

		parts = append(parts, `"details":`+q)
	}
	if s.Hint != "" {
		q, err := jsonQuoted(s.Hint)
		if err != nil {
			return nil, err
		}

		parts = append(parts, `"hint":`+q)
	}

	return []byte("{" + strings.Join(parts, ",") + "}"), nil
}
