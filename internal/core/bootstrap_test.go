package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-cli/internal/interactive"
)

func getCommands() *Commands {
	return NewCommands(
		&Command{
			Namespace: "test",
			Resource:  "interrupt",
			Verb:      "error",
			ArgsType:  reflect.TypeOf(args.RawArgs{}),
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return nil, &interactive.InterruptError{}
			},
		},
	)
}

func TestInterruptError(t *testing.T) {
	t.Run("unknown-command", Test(&TestConfig{
		Commands:     getCommands(),
		UseE2EClient: true,
		Cmd:          "scw test interrupt error",
		Check: TestCheckCombine(
			TestCheckExitCode(130),
		),
	}))
}
