package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"sort"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// loadArgsFileContent loads the content of the files referenced by arguments that
// support it (ArgSpec.CanLoadFile), i.e. arguments whose value starts with "@".
func loadArgsFileContent(cmd *Command, cmdArgs any) error {
	for _, argSpec := range cmd.ArgSpecs {
		if !argSpec.CanLoadFile {
			continue
		}

		fieldName := strcase.ToPublicGoName(argSpec.Name)
		if _, err := loadFileContentForField(
			reflect.ValueOf(cmdArgs),
			strings.Split(fieldName, "."),
		); err != nil {
			return err
		}
	}

	return nil
}

// loadFileContentForField traverses value following the path described by parts and, for each
// leaf value that starts with "@", replaces it with the content of the referenced file.
// It returns whether the path resolved to at least one leaf. Unresolvable paths are skipped.
func loadFileContentForField(value reflect.Value, parts []string) (bool, error) {
	if len(parts) == 0 {
		err := setLeafFromMaybeFile(value)

		return true, err
	}

	switch value.Kind() {
	case reflect.Pointer:
		if value.IsNil() {
			return false, nil
		}

		return loadFileContentForField(value.Elem(), parts)

	case reflect.Struct:
		fieldName := strcase.ToPublicGoName(parts[0])
		if field, ok := directStructField(value, fieldName); ok {
			return loadFileContentForField(field, parts[1:])
		}

		// If it does not exist we try to find it in a nested anonymous field
		for i := range value.NumField() {
			if !value.Type().Field(i).Anonymous {
				continue
			}

			handled, err := loadFileContentForField(value.Field(i), parts)
			if handled || err != nil {
				return handled, err
			}
		}

		return false, nil

	case reflect.Slice:
		handled := false
		for i := range value.Len() {
			fieldHandled, err := loadFileContentForField(value.Index(i), parts[1:])
			if err != nil {
				return false, err
			}

			handled = handled || fieldHandled
		}

		return handled, nil

	case reflect.Map:
		if value.IsNil() {
			return false, nil
		}

		mapKeys := value.MapKeys()
		sort.Slice(mapKeys, func(i, j int) bool {
			return mapKeys[i].String() < mapKeys[j].String()
		})

		handled := false
		for _, mapKey := range mapKeys {
			mapValue := value.MapIndex(mapKey)

			// Map values are not addressable, so operate on an addressable copy and write it back.
			mapValueCopy := reflect.New(mapValue.Type()).Elem()
			mapValueCopy.Set(mapValue)

			fieldHandled, err := loadFileContentForField(mapValueCopy, parts[1:])
			if err != nil {
				return false, err
			}

			value.SetMapIndex(mapKey, mapValueCopy)

			handled = handled || fieldHandled
		}

		return handled, nil
	}

	return false, nil
}

// directStructField returns the non-anonymous field of the struct value named fieldName, if any.
func directStructField(value reflect.Value, fieldName string) (reflect.Value, bool) {
	for i := range value.NumField() {
		field := value.Type().Field(i)
		if !field.Anonymous && field.Name == fieldName {
			return value.Field(i), true
		}
	}

	return reflect.Value{}, false
}

// setLeafFromMaybeFile replaces value with the content of the file it references when it starts with "@".
func setLeafFromMaybeFile(value reflect.Value) error {
	newValue, err := fileContentForLeaf(value)
	if err != nil {
		return err
	}

	if newValue.IsValid() {
		value.Set(newValue)
	}

	return nil
}

// fileContentForLeaf computes the value a leaf should hold after file loading.
// The returned value is invalid when no change is required.
func fileContentForLeaf(leaf reflect.Value) (reflect.Value, error) {
	switch i := leaf.Interface().(type) {
	case io.Reader:
		b, err := io.ReadAll(i)
		if err != nil {
			return reflect.Value{}, fmt.Errorf("could not read argument: %s", err)
		}

		if strings.HasPrefix(string(b), "@") {
			content, err := os.ReadFile(string(b)[1:])
			if err != nil {
				return reflect.Value{}, fmt.Errorf("could not open requested file: %s", err)
			}

			return reflect.ValueOf(bytes.NewBuffer(content)), nil
		}

		// Reader must be re-created as it can only be read once.
		return reflect.ValueOf(bytes.NewReader(b)), nil
	case *string:
		if strings.HasPrefix(*i, "@") {
			content, err := os.ReadFile((*i)[1:])
			if err != nil {
				return reflect.Value{}, fmt.Errorf("could not open requested file: %s", err)
			}

			s := string(content)

			return reflect.ValueOf(&s), nil
		}
	case string:
		if strings.HasPrefix(i, "@") {
			content, err := os.ReadFile(i[1:])
			if err != nil {
				return reflect.Value{}, fmt.Errorf("could not open requested file: %s", err)
			}

			return reflect.ValueOf(string(content)), nil
		}
	case []byte:
		if strings.HasPrefix(string(i), "@") {
			content, err := os.ReadFile(string(i[1:]))
			if err != nil {
				return reflect.Value{}, fmt.Errorf("could not open requested file: %s", err)
			}

			return reflect.ValueOf(content), nil
		}
	case nil:
	default:
		panic(fmt.Errorf("unsupported field type: %T", i))
	}

	return reflect.Value{}, nil
}
