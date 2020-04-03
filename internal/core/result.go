package core

import (
	"encoding/json"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

type SuccessResult struct {
	Message string
	Details string
}

func (s *SuccessResult) MarshalHuman() (string, error) {
	message := s.getMessage()
	if !strings.HasSuffix(message, ".") {
		message += "."
	}

	message = strcase.TitleFirstWord(message)
	message = "✅ " + terminal.Style(message, color.FgGreen)

	if s.Details != "" {
		message += s.Details
	}

	return message, nil
}

func (s *SuccessResult) MarshalJSON() ([]byte, error) {
	type tmpRes struct {
		Message string `json:"message"`
		Details string `json:"details"`
	}
	return json.Marshal(&tmpRes{Message: s.getMessage()})
}

func (s *SuccessResult) getMessage() string {
	if s.Message != "" {
		return s.Message
	}
	return "Success"
}
