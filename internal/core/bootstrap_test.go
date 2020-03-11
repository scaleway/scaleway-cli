package core

import (
	"context"
	"reflect"
	"testing"

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
}
