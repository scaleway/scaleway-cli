package human

import (
	"fmt"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/internal/tabwriter"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

// Layout describes the structure of the rendered data. Sections define
// blocks that will be rendered individually from the data in their path,
// the remaining fields being collected into the main section, defined with
// an empty path and no title by default. Sections order is conserved, with
// the main section being prepended if it isn't found in the list.
type Layout struct {
	// Title, when set, is rendered on its own line above the whole output,
	// wrapping all sections like v1's MarshalOpt.Title.
	Title string

	// NoHeaders skips the header row of table blocks, like v1's
	// MarshalOpt.NoHeaders.
	NoHeaders bool

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
	Fields []Field
}

// Field selects one column of a table or one line of a struct or map.
type Field struct {
	// Path is the dotted path of the field in the data.
	Path string

	// Label overrides the derived column header. Empty derives it from Path.
	Label string
}

// getLabel returns the field's display label: the explicit Label if set,
// otherwise the Path derived like v1 (bash-arg spaced and uppercased).
func (f Field) getLabel() string {
	if f.Label != "" {
		return f.Label
	}

	return deriveLabel(f.Path)
}

// deriveLabel renders a path as an ALL-CAPS column header, mirroring v1's
// getDefaultFieldsOpt label derivation.
func deriveLabel(path string) string {
	label := strings.ReplaceAll(path, ".", " ")
	label = strcase.ToBashArg(label)
	label = strings.ReplaceAll(label, "-", " ")

	return strings.ToUpper(label)
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

// Block is the resolved content of one section, free of any rendering
// concern. It is either a table (Header non-empty, Rows are cells) or a set
// of lines (Header empty, Rows are [key, value] pairs or a single cell).
type Block struct {
	Title  string
	Header []string
	Rows   [][]string

	// NoHeaders skips the header row when drawing a table block.
	NoHeaders bool
}

// Empty reports whether block draws nothing at all.
func (b Block) Empty() bool {
	return b.Title == "" && len(b.Header) == 0 && len(b.Rows) == 0
}

// Render resolves data into its layout blocks and draws them, separating the
// blocks with a blank line. Pointers are dereferenced at the top level, and a
// zero layout falls back to the layout registered for the concrete type of
// data.
func Render(data any, layout Layout) (string, error) {
	// Dereference pointers until the underlying value is reached, so the nil
	// shortcut and the concrete type keying the fallback apply to the same
	// value.
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

	if reflect.ValueOf(layout).IsZero() {
		if l, ok := layouts.Load(value.Type()); ok {
			layout = l.(Layout)
		}
	}

	blocks, err := layout.Resolve(data)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	for _, block := range blocks {
		if block.Empty() {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(block.Draw())
	}

	if layout.Title != "" {
		return layout.Title + ":\n" + b.String(), nil
	}

	return b.String(), nil
}

// Resolve splits data into the content blocks of the layout, without any
// rendering. It is exposed so tests and a future layout engine can reuse the
// resolution part of Render. Pointers are dereferenced first, so direct
// callers can pass them too; leaves and nil values fold into single-cell
// blocks.
func (l Layout) Resolve(data any) ([]Block, error) {
	// Dereference pointers until the underlying value is reached.
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		if value.IsNil() {
			return []Block{{Rows: [][]string{{"-"}}}}, nil
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return []Block{{Rows: [][]string{{"-"}}}}, nil
	}
	data = value.Interface()

	// Apply the default pathless section if needed.
	sections := l.ResolveSections()

	switch value.Kind() {
	// A slice or array is a single pathless section, rendered as a table.
	case reflect.Slice, reflect.Array:
		if len(sections) != 1 {
			return nil, fmt.Errorf("found %d sections for %T payload", len(sections), data)
		}

		if sections[0].Path != "" {
			return nil, fmt.Errorf("unknown path %q for %T payload", sections[0].Path, data)
		}

		block, err := sections[0].Resolve(data, nil)
		if err != nil {
			return nil, fmt.Errorf("rendering section: %w", err)
		}
		if block.Empty() {
			return nil, nil
		}

		block.NoHeaders = l.NoHeaders

		return []Block{block}, nil

	// A struct resolves into one block per section. claimed tracks the dotted
	// paths rendered as their own section so the flattening of the main block
	// skips them, keeping declaration order. Paths are validated here, when
	// they are resolved.
	case reflect.Struct:
		claimed := make(map[string]bool, len(sections))
		for _, section := range sections {
			if section.Path != "" {
				claimed[section.Path] = true
			}
		}

		blocks := make([]Block, 0, len(sections))
		for i, section := range sections {
			sectionData := data
			if section.Path != "" {
				var err error
				sectionData, err = valueAtPath(data, section.Path)
				if err != nil {
					return nil, fmt.Errorf("path %q not found in %T: %w", section.Path, data, err)
				}
			}
			block, err := section.Resolve(sectionData, claimed)
			if err != nil {
				return nil, fmt.Errorf("rendering section %d: %w", i, err)
			}
			if !block.Empty() {
				block.NoHeaders = l.NoHeaders
				blocks = append(blocks, block)
			}
		}

		return blocks, nil
	}

	// A leaf (a scalar or any other kind) renders inline through its formatter.
	str, err := Format(data)
	if err != nil {
		return nil, err
	}

	return []Block{{Rows: [][]string{{str}}}}, nil
}

// Resolve renders section into a block from data, without drawing. claimed
// tracks the dotted paths resolved as their own section so the flattening of
// the default block skips them. An omitted or empty section resolves to an
// empty block, a slice or array to a table, and a struct or map to a list of
// keyed rows.
func (s Section) Resolve(data any, claimed map[string]bool) (Block, error) {
	value := reflect.ValueOf(data)
	if (!value.IsValid() || value.IsZero()) && s.Omittable {
		return Block{}, nil
	}

	block := Block{Title: s.Title}
	if block.Title == "" && s.Path != "" {
		block.Title = deriveTitle(s.Path)
	}

	if !value.IsValid() || isFormattable(value.Type()) {
		str, err := Format(data)
		if err != nil {
			return Block{}, err
		}
		block.Rows = [][]string{{str}}

		return block, nil
	}

	switch value.Kind() {
	case reflect.Slice, reflect.Array:
		return s.resolveTable(value, block)
	case reflect.Map, reflect.Struct:
		return s.resolveList(value, block, claimed)
	}

	return block, nil
}

// resolveTable renders a slice or array section as a table. With no selected
// Fields the item's fields become the columns; values are formatted per cell.
func (s Section) resolveTable(value reflect.Value, block Block) (Block, error) {
	itemType := value.Type().Elem()
	for itemType.Kind() == reflect.Pointer {
		itemType = itemType.Elem()
	}

	// Resolve the selected fields. If no field was selected, initialize the
	// list with the item's fields.
	fields := s.Fields
	if len(s.Fields) == 0 {
		fields = make([]Field, itemType.NumField())
		for fieldIndex := range itemType.NumField() {
			fields[fieldIndex] = Field{Path: itemType.Field(fieldIndex).Name}
		}
	}

	if itemType.Kind() == reflect.Struct {
		if err := validateFields(itemType, fields); err != nil {
			return Block{}, err
		}
	}

	block.Header = make([]string, len(fields))
	for fieldIndex, field := range fields {
		block.Header[fieldIndex] = field.getLabel()
	}

	for index := range value.Len() {
		item := value.Index(index)

		// Dereference the item's value until it is either nil or not a
		// pointer so we can access the underlying's fields.
		for (item.Kind() == reflect.Pointer || item.Kind() == reflect.Interface) && !item.IsNil() {
			item = item.Elem()
		}
		isNil := item.Kind() == reflect.Pointer || item.Kind() == reflect.Interface

		row := make([]string, len(fields))
		for fieldIndex, field := range fields {
			if isNil {
				row[fieldIndex] = "-"
				continue
			}

			cell, err := valueAtPath(item.Interface(), field.Path)
			if err != nil {
				// Field paths are pre-validated, so only per-item missing values
				// (nil intermediate, absent map key) land here: degrade to "-"
				// like v1 instead of failing the whole table.
				row[fieldIndex] = "-"
				continue
			}

			str, err := Format(cell)
			if err != nil {
				return Block{}, err
			}
			row[fieldIndex] = str
		}
		block.Rows = append(block.Rows, row)
	}

	return block, nil
}

// resolveList renders a struct or map section into dotted keyed rows: one walk
// flattens composite values into dotted keys, formats formattable leaves in
// place, and skips paths claimed by other sections. Maps are just walked as
// keyed leaves.
func (s Section) resolveList(value reflect.Value, block Block, claimed map[string]bool) (Block, error) {
	valueType := value.Type()
	for valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}
	if valueType.Kind() == reflect.Struct {
		if err := validateFields(valueType, s.Fields); err != nil {
			return Block{}, err
		}
	}

	err := walk(value, "", func(keys string, node reflect.Value) (bool, error) {
		if claimed[keys] {
			return false, nil
		}

		// Field selection keeps only the selected subtrees: a node above
		// one still descends, a branch without selection is pruned.
		if len(s.Fields) > 0 && keys != "" {
			under := false
			for _, f := range s.Fields {
				if keys == f.Path || strings.HasPrefix(keys, f.Path+".") {
					under = true
					break
				}
			}
			if !under {
				for _, f := range s.Fields {
					if strings.HasPrefix(f.Path, keys+".") {
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
			block.Rows = append(block.Rows, []string{keys, str})

			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return Block{}, err
	}

	return block, nil
}

// validateFields returns an UnknownFieldError for the first selected field
// whose path does not resolve through fields of t.
func validateFields(t reflect.Type, fields []Field) error {
	for _, field := range fields {
		if !hasField(t, field.Path) {
			return &UnknownFieldError{FieldName: field.Path, ValidFields: listFields(t)}
		}
	}

	return nil
}

// hasField reports whether the dotted path resolves through t's fields.
func hasField(t reflect.Type, path string) bool {
	for part := range strings.SplitSeq(path, ".") {
		for t.Kind() == reflect.Pointer {
			t = t.Elem()
		}
		if t.Kind() != reflect.Struct {
			return false
		}

		field, ok := t.FieldByName(part)
		if !ok {
			return false
		}
		t = field.Type
	}

	return true
}

// listFields returns the selectable dotted field paths of t: structs recurse,
// other leaves (scalars, slices, maps) contribute their own path.
func listFields(t reflect.Type) []string {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	res := []string(nil)
	collectFields(t, "", &res)

	return res
}

func collectFields(t reflect.Type, prefix string, res *[]string) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() == reflect.Struct {
		for field := range t.Fields() {
			if field.PkgPath != "" {
				continue
			}

			parent := prefix
			if !field.Anonymous {
				parent = joinKeys(prefix, field.Name)
			}
			collectFields(field.Type, parent, res)
		}

		return
	}

	if prefix != "" {
		*res = append(*res, prefix)
	}
}

// deriveTitle renders a section's dotted path as a titled section header,
// mirroring v1's marshalSection title derivation.
func deriveTitle(path string) string {
	title := strings.ReplaceAll(strcase.ToBashArg(path), "-", " ")

	return cases.Title(language.English).String(strings.ReplaceAll(title, ".", " - "))
}

// Draw renders the block to text: an optional title line followed by its
// lines or table, aligned to columns. It is the layout hook of the package:
// column alignment today, terminal-width shrinking and lipgloss styling
// later, so swapping the layout engine only touches this method. An empty
// block draws nothing, and accumulated whitespace is trimmed.
func (b Block) Draw() string {
	var s strings.Builder
	if b.Title != "" {
		s.WriteString(b.Title)
		s.WriteRune(':')
		s.WriteRune('\n')
	}

	w := tabwriter.NewWriter(&s, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)
	if len(b.Header) > 0 && !b.NoHeaders {
		fmt.Fprintln(w, strings.Join(b.Header, "\t"))
	}
	for _, row := range b.Rows {
		fmt.Fprintln(w, strings.Join(row, "\t"))
	}
	w.Flush()

	return strings.TrimSpace(s.String())
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
