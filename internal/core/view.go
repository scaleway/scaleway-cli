package core

import (
	"github.com/scaleway/scaleway-cli/internal/human"
)

// View hydrates human.MarshalOpt
type View struct {
	Title    string
	Fields   []*ViewField
	Sections []*ViewSection
}

type ViewField struct {

	// Label is the string displayed as header or key for a field.
	Label string

	// FieldName is the key used to retrieve the value from a field (or nested field) of a structure.
	FieldName string
}

type ViewSection struct {
	Title     string
	FieldName string
}

func (v *View) getHumanMarshalerOpt() *human.MarshalOpt {
	if v == nil {
		return nil
	}
	opt := &human.MarshalOpt{}
	for _, field := range v.Fields {
		opt.Fields = append(opt.Fields, &human.MarshalFieldOpt{
			FieldName: field.FieldName,
			Label:     field.Label,
		})
	}
	for _, section := range v.Sections {
		opt.Sections = append(opt.Sections, &human.MarshalSection{
			Title:     section.Title,
			FieldName: section.FieldName,
		})
	}
	opt.Title = v.Title
	return opt
}
