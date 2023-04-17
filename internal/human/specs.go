package human

import (
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// MarshalOpt is hydrated by core.View
type MarshalOpt struct {
	Title    string
	Fields   []*MarshalFieldOpt
	Sections []*MarshalSection

	// Is set to true if we are marshaling a table cell
	TableCell bool

	// DisableShrinking will disable columns shrinking based on terminal size
	DisableShrinking bool
}

type MarshalFieldOpt struct {
	FieldName string
	Label     string
}

// MarshalSection describes a section to build from a given struct.
// When marshalling, this section is shown under the main struct section.
type MarshalSection struct {
	FieldName   string
	Title       string
	HideIfEmpty bool
}

func (s *MarshalFieldOpt) getLabel() string {
	if s.Label != "" {
		return s.Label
	}

	label := s.FieldName
	label = strings.ReplaceAll(label, ".", " ")
	label = strcase.ToBashArg(label)
	label = strings.ReplaceAll(label, "-", " ")
	label = strings.ToUpper(label)
	return label
}
