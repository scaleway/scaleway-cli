package core

import (
	"fmt"
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
			if v.Kind() != reflect.String {
				continue
			}

			valueString := v.String()
			if strings.HasPrefix(valueString, "@") {
				content, err := ioutil.ReadFile(valueString[1:])
				if err != nil {
					return fmt.Errorf("could not open requested file: %s", err)
				}
				v.SetString(string(content))
			}
		}
	}

	return nil
}
