package core

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

type SuccessResult struct {
	Message  string
	Details  string
	Resource string
	Verb     string
	Empty    bool
	// Used to pass resource to an AfterFunc on success
	TargetResource any
}

// This type can be return by a command that need to output specific content on stdout directly.
// When a command return this type, default printer will not be used and bytes will be directly print on stdout.
type RawResult []byte

var standardSuccessMessages = map[string]string{
	"delete": "%s has been successfully deleted.",
}

func (s *SuccessResult) MarshalHuman() (string, error) {
	if s.Empty {
		return "", nil
	}
	message := s.getMessage()
	if !strings.HasSuffix(message, ".") {
		message += "."
	}

	message = strcase.TitleFirstWord(message)
	message = "✅ " + terminal.Style(message, color.FgGreen)

	if s.Details != "" {
		message += "\n" + interactive.Indent(s.Details, 2)
	}

	return message, nil
}

func (s *SuccessResult) MarshalJSON() ([]byte, error) {
	if s.Empty {
		type emptyRes struct{}

		return json.Marshal(&emptyRes{})
	}
	type tmpRes struct {
		Message string `json:"message"`
		Details string `json:"details"`
	}

	return json.Marshal(&tmpRes{
		Message: s.getMessage(),
		Details: s.Details,
	})
}

func (s *SuccessResult) getMessage() string {
	if s.Message != "" {
		return s.Message
	}

	if messageTemplate, exists := standardSuccessMessages[s.Verb]; exists && s.Resource != "" {
		return fmt.Sprintf(messageTemplate, s.Resource)
	}

	return "Success"
}

type MultiResults []any

func (mr MultiResults) MarshalHuman() (string, error) {
	strs := []string(nil)
	for _, r := range mr {
		str, err := human.Marshal(r, nil)
		if err != nil {
			return "", err
		}
		strs = append(strs, str)
	}

	return strings.Join(strs, "\n"), nil
}
