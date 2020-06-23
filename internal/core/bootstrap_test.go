package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/interactive"
)

func TestInterruptError(t *testing.T) {
	t.Run("unknown-command", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "interrupt",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return nil, &interactive.InterruptError{}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test interrupt error",
		Check:           TestCheckExitCode(130),
	}))
	t.Run("exit-code", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "code",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return nil, &CliError{Code: 99}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test code error",
		Check:           TestCheckExitCode(99),
	}))
	t.Run("empty-error", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return nil, &CliError{Code: 99, Empty: true}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test empty error",
		Check: TestCheckCombine(
			TestCheckExitCode(99),
			TestCheckGolden(),
		),
	}))
	t.Run("empty-error-json", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return nil, &CliError{Code: 99, Empty: true}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw -o json test empty error",
		Check: TestCheckCombine(
			TestCheckExitCode(99),
			TestCheckGolden(),
		),
	}))
	t.Run("empty-success", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return &SuccessResult{
						Empty:    true,
						Message:  "dummy",
						Details:  "dummy",
						Resource: "dummy",
						Verb:     "dummy",
					}, nil
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test empty success",
		Check:           TestCheckGolden(),
	}))
	t.Run("empty-success-json", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return &SuccessResult{
						Empty:    true,
						Message:  "dummy",
						Details:  "dummy",
						Resource: "dummy",
						Verb:     "dummy",
					}, nil
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw -o json test empty success",
		Check:           TestCheckGolden(),
	}))
	t.Run("empty-list-json", Test(&TestConfig{
		Commands: NewCommands(
			&Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
					return []int(nil), nil
				},
				AllowAnonymousClient: true,
			},
		),
		Cmd: "scw -o json test empty success",
		Check: func(t *testing.T, ctx *CheckFuncCtx) {
			assert.Equal(t, "[]\n", string(ctx.Stdout))
		},
	}))
}
