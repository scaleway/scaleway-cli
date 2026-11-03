package core

import (
	"github.com/scaleway/scaleway-cli/v2/core/human"
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
	Title       string
	FieldName   string
	HideIfEmpty bool
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
			Title:       section.Title,
			FieldName:   section.FieldName,
			HideIfEmpty: section.HideIfEmpty,
		})
	}
	opt.Title = v.Title

	return opt
}

// ListViewBuilder returns an override function that sets an explicit view on a
// list command: the displayed columns are the given field names, stable
// regardless of the model fields.
func ListViewBuilder(fieldNames ...string) func(c *Command) *Command {
	return func(c *Command) *Command {
		view := &View{Fields: make([]*ViewField, 0, len(fieldNames))}
		for _, fieldName := range fieldNames {
			view.Fields = append(view.Fields, &ViewField{FieldName: fieldName})
		}
		c.View = view

		return c
	}
}
