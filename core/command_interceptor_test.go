package core_test

import (
	"context"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/stretchr/testify/assert"
)

func Test_CombineCommandInterceptor(t *testing.T) {
	runner := func(context.Context, any) (any, error) {
		return []string{"runner"}, nil
	}

	newInterceptor := func(name string) core.CommandInterceptor {
		return func(ctx context.Context, args any, runner core.CommandRunner) (any, error) {
			res, _ := runner(ctx, args)

			return append([]string{name}, res.([]string)...), nil
		}
	}

	type TestCase struct {
		Interceptors []core.CommandInterceptor
		Expected     []string
	}

	run := func(tc *TestCase) func(t *testing.T) {
		return func(t *testing.T) {
			t.Helper()
			interceptor := core.CombineCommandInterceptor(tc.Interceptors...)
			res, _ := interceptor(nil, nil, runner)
			assert.Equal(t, tc.Expected, res)
		}
	}

	t.Run("simple", run(&TestCase{
		Interceptors: []core.CommandInterceptor{
			newInterceptor("A"),
		},
		Expected: []string{"A", "runner"},
	}))

	t.Run("with two interceptor", run(&TestCase{
		Interceptors: []core.CommandInterceptor{
			newInterceptor("A"),
			newInterceptor("B"),
		},
		Expected: []string{"A", "B", "runner"},
	}))

	t.Run("with nil", run(&TestCase{
		Interceptors: []core.CommandInterceptor{
			newInterceptor("A"),
			nil,
			newInterceptor("B"),
		},
		Expected: []string{"A", "B", "runner"},
	}))

	t.Run("with three interceptor", run(&TestCase{
		Interceptors: []core.CommandInterceptor{
			newInterceptor("A"),
			newInterceptor("B"),
			newInterceptor("C"),
		},
		Expected: []string{"A", "B", "C", "runner"},
	}))
}
