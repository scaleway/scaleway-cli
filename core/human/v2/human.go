package human

import (
	"fmt"
	"reflect"
	"slices"
	"strings"
	"sync"

	"github.com/kr/pretty"
	"github.com/scaleway/scaleway-cli/v2/internal/tabwriter"
)

// Layout describes the structure of the rendered data. Sections define
// blocks that will be rendered individually from the data in their path,
// the remaining fields being collected into the main section, defined with
// an empty path and no title by default. Sections order is conserved, with
// the main section being prepended if it isn't found in the list.
type Layout struct {
	Sections []Section
}

type Section struct {
	// Path is the dotted-path to the field composing the section in the data.
	// The main section is denoted with an empty path.
	Path string

	// Title is the string for the section's header. It will be omitted if empty.
	Title string

	// Omittable controls whether the section should be rendered if empty.
	Omittable bool
}

// ResolveSections return the final list of sections, with any default applied
// like the default section.
func (l Layout) ResolveSections() []Section {
	if slices.ContainsFunc(l.Sections, func(s Section) bool { return s.Path == "" }) {
		return l.Sections
	}

	return append([]Section{{Path: "", Title: ""}}, l.Sections...)
}

var layouts sync.Map

// SetLayout stores a custom layout for a particular type. This custom layout
// will override the default one, and be eventually overridden by the render
// options themselves. This is useful because we don't want to have to
// maintain the layouts in the SDK itself, nor wrap the types to implement
// an interface.
func SetLayout(t reflect.Type, l Layout) {
	layouts.Store(t, l)
}

// Render data into a specific layout.
func Render(data any, layout Layout) (string, error) {
	sections := layout.ResolveSections()
	value := reflect.ValueOf(data)

	switch value.Kind() {
	// If the value is a slice or array, we ensure that the layout is a single
	// pathless section and write it.
	case reflect.Slice, reflect.Array:
		if len(sections) != 1 {
			return "", fmt.Errorf("found %d sections for %T payload", len(sections), data)
		}

		if sections[0].Path != "" {
			return "", fmt.Errorf("unknown path %q for %T payload", sections[0].Path, data)
		}

		var b strings.Builder
		err := writeSection(&b, sections[0], data)
		if err != nil {
			return "", fmt.Errorf("rendering section: %w", err)
		}

		return b.String(), nil

	// If the value is a struct, we initialize a map of the data fields, then
	// split them by section.
	case reflect.Struct:
		fields := make(map[string]any, value.NumField())
		for i := range value.NumField() {
			fields[value.Type().Field(i).Name] = value.Field(i).Interface()
		}

		fieldsBySection := make(map[int]any, len(sections))
		for i, section := range sections {
			if section.Path == "" {
				fieldsBySection[i] = fields
			}
			if field, ok := fields[section.Path]; ok {
				fieldsBySection[i] = field
				delete(fields, section.Path)
			} else {
				return "", fmt.Errorf("path %q not found in %T", section.Path, data)
			}
		}

		var b strings.Builder
		for i, section := range sections {
			err := writeSection(&b, section, fieldsBySection[i])
			if err != nil {
				return "", fmt.Errorf("rendering section %d: %w", i, err)
			}
			b.WriteRune('\n')
		}

		return b.String(), nil
	}

	return "", nil
}

func writeSection(b *strings.Builder, section Section, data any) error {
	if section.Title != "" {
		b.WriteString(section.Title)
		b.WriteRune(':')
		b.WriteRune('\n')
	}

	value := reflect.ValueOf(data)

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		elem := value.Type().Elem()

		w := tabwriter.NewWriter(b, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)
		for fieldIndex := range elem.NumField() {
			if fieldIndex != 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, elem.Field(fieldIndex).Name)
		}
		fmt.Fprint(w, "\n")
		for index := range value.Len() {
			for fieldIndex := range value.Index(index).NumField() {
				if fieldIndex != 0 {
					fmt.Fprint(w, "\t")
				}
				fmt.Fprint(w, value.Index(index).Field(fieldIndex).Interface())
			}
			fmt.Fprint(w, "\n")
		}
		w.Flush()

	default:
		b.WriteString(pretty.Sprint(data))
	}

	return nil
}
