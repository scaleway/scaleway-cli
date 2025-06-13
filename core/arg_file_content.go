package core

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// loadArgsFileContent will hydrate args with default values.
func loadArgsFileContent(cmd *Command, cmdArgs any) error {
	for _, argSpec := range cmd.ArgSpecs {
		if !argSpec.CanLoadFile {
			continue
		}

		fieldName := strcase.ToPublicGoName(argSpec.Name)
		fieldValues, err := GetValuesForFieldByName(
			reflect.ValueOf(cmdArgs),
			strings.Split(fieldName, "."),
		)
		if err != nil {
			continue
		}

		for _, v := range fieldValues {
			switch i := v.Interface().(type) {
			case io.Reader:
				b, err := io.ReadAll(i)
				if err != nil {
					return fmt.Errorf("could not read argument: %s", err)
				}

				if strings.HasPrefix(string(b), "@") {
					content, err := os.ReadFile(string(b)[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					test := bytes.NewBuffer(content)
					v.Set(reflect.ValueOf(test))
				} else {
					// Reader must be re-created as it can only be read once.
					r := bytes.NewReader(b)
					v.Set(reflect.ValueOf(r))
				}
			case *string:
				if strings.HasPrefix(*i, "@") {
					content, err := os.ReadFile((*i)[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					v.SetString(string(content))
				}
			case string:
				if strings.HasPrefix(i, "@") {
					content, err := os.ReadFile(i[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					v.SetString(string(content))
				}
			case []byte:
				if strings.HasPrefix(string(i), "@") {
					content, err := os.ReadFile(string(i)[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					v.SetBytes(content)
				}
			case nil:
				continue
			default:
				panic(fmt.Errorf("unsupported field type: %T", v.Interface()))
			}
		}
	}

	return nil
}
