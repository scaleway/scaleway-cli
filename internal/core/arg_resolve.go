package core

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/scaleway/scaleway-sdk-go/strcase"
)

type ResolveArgFunc func(ctx context.Context, arg string) (string, error)

func resolveArgs(ctx context.Context, cmd *Command, cmdArgs interface{}) error {
	for _, argSpec := range cmd.ArgSpecs {
		if argSpec.ResolveFunc == nil {
			continue
		}
		fieldName := strcase.ToPublicGoName(argSpec.Name)
		fieldValues, err := getValuesForFieldByName(reflect.ValueOf(cmdArgs), strings.Split(fieldName, "."))
		if err != nil || len(fieldValues) > 1 {
			continue
		}

		newArg, err := argSpec.ResolveFunc(ctx, fieldValues[0].String())
		if err != nil {
			return fmt.Errorf("failed to resolve argument: %w", err)
		}
		fieldValues[0].Set(reflect.ValueOf(newArg))
	}

	return nil
}
