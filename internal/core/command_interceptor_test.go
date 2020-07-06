package core

import (
	"context"
	"testing"

	"github.com/alecthomas/assert"
)

func Test_CombineCommandInterceptor(t *testing.T) {
	runner := func(context.Context, interface{}) (interface{}, error) {
		return []string{"runner"}, nil
	}

	newInterceptor := func(name string) CommandInterceptor {
		return func(ctx context.Context, args interface{}, runner CommandRunner) (interface{}, error) {
			res, _ := runner(ctx, args)
			return append([]string{name}, res.([]string)...), nil
		}
	}

	type TestCase struct {
		Interceptors []CommandInterceptor
		Expected     []string
	}

	run := func(tc *TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			interceptor := combineCommandInterceptor(tc.Interceptors...)
			res, _ := interceptor(nil, nil, runner)
			assert.Equal(t, tc.Expected, res)
		}
	}

	t.Run("simple", run(&TestCase{
		Interceptors: []CommandInterceptor{
			newInterceptor("A"),
		},
		Expected: []string{"A", "runner"},
	}))

	t.Run("with two interceptor", run(&TestCase{
		Interceptors: []CommandInterceptor{
			newInterceptor("A"),
			newInterceptor("B"),
		},
		Expected: []string{"A", "B", "runner"},
	}))

	t.Run("with nil", run(&TestCase{
		Interceptors: []CommandInterceptor{
			newInterceptor("A"),
			nil,
			newInterceptor("B"),
		},
		Expected: []string{"A", "B", "runner"},
	}))

	t.Run("with three interceptor", run(&TestCase{
		Interceptors: []CommandInterceptor{
			newInterceptor("A"),
			newInterceptor("B"),
			newInterceptor("C"),
		},
		Expected: []string{"A", "B", "C", "runner"},
	}))
}
