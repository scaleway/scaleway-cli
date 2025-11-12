package gofields

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

type NilValueError struct {
	Path string
}

func NewNilValueError(path string) *NilValueError {
	return &NilValueError{Path: path}
}

func (e *NilValueError) Error() string {
	return fmt.Sprintf("field %s is nil", e.Path)
}

// GetValue will extract the value at the given path from the data Go struct
// E.g data = { Friends: []Friend{ { Name: "John" } }, path = "Friends.0.Name"  will return "John"
func GetValue(data any, path string) (any, error) {
	value := reflect.ValueOf(data)
	v, err := getValue(value, []string{"$"}, strings.Split(path, "."))
	if err != nil {
		return nil, err
	}

	return v.Interface(), nil
}

func getValue(value reflect.Value, parents []string, path []string) (reflect.Value, error) {
	if len(path) == 0 {
		return value, nil
	}

	if !value.IsValid() || IsNil(value) {
		return reflect.Value{}, NewNilValueError(strings.Join(parents, "."))
	}

	if value.Type().Kind() == reflect.Ptr {
		return getValue(value.Elem(), parents, path)
	}

	switch value.Kind() {
	case reflect.Slice:
		idx, err := strconv.Atoi(path[0])
		if err != nil {
			return reflect.Value{}, fmt.Errorf(
				"trying to access array %s but %s is not a numerical index",
				strings.Join(parents, "."),
				path[0],
			)
		}
		if idx >= value.Len() {
			return reflect.Value{}, fmt.Errorf(
				"trying to access array %s but %d is out of range",
				strings.Join(parents, "."),
				idx,
			)
		}

		return getValue(value.Index(idx), append(parents, path[0]), path[1:])
	case reflect.Map:
		v := value.MapIndex(reflect.ValueOf(path[0]))
		if !v.IsValid() {
			return reflect.Value{}, fmt.Errorf(
				"trying to access map %s but %s key does not exist",
				strings.Join(parents, "."),
				path[0],
			)
		}

		return getValue(v, append(parents, path[0]), path[1:])
	case reflect.Struct:
		f, exist := value.Type().FieldByName(path[0])
		if !exist {
			return reflect.Value{}, fmt.Errorf(
				"field %s does not exist in %s",
				path[0],
				strings.Join(parents, "."),
			)
		}
		if !isFieldPublic(f) {
			return reflect.Value{}, fmt.Errorf(
				"field %s is private in %s",
				path[0],
				strings.Join(parents, "."),
			)
		}
		v := value.FieldByIndex(f.Index)

		return getValue(v, append(parents, path[0]), path[1:])
	default:
		return reflect.Value{}, fmt.Errorf(
			"cannot get %s in field %s",
			strings.Join(path, "."),
			strings.Join(parents, "."),
		)
	}
}

// GetType will extract the type at the given path from the data Go struct
// E.g data = { Friends: []Friend{ { Name: "John" } }, path = "Friends.0.Name"  will return "John"
func GetType(t reflect.Type, path string) (reflect.Type, error) {
	return getType(t, []string{"$"}, strings.Split(path, "."))
}

func getType(t reflect.Type, parents []string, path []string) (reflect.Type, error) {
	if len(path) == 0 {
		return t, nil
	}

	if t.Kind() == reflect.Ptr {
		return getType(t.Elem(), parents, path)
	}

	switch t.Kind() {
	case reflect.Slice:
		_, err := strconv.Atoi(path[0])
		if err != nil {
			return nil, fmt.Errorf(
				"trying to access array %s but %s is not a numerical index",
				strings.Join(parents, "."),
				path[0],
			)
		}

		return getType(t.Elem(), append(parents, path[0]), path[1:])
	case reflect.Map:
		return getType(t.Elem(), append(parents, path[0]), path[1:])
	case reflect.Struct:
		field, exist := t.FieldByName(path[0])
		if !exist {
			return nil, fmt.Errorf(
				"field %s does not exist in %s",
				path[0],
				strings.Join(parents, "."),
			)
		}
		if !isFieldPublic(field) {
			return nil, fmt.Errorf("field %s is private in %s", path[0], strings.Join(parents, "."))
		}

		return getType(field.Type, append(parents, path[0]), path[1:])
	default:
		return nil, fmt.Errorf(
			"cannot get %s in field %s",
			strings.Join(path, "."),
			strings.Join(parents, "."),
		)
	}
}

type ListFieldFilter func(reflect.Type, string) bool

// ListFields will recursively list all fields path that can be used with GetType or GetValue
func ListFields(t reflect.Type) []string {
	return listFields(t, []string(nil), nil)
}

// ListFieldsWithFilter is the same as ListFields but accept a filter method that will be call for each fields
func ListFieldsWithFilter(t reflect.Type, filter ListFieldFilter) []string {
	return listFields(t, []string(nil), filter)
}

func listFields(t reflect.Type, parents []string, filter ListFieldFilter) []string {
	if t.Kind() == reflect.Ptr {
		return listFields(t.Elem(), parents, filter)
	}

	switch t.Kind() {
	case reflect.Slice:
		return listFields(t.Elem(), append(parents, "<index>"), filter)
	case reflect.Map:
		return listFields(t.Elem(), append(parents, "<key>"), filter)
	case reflect.Struct:
		res := []string(nil)
		for i := range t.NumField() {
			field := t.Field(i)

			if !isFieldPublic(field) {
				continue
			}

			fieldParents := parents
			if !field.Anonymous {
				fieldParents = append(parents, field.Name) //nolint:gocritic
			}

			res = append(res, listFields(field.Type, fieldParents, filter)...)
		}

		return res
	default:
		path := strings.Join(parents, ".")
		if filter != nil && !filter(t, path) {
			return nil
		}

		return []string{path}
	}
}

// IsNil test if a given value is nil. It is saf to call the method with non nillable value like scalar types
func IsNil(value reflect.Value) bool {
	return (value.Kind() == reflect.Ptr || value.Kind() == reflect.Slice || value.Kind() == reflect.Map) &&
		value.IsNil()
}

// isFieldPublic returns true if the given field is public (Name starts with an uppercase)
func isFieldPublic(field reflect.StructField) bool {
	return field.Name[0] >= 'A' && field.Name[0] <= 'Z'
}
