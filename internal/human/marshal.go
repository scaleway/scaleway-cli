package human

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/internal/terminal"
	"github.com/scaleway/scaleway-sdk-go/logger"
	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// Padding between column
const colPadding = 2

// Marshaler allow custom display for some type when printed using HumanPrinter
type Marshaler interface {
	MarshalHuman() (string, error)
}

func Marshal(data interface{}, opt *MarshalOpt) (string, error) {

	// Debug infos
	logger.Debugf("marshalling type '%v'", reflect.TypeOf(data))

	// If data is nil there is nothing to print
	if data == nil {
		return "", nil
	}

	if opt == nil {
		opt = &MarshalOpt{}
	}

	if opt.Title != "" {
		subOpt := *opt
		subOpt.Title = ""
		body, err := Marshal(data, &subOpt)
		return terminal.Style(opt.Title+":", color.Bold) + "\n" + body, err
	}

	rValue := reflect.ValueOf(data)
	rType := rValue.Type()

	switch {

	// If data has a registered MarshalerFunc call it
	case marshalerFuncs[rType] != nil:
		return marshalerFuncs[rType](rValue.Interface(), opt)

	// Handle special well known interface
	case rType.Implements(reflect.TypeOf((*Marshaler)(nil)).Elem()):
		return rValue.Interface().(Marshaler).MarshalHuman()

	// Handle errors
	case rType.Implements(reflect.TypeOf((*error)(nil)).Elem()):
		return terminal.Style(Capitalize(rValue.Interface().(error).Error()), color.FgRed), nil

	// Handle stringers
	case rType.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()):
		return rValue.Interface().(fmt.Stringer).String(), nil

	// If data is a pointer dereference an call Marshal again
	case rType.Kind() == reflect.Ptr:
		return Marshal(rValue.Elem().Interface(), opt)

	// If data is a slice uses marshalSlice
	case rType.Kind() == reflect.Slice:
		return marshalSlice(rValue, opt)

	// If data is a struct uses marshalStruct
	case rType.Kind() == reflect.Struct:
		return marshalStruct(rValue, opt)

	// by default we use defaultMarshalerFunc
	default:
		return defaultMarshalerFunc(rValue.Interface(), opt)
	}
}

func marshalStruct(value reflect.Value, opt *MarshalOpt) (string, error) {

	// subOpts that will be passed down to sub marshaller (field, sub struct, slice item ...)
	subOpts := &MarshalOpt{}

	// Marshal sections
	sectionsStrs := []string(nil)
	sectionFieldNames := map[string]bool{}
	for _, section := range opt.Sections {
		sectionStr, err := marshalSection(section, value, subOpts)
		if err != nil {
			return "", err
		}
		sectionsStrs = append(sectionsStrs, sectionStr)
		sectionFieldNames[section.FieldName] = true
	}

	var marshal func(reflect.Value, []string) ([][]string, error)

	marshal = func(value reflect.Value, keys []string) ([][]string, error) {

		if _, isSection := sectionFieldNames[strings.Join(keys, ".")]; isSection {
			return nil, nil
		}
		rType := value.Type()

		switch {

		// If data as a register MarshalerFunc call it
		case marshalerFuncs[rType] != nil:
			str, err := marshalerFuncs[rType](value.Interface(), subOpts)
			return [][]string{{strings.Join(keys, "."), str}}, err

		case rType.Kind() == reflect.Ptr:
			// If src is nil we do not marshal it
			if value.IsNil() {
				return nil, nil
			}
			// If type is a pointer we Marshal pointer.Elem()
			return marshal(value.Elem(), keys)
		case rType.Kind() == reflect.Slice:
			// If type is a slice:
			// We loop through all items and marshal them with key = key.0, key.1, ....
			data := [][]string(nil)
			for i := 0; i < value.Len(); i++ {
				subData, err := marshal(value.Index(i), append(keys, strconv.Itoa(i)))
				if err != nil {
					return nil, err
				}
				data = append(data, subData...)
			}
			return data, nil

		case rType.Kind() == reflect.Map:
			// If map is nil we do not marshal it
			if value.IsNil() {
				return nil, nil
			}

			// If type is a map:
			// We loop through all items and marshal them with key = key.0, key.1, ....
			data := [][]string(nil)

			// Get all map keys and sort them. We assume keys are string
			mapKeys := value.MapKeys()
			sort.Slice(mapKeys, func(i, j int) bool {
				return mapKeys[i].String() < mapKeys[j].String()
			})

			for _, mapKey := range mapKeys {
				mapValue := value.MapIndex(mapKey)
				subData, err := marshal(mapValue, append(keys, mapKey.String()))
				if err != nil {
					return nil, err
				}
				data = append(data, subData...)
			}
			return data, nil

		case rType.Kind() == reflect.Struct:
			// If type is a struct
			// We loop through all struct field
			data := [][]string(nil)
			for i := 0; i < value.NumField(); i++ {
				subData, err := marshal(value.Field(i), append(keys, strcase.ToBashArg(value.Type().Field(i).Name)))
				if err != nil {
					return nil, err
				}
				data = append(data, subData...)
			}

			return data, nil
		default:
			str, err := defaultMarshalerFunc(value.Interface(), subOpts)
			if err != nil {
				return nil, err
			}
			return [][]string{{strings.Join(keys, "."), str}}, nil
		}
	}

	data, err := marshal(value, nil)
	if err != nil {
		return "", err
	}

	buffer := bytes.Buffer{}
	w := tabwriter.NewWriter(&buffer, 5, 1, colPadding, ' ', tabwriter.ANSIGraphicsRendition)
	for _, line := range data {
		fmt.Fprintln(w, strings.Join(line, "\t"))
	}

	if len(sectionsStrs) > 0 {
		fmt.Fprintln(w, "\n"+strings.Join(sectionsStrs, "\n\n"))
	}

	w.Flush()

	return strings.TrimSpace(buffer.String()), nil
}

func marshalSlice(slice reflect.Value, opt *MarshalOpt) (string, error) {

	// Resole itemType and get rid of all pointer level if needed.
	itemType := slice.Type().Elem()
	for itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	// if itemType is not a struct (e.g []string) we just stringify it
	if itemType.Kind() != reflect.Struct {
		return fmt.Sprint(slice.Interface()), nil
	}

	// If there is no Field in opt we generated default one using reflect
	if len(opt.Fields) == 0 {
		opt.Fields = getDefaultFieldsOpt(itemType)
	}

	// We create a in memory grid of the content we want to print
	grid := make([][]string, 0, slice.Len()+1)

	// Generate header row
	headerRow := []string(nil)
	for _, fieldSpec := range opt.Fields {
		headerRow = append(headerRow, fieldSpec.getLabel())
	}
	grid = append(grid, headerRow)

	// For each item in the slice
	for i := 0; i < slice.Len(); i++ {
		item := slice.Index(i)
		row := []string(nil)
		for _, fieldSpec := range opt.Fields {
			fieldValue := getFieldValue(item, fieldSpec.FieldName)
			str := ""
			if fieldValue.IsValid() {
				err := error(nil)
				switch {
				// Handle inline slice.
				case fieldValue.Type().Kind() == reflect.Slice:
					str, err = marshalInlineSlice(fieldValue)
					if err != nil {
						return "", err
					}

				default:
					str, err = Marshal(fieldValue.Interface(), opt)
					if err != nil {
						return "", err
					}
				}
			} else {
				logger.Debugf("invalid getFieldValue(): '%v' might not be exported", fieldSpec.FieldName)
			}
			row = append(row, str)
		}
		grid = append(grid, row)
	}
	return formatGrid(grid)
}

// marshalInlineSlice transform nested scalar slices in an inline string representation
// and other types of slices in a count representation.
func marshalInlineSlice(slice reflect.Value) (string, error) {
	itemType := slice.Type().Elem()
	for itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	switch {
	// If marshaler func is available.
	// As we cannot set MarshalOpt of a nested slice opt will always be nil here.
	case marshalerFuncs[itemType] != nil:
		return marshalerFuncs[itemType](slice.Interface(), nil)

	// If it's a scalar value we marshall it inline.
	case itemType.Kind() != reflect.Slice &&
		itemType.Kind() != reflect.Map &&
		itemType.Kind() != reflect.Struct:
		return fmt.Sprint(slice), nil

	// Use slice count by default.
	default:
		return strconv.Itoa(slice.Len()), nil
	}
}

// marshalSection transforms a field from a struct into a section.
func marshalSection(section *MarshalSection, value reflect.Value, opt *MarshalOpt) (string, error) {
	subOpt := *opt

	title := section.Title
	if title == "" {
		title = strings.ReplaceAll(strcase.ToBashArg(section.FieldName), "-", " ")
		title = strings.Title(strings.ReplaceAll(title, ".", " - "))
	}
	subOpt.Title = title

	field := getFieldValue(value, section.FieldName)
	return Marshal(field.Interface(), &subOpt)
}

func formatGrid(grid [][]string) (string, error) {
	buffer := bytes.Buffer{}
	maxCols := computeMaxCols(grid)
	w := tabwriter.NewWriter(&buffer, 5, 1, colPadding, ' ', tabwriter.ANSIGraphicsRendition)
	for _, line := range grid {
		fmt.Fprintln(w, strings.Join(line[:maxCols], "\t"))
	}
	w.Flush()
	return strings.TrimSpace(buffer.String()), nil
}

// computeMaxCols calculates how many row we can fit in terminal width.
func computeMaxCols(grid [][]string) int {
	maxCols := len(grid[0])
	// If we are not writing to Stdout or through a tty Stdout, returns max length
	if color.NoColor {
		return maxCols
	}
	width := terminal.GetWidth()
	colMaxSize := make([]int, len(grid[0]))
	for i := 0; i < len(grid); i++ {
		lineSize := 0
		for j := 0; j < maxCols; j++ {
			size := len(grid[i][j]) + colPadding
			if size >= colMaxSize[j] {
				colMaxSize[j] = size
			}
			lineSize += colMaxSize[j]
			if lineSize > width {
				maxCols = j
			}
		}
	}
	return maxCols
}

// Generate default []*MarshalFieldOpt using reflect
// It will detect item type of a slice an keep all root level field that are marshalable
func getDefaultFieldsOpt(t reflect.Type) []*MarshalFieldOpt {

	results := []*MarshalFieldOpt(nil)
	// Loop through all struct field
	for fieldIdx := 0; fieldIdx < t.NumField(); fieldIdx++ {
		field := t.Field(fieldIdx)
		fieldType := field.Type

		if isMarshalable(fieldType) {
			spec := &MarshalFieldOpt{
				FieldName: field.Name,
			}
			results = append(results, spec)
		} else {
			logger.Debugf("fieldType '%v' is not marshallable", fieldType)
		}
	}

	return results
}

// getFieldValue will extract a nested field from a Name ( e.g User.Address.Line1 )
func getFieldValue(value reflect.Value, fieldName string) reflect.Value {
	parts := strings.Split(fieldName, ".")

	// Resolve all pointer level
	for value.Kind() == reflect.Ptr {
		value = value.Elem()
	}

	for _, part := range parts {
		value = value.FieldByName(strcase.ToPublicGoName(part))
		if !value.IsValid() {
			return value
		}
		// If value is Nil return invalid value
		if value.Kind() == reflect.Ptr && value.IsNil() {
			return reflect.Value{}
		}

		// Resolve all pointer level
		for value.Kind() == reflect.Ptr {
			value = value.Elem()
		}
	}

	return value
}
