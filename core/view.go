package core

import (
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

// View hydrates humanv2.Layout
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

func (v *View) getHumanLayout() *humanv2.Layout {
	if v == nil {
		return nil
	}

	main := humanv2.Section{Path: ""}
	for _, field := range v.Fields {
		main.Fields = append(main.Fields, humanv2.Field{
			Path:  field.FieldName,
			Label: field.Label,
		})
	}
	sections := []humanv2.Section{main}
	for _, section := range v.Sections {
		sections = append(sections, humanv2.Section{
			Title:     section.Title,
			Path:      section.FieldName,
			Omittable: section.HideIfEmpty,
		})
	}

	return &humanv2.Layout{
		Title:    v.Title,
		Sections: sections,
	}
}
