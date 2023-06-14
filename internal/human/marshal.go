package human

import (
	"bytes"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/gofields"
	"github.com/scaleway/scaleway-cli/v2/internal/tabwriter"
	"github.com/scaleway/scaleway-cli/v2/internal/terminal"
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
	if !rValue.IsValid() {
		return defaultMarshalerFunc(nil, opt)
	}

	rType := rValue.Type()

	// safely get the marshalerFunc
	marshalerFunc, _ := getMarshalerFunc(rType)
	isNil := isInterfaceNil(data)
	isSlice := rType.Kind() == reflect.Slice

	switch {
	// If data is nil and is not a slice ( In case of a slice we want to print header row and not a simple dash )
	case isNil && !isSlice:
		return defaultMarshalerFunc(nil, opt)

	// If data has a registered MarshalerFunc call it
	case marshalerFunc != nil:
		return marshalerFunc(rValue.Interface(), opt)

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
	case isSlice:
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
	// subOpts that will be passed down to sub marshaler (field, sub struct, slice item ...)
	subOpts := &MarshalOpt{}

	// Marshal sections
	sectionsStrs := []string(nil)
	sectionFieldNames := map[string]bool{}
	for _, section := range opt.Sections {
		sectionStr, err := marshalSection(section, value, subOpts)
		if err != nil {
			return "", err
		}

		sectionFieldNames[section.FieldName] = true

		if sectionStr != "" {
			sectionsStrs = append(sectionsStrs, sectionStr)
		}
	}

	var marshal func(reflect.Value, []string) ([][]string, error)

	marshal = func(value reflect.Value, keys []string) ([][]string, error) {
		if _, isSection := sectionFieldNames[strings.Join(keys, ".")]; isSection {
			return nil, nil
		}
		rType := value.Type()

		// safely get the marshaler func
		marshalerFunc, _ := getMarshalerFunc(rType)

		switch {
		// If data is nil
		case isInterfaceNil(value.Interface()):
			return nil, nil

		// If data has a registered MarshalerFunc call it.
		case marshalerFunc != nil:
			str, err := marshalerFunc(value.Interface(), subOpts)
			return [][]string{{strings.Join(keys, "."), str}}, err

		// If data is a stringers
		case rType.Implements(reflect.TypeOf((*fmt.Stringer)(nil)).Elem()):
			return [][]string{{strings.Join(keys, "."), value.Interface().(fmt.Stringer).String()}}, nil

		case rType.Kind() == reflect.Ptr:
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
			for _, fieldIndex := range getStructFieldsIndex(value.Type()) {
				subData, err := marshal(value.FieldByIndex(fieldIndex), append(keys, value.Type().FieldByIndex(fieldIndex).Name))
				if err != nil {
					return nil, err
				}
				data = append(data, subData...)
			}

			return data, nil
		case rType.Kind() == reflect.Interface:
			// If type is interface{}
			// marshal the underlying type
			return marshal(value.Elem(), keys)
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

// getStructFieldsIndex will return a list of fieldIndex ([]int) sorted by their position in the Go struct.
// This function will handle anonymous field and make sure that if a field is overwritten only the highest is returned.
// You can use reflect GetFieldByIndex([]int) to get the correct field.
func getStructFieldsIndex(v reflect.Type) [][]int {
	// Using a map we make sure only the field with the highest order is returned for a given Name
	found := map[string][]int{}

	var recFunc func(v reflect.Type, parent []int)
	recFunc = func(v reflect.Type, parent []int) {
		for v.Kind() == reflect.Ptr {
			v = v.Elem()
		}

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// If a field is anonymous we start recursive call
			if field.Anonymous {
				recFunc(v.Field(i).Type, append(parent, i))
			} else {
				// else we add the field in the found map
				found[field.Name] = append(parent, i)
			}
		}
	}
	recFunc(v, []int(nil))

	result := [][]int(nil)
	for _, value := range found {
		result = append(result, value)
	}

	sort.Slice(result, func(i, j int) bool {
		n := 0
		for n < len(result[i]) && n < len(result[j]) {
			if result[i][n] != result[j][n] {
				return result[i][n] < result[j][n]
			}
			n++
		}
		// if equal, less should be false
		return false
	})

	return result
}

func marshalSlice(slice reflect.Value, opt *MarshalOpt) (string, error) {
	// Resole itemType and get rid of all pointer level if needed.
	itemType := slice.Type().Elem()
	for itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	// If itemType is not a struct (e.g []string) we just stringify it
	if itemType.Kind() != reflect.Struct {
		return fmt.Sprint(slice.Interface()), nil
	}

	// If there is no Field in opt we generated default one using reflect
	if len(opt.Fields) == 0 {
		opt.Fields = getDefaultFieldsOpt(itemType)
	}

	subOpts := &MarshalOpt{TableCell: true}

	// Validate that all field exist
	for _, f := range opt.Fields {
		_, err := gofields.GetType(itemType, f.FieldName)
		if err != nil {
			return "", &UnknownFieldError{
				FieldName:   f.FieldName,
				ValidFields: gofields.ListFields(itemType),
			}
		}
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
			v, err := gofields.GetValue(item.Interface(), fieldSpec.FieldName)
			if err != nil {
				logger.Debugf("invalid getFieldValue(): '%v' might not be exported", fieldSpec.FieldName)
				row = append(row, "")
				continue
			}
			fieldValue := reflect.ValueOf(v)

			str := ""
			switch {
			// Handle inline slice.
			case fieldValue.Type().Kind() == reflect.Slice:
				str, err = marshalInlineSlice(fieldValue)
			default:
				str, err = Marshal(fieldValue.Interface(), subOpts)
			}
			if err != nil {
				return "", err
			}
			row = append(row, str)
		}
		grid = append(grid, row)
	}
	return formatGrid(grid, !opt.DisableShrinking)
}

// marshalInlineSlice transforms nested scalar slices in an inline string representation
// and other types of slices in a count representation.
func marshalInlineSlice(slice reflect.Value) (string, error) {
	if slice.IsNil() {
		return defaultMarshalerFunc(nil, nil)
	}

	itemType := slice.Type().Elem()
	for itemType.Kind() == reflect.Ptr {
		itemType = itemType.Elem()
	}

	// safely get the marshalerFunc
	marshalerFunc, _ := getMarshalerFunc(slice.Type())

	switch {
	// If marshaler func is available.
	// As we cannot set MarshalOpt of a nested slice opt will always be nil here.
	case marshalerFunc != nil:
		return marshalerFunc(slice.Interface(), nil)

	// If it is a slice of scalar values.
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
		title = cases.Title(language.English).String(strings.ReplaceAll(title, ".", " - "))
	}
	subOpt.Title = title

	field, err := gofields.GetValue(value.Interface(), section.FieldName)
	if err != nil {
		if section.HideIfEmpty {
			if _, ok := err.(*gofields.NilValueError); ok {
				return "", nil
			}
		}

		return "", err
	}

	if section.HideIfEmpty && reflect.ValueOf(field).IsZero() {
		return "", nil
	}

	return Marshal(field, &subOpt)
}

func formatGrid(grid [][]string, shrinkColumns bool) (string, error) {
	buffer := bytes.Buffer{}
	maxCols := computeMaxCols(grid)
	w := tabwriter.NewWriter(&buffer, 5, 1, colPadding, ' ', tabwriter.ANSIGraphicsRendition)
	for _, line := range grid {
		if shrinkColumns {
			line = line[:maxCols]
		}
		fmt.Fprintln(w, strings.Join(line, "\t"))
	}
	w.Flush()
	return strings.TrimSpace(buffer.String()), nil
}

// computeMaxCols calculates how many row we can fit in terminal width.
func computeMaxCols(grid [][]string) int {
	maxCols := len(grid[0])
	width := terminal.GetWidth()
	// If we are not writing to Stdout or through a tty Stdout, returns max length
	if !terminal.IsTerm() || width == 0 {
		return maxCols
	}
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

		if field.Anonymous {
			results = append(results, getDefaultFieldsOpt(fieldType)...)
			continue
		}

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
