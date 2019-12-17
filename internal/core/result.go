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
}

func (s *SuccessResult) MarshalHuman() (string, error) {
	message := s.getMessage()
	if !strings.HasSuffix(message, ".") {
		message += "."
	}
	message = strcase.TitleFirstWord(message)
	return "✅ " + terminal.Style(s.getMessage(), color.FgGreen), nil
}

func (s *SuccessResult) MarshalJSON() ([]byte, error) {
	type tmpRes struct {
		Message string `json:"message"`
	}
	return json.Marshal(&tmpRes{Message: s.getMessage()})
}

func (s *SuccessResult) getMessage() string {
	if s.Message != "" {
		return s.Message
	}
	return "Success"
}
