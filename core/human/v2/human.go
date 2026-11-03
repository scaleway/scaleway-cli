package human

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

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

	// Fields limits the section to the selected dotted field paths: for a
	// table it selects the columns, for a struct or map it selects the lines.
	// Empty means all fields.
	Fields []string
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

// isFormattable reports whether t should be rendered inline by Format rather
// than as a composite table.
func isFormattable(t reflect.Type) bool {
	// A formatter at any pointer level escapes through Format, which checks
	// each level in its deref loop, so mirror it here.
	for t.Kind() == reflect.Pointer {
		if t.Implements(reflect.TypeFor[Formatter]()) {
			return true
		}
		t = t.Elem()
	}
	switch {
	// A struct is a table, unless it escapes through Formatter.
	case t.Kind() == reflect.Struct:
		return t.Implements(reflect.TypeFor[Formatter]())
	// A slice of structs breaks into a table; a scalar slice is inline.
	case t.Kind() == reflect.Slice || t.Kind() == reflect.Array:
		elem := t.Elem()
		for elem.Kind() == reflect.Pointer {
			elem = elem.Elem()
		}

		return elem.Kind() != reflect.Struct
	// A map is a key-value table.
	case t.Kind() == reflect.Map:
		return false
	default:
		return true
	}
}

// Render data into a specific layout.
func Render(data any, layout Layout) (string, error) {
	// Dereference pointers until the underlying value is reached.
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return "-", nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return "-", nil
	}
	data = value.Interface()

	// Fall back to a per-type registered layout when the caller gave none,
	// then apply the default pathless section if needed.
	if reflect.ValueOf(layout).IsZero() {
		if l, ok := layouts.Load(value.Type()); ok {
			layout = l.(Layout)
		}
	}
	sections := layout.ResolveSections()

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
		_, err := writeSection(&b, sections[0], data, nil)
		if err != nil {
			return "", fmt.Errorf("rendering section: %w", err)
		}

		return b.String(), nil

	// If the value is a struct, we resolve the sections from the data, then
	// split them by section. claimed tracks the dotted paths rendered as their
	// own section so the flattening of the main block skips them, keeping
	// declaration order.
	case reflect.Struct:
		// claimed tracks the dotted paths rendered as their own section so
		// the flattening of the main block skips them, keeping declaration
		// order. Paths are validated when they are resolved below.
		claimed := make(map[string]bool, len(sections))
		for _, section := range sections {
			if section.Path != "" {
				claimed[section.Path] = true
			}
		}

		var b strings.Builder
		for i, section := range sections {
			sectionData := data
			if section.Path != "" {
				var err error
				sectionData, err = valueAtPath(data, section.Path)
				if err != nil {
					return "", fmt.Errorf("path %q not found in %T: %w", section.Path, data, err)
				}
			}
			wrote, err := writeSection(&b, section, sectionData, claimed)
			if err != nil {
				return "", fmt.Errorf("rendering section %d: %w", i, err)
			}
			// Only separate sections that wrote content, so a skipped omittable
			// section leaves no orphan blank line.
			if wrote {
				b.WriteRune('\n')
			}
		}

		return strings.TrimSpace(b.String()), nil
	}

	return Format(data)
}

// writeSection data to the given string builder, with a map of the already
// claimed paths so we can skip the extracted sections from being rendered in
// the default section. This is necessary because the section data must be
// rendered in declaration order, so we can't convert the original data to
// a map without having to add additional complexity. It returns whether
// anything was written.
func writeSection(b *strings.Builder, section Section, data any, claimed map[string]bool) (bool, error) {
	value := reflect.ValueOf(data)
	if (!value.IsValid() || value.IsZero()) && section.Omittable {
		return false, nil
	}

	if section.Title != "" {
		b.WriteString(section.Title)
		b.WriteRune(':')
		b.WriteRune('\n')
	}

	if !value.IsValid() || isFormattable(value.Type()) {
		str, err := Format(data)
		if err != nil {
			return false, err
		}
		b.WriteString(str)

		return true, nil
	}

	switch value.Kind() {
	// Slices and arrays renders as table.
	case reflect.Slice, reflect.Array:
		itemType := value.Type().Elem()
		for itemType.Kind() == reflect.Pointer {
			itemType = itemType.Elem()
		}

		// Resolve the selected fields. If no field was selected, initialize the list
		// with the item's fields.
		fields := section.Fields
		if len(section.Fields) == 0 {
			fields = make([]string, itemType.NumField())
			for fieldIndex := range itemType.NumField() {
				fields[fieldIndex] = itemType.Field(fieldIndex).Name
			}
		}

		w := tabwriter.NewWriter(b, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)
		for fieldIndex, field := range fields {
			if fieldIndex != 0 {
				fmt.Fprint(w, "\t")
			}
			fmt.Fprint(w, field)
		}
		fmt.Fprint(w, "\n")

		for index := range value.Len() {
			item := value.Index(index)

			// Dereference the item's value until it is either nil or not a pointer so
			// we can access the underlying's fields.
			for (item.Kind() == reflect.Pointer || item.Kind() == reflect.Interface) && !item.IsNil() {
				item = item.Elem()
			}
			isNil := item.Kind() == reflect.Pointer || item.Kind() == reflect.Interface

			for fieldIndex, field := range fields {
				if fieldIndex != 0 {
					fmt.Fprint(w, "\t")
				}

				if isNil {
					fmt.Fprint(w, "-")
					continue
				}

				cell, err := valueAtPath(item.Interface(), field)
				if err != nil {
					return false, err
				}

				str, err := Format(cell)
				if err != nil {
					return false, err
				}

				fmt.Fprint(w, str)
			}
			fmt.Fprint(w, "\n")
		}
		w.Flush()

		return true, nil

	// A struct and a map render the same: one walk flattens composites into
	// dotted rows, renders formattable leaves in place, and skips paths claimed
	// by other sections. Maps are just walked as keyed leaves.
	case reflect.Map, reflect.Struct:
		w := tabwriter.NewWriter(b, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)
		err := walk(value, "", func(keys string, node reflect.Value) (bool, error) {
			if claimed[keys] {
				return false, nil
			}

			// Field selection keeps only the selected subtrees: a node above
			// one still descends, a branch without selection is pruned.
			if len(section.Fields) > 0 && keys != "" {
				under := false
				for _, f := range section.Fields {
					if keys == f || strings.HasPrefix(keys, f+".") {
						under = true
						break
					}
				}
				if !under {
					for _, f := range section.Fields {
						if strings.HasPrefix(f, keys+".") {
							return true, nil
						}
					}

					return false, nil
				}
			}

			// Scalars, formatter structs (and pointers to them) render in place.
			// Everything else descends into dotted leaves. isFormattable reports
			// scalar slices as inline for section bodies, but a walk always
			// flattens slices into indexed rows.
			t := node.Type()
			for t.Kind() == reflect.Pointer {
				t = t.Elem()
			}
			if isFormattable(node.Type()) && t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
				str, err := Format(node.Interface())
				if err != nil {
					return false, err
				}
				fmt.Fprintf(w, "%s\t%s\n", keys, str)

				return false, nil
			}

			return true, nil
		})
		if err != nil {
			return false, err
		}
		w.Flush()

		return true, nil
	}

	return false, nil
}

// walk represents value and its nodes to f; f installs its own walk logic:
// it's called for every node with its dotted-path keys (Sub.X, Strings.0,
// Map.key) and returns whether to keep descending after it. Nil pointers and
// interfaces yield no node at all.
func walk(value reflect.Value, keys string, f func(keys string, value reflect.Value) (bool, error)) error {
	if (value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface) && value.IsNil() {
		return nil
	}

	descend, err := f(keys, value)
	if err != nil {
		return err
	}
	if !descend {
		return nil
	}

	switch value.Kind() {
	case reflect.Pointer, reflect.Interface:
		return walk(value.Elem(), keys, f)

	case reflect.Struct:
		for i := range value.NumField() {
			name := value.Type().Field(i).Name
			err := walk(value.Field(i), joinKeys(keys, name), f)
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Slice, reflect.Array:
		for i := range value.Len() {
			err := walk(value.Index(i), joinKeys(keys, strconv.Itoa(i)), f)
			if err != nil {
				return err
			}
		}

		return nil

	case reflect.Map:
		mapKeys := value.MapKeys()
		slices.SortFunc(mapKeys, func(i, j reflect.Value) int {
			return strings.Compare(i.String(), j.String())
		})
		for _, mapKey := range mapKeys {
			err := walk(value.MapIndex(mapKey), joinKeys(keys, mapKey.String()), f)
			if err != nil {
				return err
			}
		}

		return nil

	default:
		// Leaf: f already decided whether to descend.
		return nil
	}
}

// joinKeys appends part to keys, dropping the separator at the root.
func joinKeys(keys, part string) string {
	if keys == "" {
		return part
	}

	return keys + "." + part
}

// valueAtPath returns the value at the given dotted path, dereferencing
// pointers and interfaces along the way.
func valueAtPath(data any, path string) (any, error) {
	value := reflect.ValueOf(data)

	for part := range strings.SplitSeq(path, ".") {
		for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
			if value.IsNil() {
				return nil, fmt.Errorf("path %q not found in %T: nil value", path, data)
			}
			value = value.Elem()
		}

		switch value.Kind() {
		case reflect.Struct:
			value = value.FieldByName(part)
		case reflect.Slice, reflect.Array:
			index, err := strconv.Atoi(part)
			if err != nil || index < 0 || index >= value.Len() {
				return nil, pathNotFound(path, data)
			}
			value = value.Index(index)
		case reflect.Map:
			value = value.MapIndex(reflect.ValueOf(part))
		default:
			return nil, pathNotFound(path, data)
		}
		if !value.IsValid() || !value.CanInterface() {
			return nil, pathNotFound(path, data)
		}
	}

	return value.Interface(), nil
}

// pathNotFound reports a dotted path component that does not resolve in data.
func pathNotFound(path string, data any) error {
	return fmt.Errorf("path %q not found in %T", path, data)
}
