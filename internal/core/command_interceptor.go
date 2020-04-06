package core

import (
	"context"
	"fmt"

	"github.com/scaleway/scaleway-sdk-go/scw"
)

func CombineInterceptor(interceptors ...CommandInterceptor) CommandInterceptor {
	var result CommandInterceptor
	for _, _interceptor := range interceptors {
		// Assure interceptor do not escape the context of this loop iteration.
		interceptor := _interceptor

		if interceptor == nil {
			continue
		}
		if result == nil {
			result = interceptor
		}

		previousInterceptor := result
		result = func(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
			return previousInterceptor(ctx, args, func(ctx context.Context, arg interface{}) (interface{}, error) {
				return interceptor(ctx, args, runner)
			})
		}
	}
	return result
}

func sdkStdErrorInterceptor(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
	res, err := runner(ctx, args)
	switch sdkError := err.(type) {
	case *scw.ResourceNotFoundError:
		return nil, &CliError{
			Err: fmt.Errorf("cannot find resource '%v' with ID '%v'", sdkError.Resource, sdkError.ResourceID),
		}
	}
	return res, err
}
