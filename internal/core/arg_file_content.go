package core

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

// loadArgsFileContent will hydrate args with default values.
func loadArgsFileContent(cmd *Command, cmdArgs interface{}) error {
	for _, argSpec := range cmd.ArgSpecs {
		if !argSpec.CanLoadFile {
			continue
		}

		fieldName := strcase.ToPublicGoName(argSpec.Name)
		fieldValues, err := getValuesForFieldByName(reflect.ValueOf(cmdArgs), strings.Split(fieldName, "."))
		if err != nil {
			continue
		}

		for _, v := range fieldValues {
			switch i := v.Interface().(type) {
			case io.Reader:
				b, err := ioutil.ReadAll(i)
				if err != nil {
					return fmt.Errorf("could not read argument: %s", err)
				}

				if strings.HasPrefix(string(b), "@") {
					content, err := ioutil.ReadFile(string(b)[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					test := bytes.NewBuffer(content)
					v.Set(reflect.ValueOf(test))
				}
			case *string:
				if strings.HasPrefix(*i, "@") {
					content, err := ioutil.ReadFile((*i)[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					v.SetString(string(content))
				}
			case string:
				if strings.HasPrefix(i, "@") {
					content, err := ioutil.ReadFile(i[1:])
					if err != nil {
						return fmt.Errorf("could not open requested file: %s", err)
					}
					v.SetString(string(content))
				}
			default:
				panic(fmt.Errorf("unsupported field type: %T", v.Interface()))
			}
		}
	}

	return nil
}
