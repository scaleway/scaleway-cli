package core

import (
	"context"
)

func CombineCommandInterceptor(interceptors ...CommandInterceptor) CommandInterceptor {
	var combinedInterceptors CommandInterceptor
	for _, interceptor := range interceptors {
		if interceptor == nil {
			continue
		}
		if combinedInterceptors == nil {
			combinedInterceptors = interceptor
			continue
		}

		previousInterceptor := combinedInterceptors
		combinedInterceptors = func(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
			return previousInterceptor(ctx, args, func(ctx context.Context, arg interface{}) (interface{}, error) {
				return interceptor(ctx, args, runner)
			})
		}
	}
	return combinedInterceptors
}

// sdkStdErrorInterceptor is a command interceptor that will catch sdk standard error and return more friendly CLI error.
func sdkStdErrorInterceptor(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
	// TODO implement error conversion.
	return runner(ctx, args)
}
