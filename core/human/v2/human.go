package humanv2

import (
	"bytes"
	"fmt"
	"reflect"
	"slices"
	"sort"
	"strconv"
	"strings"
	"sync"

	"github.com/scaleway/scaleway-cli/v2/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/strcase"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

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

	if reflect.ValueOf(layout).IsZero() {
		if l, ok := layouts.Load(value.Type()); ok {
			layout = l.(Layout)
		}
	}

	// Pass the raw data through: Resolve checks isLeaf on the pointer type,
	// whose error and Formatter methods are often pointer-receiver.
	blocks, err := layout.Resolve(data)
	if err != nil {
		return "", err
	}

	var b strings.Builder
	// Width-driven shrinking only applies to tables and only when the output
	// is a terminal; disabling it (wide output) or a missing width keeps
	// every column.
	width := 0
	if !layout.NoShrinking {
		if w := terminal.GetWidth(); w > 0 {
			width = w
		}
	}

	for _, block := range blocks {
		if block.Empty() {
			continue
		}
		if b.Len() > 0 {
			b.WriteString("\n\n")
		}
		b.WriteString(block.Draw(width, layout.ColPadding))
	}

	if layout.Title != "" {
		return layout.Title + ":\n" + b.String(), nil
	}

	return b.String(), nil
}

// Layout describes the structure of the rendered data. Sections define
// blocks that will be rendered individually from the data in their path,
// the remaining fields being collected into the main section, defined with
// an empty path and no title by default. Sections order is conserved, with
// the main section being prepended if it isn't found in the list.
type Layout struct {
	// Title, when set, is rendered on its own line above the whole output.
	Title string

	// NoHeaders skips the header row of table blocks.
	NoHeaders bool

	// NoShrinking keeps every table column regardless of the terminal width.
	NoShrinking bool

	// ColPadding is the number of spaces between table columns, 0 means default.
	ColPadding int

	Sections []Section
}

// ResolveSections returns the layout's sections, validated for a single main
// section: at most one section may be the main (the Main flag or the pathless
// struct-root receiver), and one is prepended as the default pathless main
// when the layout declares none. All mainless layouts get that default, so
// the main block is always rendered.
func (l Layout) ResolveSections() ([]Section, error) {
	main := -1
	for i, s := range l.Sections {
		if !s.isMain() {
			continue
		}
		if main >= 0 {
			return nil, fmt.Errorf("only one main section allowed, got %v", l.Sections)
		}
		main = i
	}

	if main >= 0 {
		return l.Sections, nil
	}

	return append([]Section{{Path: ""}}, l.Sections...), nil
}

// Resolve splits data into the content blocks of the layout, without any
// rendering. It is exposed so tests and a future layout engine can reuse the
// resolution part of Render. Pointers are dereferenced first, so direct
// callers can pass them too; leaves and nil values fold into single-cell
// blocks.
func (l Layout) Resolve(data any) ([]Block, error) {
	return l.resolve(data, 0)
}

// maxLayoutNesting caps recursive section/layout descent, matching v1's
// recursive RegisterMarshalerFunc application to section payloads.
const maxLayoutNesting = 16

func (l Layout) resolve(data any, depth int) ([]Block, error) {
	if depth > maxLayoutNesting {
		return nil, fmt.Errorf("layout nesting too deep for %T", data)
	}

	// A formattable value on its raw type escapes to a single-cell block:
	// error and Formatter methods often hang on the pointer, so this has to
	// run before pointers are dereferenced.
	if reflect.ValueOf(data).IsValid() && isLeaf(reflect.TypeOf(data)) {
		str, err := Format(data)
		if err != nil {
			return nil, err
		}

		return []Block{{Rows: [][]string{{str}}}}, nil
	}

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

	// Validate the section list and guarantee a main section.
	sections, err := l.ResolveSections()
	if err != nil {
		return nil, err
	}

	switch value.Kind() {
	// A slice or array is a single pathless section, rendered as a table.
	case reflect.Slice, reflect.Array:
		if len(sections) != 1 {
			return nil, fmt.Errorf("found %d sections for %T payload", len(sections), data)
		}

		if sections[0].Path != "" {
			return nil, fmt.Errorf("unknown path %q for %T payload", sections[0].Path, data)
		}

		block, err := sections[0].Resolve(data, data, nil)
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
		return l.resolveStruct(sections, data, depth)
	}

	// A leaf (a scalar or any other kind) renders inline through its formatter.
	str, err := Format(data)
	if err != nil {
		return nil, err
	}

	return []Block{{Rows: [][]string{{str}}}}, nil
}

func (l Layout) resolveStruct(sections []Section, data any, depth int) ([]Block, error) {
	// The main section (Main flag, or the pathless default) renders the
	// struct root. With a non-empty path the field at that path becomes
	// the whole payload (v1's marshallNestedField), rendered flat and
	// titleless; it then must be the only section, since other paths are
	// resolved against the root struct. A pathless main leaves the whole
	// struct as the payload.
	if mainIdx := slices.IndexFunc(
		sections,
		func(s Section) bool { return s.isMain() },
	); mainIdx >= 0 &&
		sections[mainIdx].Path != "" {
		if len(sections) != 1 {
			return nil, fmt.Errorf(
				"main section with a path must be the only one, got %v",
				sections,
			)
		}
		sectionData, err := valueAtPath(data, sections[mainIdx].Path)
		if err != nil {
			return nil, fmt.Errorf(
				"path %q not found in %T: %w",
				sections[mainIdx].Path,
				data,
				err,
			)
		}
		data = sectionData
		sections = []Section{{}}
	}

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

		// A section over a type with its own registered layout renders
		// through that layout (v1 recursed into RegisterMarshalerFunc for
		// section payloads). It can expand into several blocks, e.g.
		// lbBackendLayout's "LB" section resolves lbLayout into
		// Lb/IPs/LB Instances.
		if section.Path != "" {
			if nested, ok := registeredLayout(sectionData); ok {
				if section.Omittable && isOmittable(reflect.ValueOf(sectionData)) {
					continue
				}
				nestedBlocks, err := nested.resolve(sectionData, depth+1)
				if err != nil {
					return nil, fmt.Errorf("rendering section %d: %w", i, err)
				}
				title := section.Title
				if title == "" {
					title = deriveTitle(section.Path)
				}
				for j, block := range nestedBlocks {
					if j == 0 {
						block.Title = title
					}
					block.NoHeaders = block.NoHeaders || l.NoHeaders
					if !block.Empty() {
						blocks = append(blocks, block)
					}
				}

				continue
			}
		}

		block, err := section.Resolve(sectionData, data, claimed)
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

var layouts sync.Map

// RegisterLayout stores a custom layout for a particular type. This custom layout
// will override the default one, and be eventually overridden by the render
// options themselves.
func RegisterLayout(t reflect.Type, l Layout) {
	layouts.Store(t, l)
}

type Section struct {
	// Path is the dotted-path to the field composing the section in the data.
	// The main section is denoted with an empty path: that section renders the
	// struct root. The Main flag marks the main explicitly and may carry a
	// Path to the field rendered as the whole root instead.
	Path string

	// Main marks this section as the layout's main section. With a non-empty
	// Path the field at that path is rendered as the whole payload, without a
	// main section or a derived title (v1's marshallNestedField); the section
	// then must be the only one. An empty Path, or no Main flag at all, makes
	// the section the default main receiving the struct root. At most one main
	// section is allowed in a layout.
	Main bool

	// Title is the string for the section's header. It will be omitted if empty.
	Title string

	// Omittable controls whether the section is rendered when empty.
	Omittable bool

	// Fields limits the section to the selected dotted field paths: for a
	// table it selects the columns, for a struct or map it selects the lines.
	// Empty means all fields.
	Fields []Field

	// Footer is free text drawn below the block, after a blank line.
	// FooterFunc, when set, computes it from the rendered payload instead.
	Footer     string
	FooterFunc func(any) (string, error)
}

// Resolve renders section into a block from data, without drawing. claimed
// tracks the dotted paths resolved as their own section so the flattening of
// the default block skips them. An omitted or empty section resolves to an
// empty block, a slice or array to a table, and a struct or map to a list of
// keyed rows.
func (s Section) Resolve(data, root any, claimed map[string]bool) (Block, error) {
	value := reflect.ValueOf(data)
	for value.Kind() == reflect.Pointer {
		value = value.Elem()
	}
	if s.Omittable && isOmittable(value) {
		return Block{}, nil
	}

	block := Block{Title: s.Title}
	if block.Title == "" && s.Path != "" {
		block.Title = deriveTitle(s.Path)
	}

	// isLeaf on the raw data type, not the dereferenced one: errors and
	// Formatter methods often hang on the pointer.
	var err error
	if !value.IsValid() || isLeaf(reflect.TypeOf(data)) {
		var str string
		str, err = Format(data)
		block.Rows = [][]string{{str}}
	} else {
		switch value.Kind() {
		case reflect.Slice, reflect.Array:
			itemType := value.Type().Elem()
			for itemType.Kind() == reflect.Pointer {
				itemType = itemType.Elem()
			}
			switch itemType.Kind() {
			case reflect.Struct:
				block.Header, block.Rows, err = s.resolveTable(value)
			default:
				// Non-struct item slices inline like v1's marshalSlice:
				// []string -> "[a b]", nil -> "[]".
				var str string
				str, err = Format(data)
				block.Rows = [][]string{{str}}
			}
		case reflect.Map, reflect.Struct:
			block.Rows, err = s.resolveList(value, claimed)
		default:
			return block, fmt.Errorf("unresolvable type %T", value.Interface())
		}
	}
	if err != nil {
		return Block{}, err
	}

	if s.FooterFunc != nil {
		block.Footer, err = s.FooterFunc(root)
	} else {
		block.Footer = s.Footer
	}
	if err != nil {
		return Block{}, err
	}

	return block, nil
}

// isMain reports whether the section is the layout's main: either explicitly
// marked Main, or pathless as the default struct-root receiver.
func (s Section) isMain() bool {
	return s.Main || s.Path == ""
}

// resolveList renders a struct or map section into dotted keyed rows: one walk
// flattens composite values into dotted keys, formats formattable leaves in
// place, and skips paths claimed by other sections. Maps are just walked as
// keyed leaves.
func (s Section) resolveList(value reflect.Value, claimed map[string]bool) ([][]string, error) {
	valueType := value.Type()
	for valueType.Kind() == reflect.Pointer {
		valueType = valueType.Elem()
	}
	if valueType.Kind() == reflect.Struct {
		if err := validateFields(valueType, s.Fields); err != nil {
			return nil, err
		}
	}

	var rows [][]string
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

		// Leaves render in place; everything else descends into dotted rows.
		// Slices are never leaves, so a walk always flattens them into
		// indexed rows; scalar-slice inlining lives in the section, cell and
		// layout resolvers.
		if isLeaf(node.Type()) {
			str, err := Format(node.Interface())
			if err != nil {
				return false, err
			}
			rows = append(rows, []string{keys, str})

			return false, nil
		}

		return true, nil
	})
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// resolveTable renders a slice or array section as a table. With no selected
// Fields the item's fields become the columns; values are formatted per cell.
func (s Section) resolveTable(value reflect.Value) ([]string, [][]string, error) {
	itemType := value.Type().Elem()
	for itemType.Kind() == reflect.Pointer {
		itemType = itemType.Elem()
	}

	// Resolve the selected fields. If no field was selected, initialize the
	// list with the item's fields.
	fields := s.Fields
	if len(s.Fields) == 0 {
		fields = defaultTableFields(itemType)
	}

	if itemType.Kind() == reflect.Struct {
		if err := validateFields(itemType, fields); err != nil {
			return nil, nil, err
		}
	}

	header := make([]string, len(fields))
	for fieldIndex, field := range fields {
		header[fieldIndex] = field.getLabel()
	}

	var rows [][]string
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

			str, err := formatCell(cell)
			if err != nil {
				return nil, nil, err
			}
			row[fieldIndex] = str
		}
		rows = append(rows, row)
	}

	return header, rows, nil
}

// Field selects one column of a table or one line of a struct or map.
type Field struct {
	// Path is the dotted path of the field in the data.
	Path string

	// Label overrides the derived column header. Empty derives it from Path.
	Label string
}

// getLabel returns the field's display label: the explicit Label if set,
// otherwise the Path bash-arg spaced and uppercased.
func (f Field) getLabel() string {
	if f.Label != "" {
		return f.Label
	}

	return deriveLabel(f.Path)
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

	// Footer is free text drawn below the block, after a blank line.
	Footer string
}

// Draw renders the block to text: an optional title line followed by its
// lines or table, aligned to columns. It is the layout hook of the package:
// column alignment today, terminal-width shrinking and lipgloss styling
// later, so swapping the layout engine only touches this method. An empty
// Draw renders the block to text: an optional title line followed by its
// lines, table columns shrunk to fit the given width (0 or negative keeps
// every column). Only table blocks shrink, like v1's formatGrid.
func (b Block) Draw(width, colPadding int) string {
	if colPadding == 0 {
		colPadding = 2 // v1 default
	}

	// A single-cell block is a pre-formatted string (a Formatter, format func
	// or leaf escaping to one cell, possibly containing newlines); echo it
	// verbatim, tabulating it would pad every line to the widest one.
	if len(b.Header) == 0 && len(b.Rows) == 1 && len(b.Rows[0]) == 1 {
		var s strings.Builder
		if b.Title != "" {
			s.WriteString(b.Title)
			s.WriteString(":\n")
		}
		s.WriteString(strings.TrimSpace(b.Rows[0][0]))
		if b.Footer != "" {
			s.WriteString("\n\n")
			s.WriteString(b.Footer)
		}

		return s.String()
	}

	var s strings.Builder
	if b.Title != "" {
		s.WriteString(b.Title)
		s.WriteRune(':')
		s.WriteRune('\n')
	}

	lines := make([][]string, 0, len(b.Rows)+1)
	if len(b.Header) > 0 && !b.NoHeaders {
		lines = append(lines, b.Header)
	}
	lines = append(lines, b.Rows...)

	cols := len(b.Header)
	if cols == 0 && len(b.Rows) > 0 {
		cols = len(b.Rows[0])
	}
	if cols > 0 && width > 0 && len(b.Header) > 0 {
		if fit := MaxCols(lines, width, colPadding); fit < cols {
			cols = fit
		}
	}

	w := tabwriter.NewWriter(&s, 5, 1, colPadding, ' ', tabwriter.ANSIGraphicsRendition)
	for _, line := range lines {
		if cols > 0 && len(line) > cols {
			line = line[:cols]
		}
		fmt.Fprintln(w, strings.Join(line, "\t"))
	}
	w.Flush()

	out := strings.TrimSpace(s.String())
	if b.Footer != "" {
		if out != "" {
			out += "\n\n"
		}
		out += b.Footer
	}

	return out
}

// Empty reports whether block draws nothing at all.
func (b Block) Empty() bool {
	return b.Title == "" && len(b.Header) == 0 && len(b.Rows) == 0 && b.Footer == ""
}

// deriveLabel renders a path as an ALL-CAPS column header.
func deriveLabel(path string) string {
	label := strings.ReplaceAll(path, ".", " ")
	label = strcase.ToBashArg(label)
	label = strings.ReplaceAll(label, "-", " ")

	return strings.ToUpper(label)
}

// isOmittable reports whether v is empty per encoding/json's omitempty:
// an invalid or zero value, or a length-zero slice, array or map.
func isOmittable(v reflect.Value) bool {
	// IsZero and Len panic on invalid values or non-container kinds.
	if !v.IsValid() || v.IsZero() {
		return true
	}

	switch v.Kind() {
	case reflect.Slice, reflect.Array, reflect.Map:
		return v.Len() == 0
	default:
		return false
	}
}

// defaultTableFields lists the default columns for a table of t: each field
// becomes one column, and anonymous (embedded) fields contribute their
// promoted leaf fields instead of a single column, so tables agree with the
// list walk on embedded structs (v1's getDefaultFieldsOpt behavior).
func defaultTableFields(t reflect.Type) []Field {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}

	if t.Kind() != reflect.Struct {
		return nil
	}

	fields := []Field(nil)
	for field := range t.Fields() {
		if field.PkgPath != "" {
			continue
		}
		if field.Anonymous {
			fields = append(fields, defaultTableFields(field.Type)...)

			continue
		}
		fields = append(fields, Field{Path: field.Name})
	}

	return fields
}

// formatCell renders one table cell, mirroring v1's marshalInlineSlice: nil
// slices fold into "-", registered slices escape through Format, scalar
// slices inline as "[a b]", composite slices as their count.
func formatCell(v any) (string, error) {
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return "-", nil
		}
		value = value.Elem()
	}

	// NOTE: cell marshalling follows v1: a registered format func wins for the
	// map's type (e.g. instance's map[string]*instance.Volume -> len), otherwise
	// the value is dumped as sorted "key  value" lines like v1's MarshalMap
	// (e.g. k8s version list's AvailableKubeletArgs).
	if value.Kind() == reflect.Map {
		if value.IsNil() {
			return "-", nil
		}
		if _, ok := getFormatFunc(value.Type()); ok {
			return Format(v)
		}

		return marshalMapCell(value)
	}
	if value.Kind() != reflect.Slice && value.Kind() != reflect.Array {
		return Format(v)
	}
	if value.IsNil() {
		return "-", nil
	}
	if _, ok := getFormatFunc(value.Type()); ok {
		return Format(v)
	}

	itemType := value.Type().Elem()
	for itemType.Kind() == reflect.Pointer {
		itemType = itemType.Elem()
	}
	if itemType.Kind() == reflect.Struct || itemType.Kind() == reflect.Map ||
		itemType.Kind() == reflect.Slice {
		return strconv.Itoa(value.Len()), nil
	}

	return fmt.Sprint(value.Interface()), nil
}

// marshalMapCell renders a map cell like v1's MarshalMap: keys sorted, one
// tabwriter-aligned "key\tvalue" line each, slice values sorted and
// space-joined, nested maps dumped recursively.
func marshalMapCell(value reflect.Value) (string, error) {
	buffer := bytes.Buffer{}

	w := tabwriter.NewWriter(&buffer, 5, 1, 2, ' ', tabwriter.ANSIGraphicsRendition)

	mapKeys := value.MapKeys()
	sort.Slice(mapKeys, func(i, j int) bool {
		return mapKeys[i].String() < mapKeys[j].String()
	})

	for _, mapKey := range mapKeys {
		mapValue := value.MapIndex(mapKey)

		if mapValue.Kind() == reflect.Slice || mapValue.Kind() == reflect.Array {
			values := make([]string, mapValue.Len())
			for i := range mapValue.Len() {
				values[i] = fmt.Sprint(mapValue.Index(i).Interface())
			}
			sort.Strings(values)
			fmt.Fprintf(w, "%s\t%s\n", mapKey.String(), strings.Join(values, " "))
		} else {
			content, err := Format(mapValue.Interface())
			if err != nil {
				return "", err
			}
			fmt.Fprintf(w, "%s\t%s\n", mapKey.String(), content)
		}
	}

	w.Flush()

	return strings.TrimSpace(buffer.String()), nil
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

// MaxCols returns the number of columns that fit within width, mirroring v1:
// each column takes its widest cell plus colPadding, and the first column
// whose running width exceeds width is dropped.
func MaxCols(lines [][]string, width, colPadding int) int {
	cols := len(lines[0])
	colMaxSize := make([]int, len(lines[0]))

	for _, line := range lines {
		lineSize := 0
		for j := 0; j < cols; j++ {
			size := len(line[j]) + colPadding
			if size >= colMaxSize[j] {
				colMaxSize[j] = size
			}
			lineSize += colMaxSize[j]
			if lineSize > width {
				cols = j
			}
		}
	}

	return cols
}

// walk represents value and its nodes to f; f installs its own walk logic:
// it's called for every node with its dotted-path keys (Sub.X, Strings.0,
// Map.key) and returns whether to keep descending after it. Nil pointers and
// interfaces yield no node at all.
func walk(
	value reflect.Value,
	keys string,
	f func(keys string, value reflect.Value) (bool, error),
) error {
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
			field := value.Type().Field(i)
			if !value.Field(i).CanInterface() {
				// Unexported fields cannot be turned into an interface by the
				// visitor; reflect would panic on Interface(). Match
				// encoding/json, which still flattens the promoted exported
				// fields of an anonymous embedded struct or pointer (e.g. a
				// *customServer embed) even when its own type is unexported.
				// NOTE: we only descend into struct-like embeds; an unexported
				// leaf could not be formatted anyway.
				if !field.Anonymous {
					continue
				}
				switch value.Field(i).Kind() {
				case reflect.Struct, reflect.Pointer, reflect.Interface:
				default:
					continue
				}
			}
			parent := keys
			if !field.Anonymous {
				parent = joinKeys(keys, field.Name)
			}
			err := walk(value.Field(i), parent, f)
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

// registeredLayout returns the layout registered for the pointed-to type of
// v, mirroring Render's dereference before the layouts lookup.
func registeredLayout(v any) (Layout, bool) {
	value := reflect.ValueOf(v)
	for value.Kind() == reflect.Pointer || value.Kind() == reflect.Interface {
		if value.IsNil() {
			return Layout{}, false
		}
		value = value.Elem()
	}
	if !value.IsValid() {
		return Layout{}, false
	}
	l, ok := layouts.Load(value.Type())
	if !ok {
		return Layout{}, false
	}

	return l.(Layout), true
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
