package core_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/args"
	"github.com/scaleway/scaleway-cli/v2/internal/interactive"
	"github.com/stretchr/testify/assert"
)

func TestInterruptError(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("unknown-command", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "interrupt",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return nil, &interactive.InterruptError{}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test interrupt error",
		Check:           core.TestCheckExitCode(130),
	}))
	t.Run("exit-code", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "code",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return nil, &core.CliError{Code: 99}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test code error",
		Check:           core.TestCheckExitCode(99),
	}))
	t.Run("empty-error", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return nil, &core.CliError{Code: 99, Empty: true}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test empty error",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(99),
			core.TestCheckGolden(),
		),
	}))
	t.Run("empty-error-json", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "error",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return nil, &core.CliError{Code: 99, Empty: true}
				},
			},
		),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw -o json test empty error",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(99),
			core.TestCheckGolden(),
		),
	}))
	t.Run("empty-success", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return &core.SuccessResult{
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
		Check:           core.TestCheckGolden(),
	}))
	t.Run("empty-success-json", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return &core.SuccessResult{
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
		Check:           core.TestCheckGolden(),
	}))
	t.Run("empty-list-json", core.Test(&core.TestConfig{
		Commands: core.NewCommands(
			&core.Command{
				Namespace: "test",
				Resource:  "empty",
				Verb:      "success",
				ArgsType:  reflect.TypeOf(args.RawArgs{}),
				Run: func(_ context.Context, _ any) (i any, e error) {
					return []int(nil), nil
				},
				AllowAnonymousClient: true,
			},
		),
		Cmd: "scw -o json test empty success",
		Check: func(t *testing.T, ctx *core.CheckFuncCtx) {
			t.Helper()
			assert.Equal(t, "[]\n", string(ctx.Stdout))
		},
	}))
}
