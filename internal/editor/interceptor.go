package editor

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

func Interceptor(getter *core.Command) core.CommandInterceptor {
	return func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		argsV := reflect.ValueOf(argsI)

		editedArgs, err := UpdateResourceEditor(argsI, getter.ArgsType, func(i interface{}) (interface{}, error) {
			return getter.Run(ctx, i)
		})
		if err != nil {
			return nil, fmt.Errorf("failed to edit args: %w", err)
		}

		// TODO: only map diff
		valueMapper(argsV, reflect.ValueOf(editedArgs))

		return runner(ctx, argsI)
	}
}
