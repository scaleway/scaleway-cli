package core

import (
	"encoding/json"
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
