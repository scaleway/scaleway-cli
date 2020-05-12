package core

import (
	"context"
	"reflect"
	"testing"

	"github.com/scaleway/scaleway-cli/internal/args"
	"github.com/scaleway/scaleway-sdk-go/logger"
)

func testLogCommands() *Commands {
	return NewCommands(
		&Command{
			Namespace:            "log",
			ArgsType:             reflect.TypeOf(args.RawArgs{}),
			AllowAnonymousClient: true,
			Run: func(ctx context.Context, argsI interface{}) (i interface{}, e error) {
				return "Hello World!", nil
			},
		},
	)
}

func Test_Logger(t *testing.T) {
	t.Run("Simple", Test(&TestConfig{
		Commands: testLogCommands(),
		BeforeFunc: func(ctx *BeforeFuncCtx) error {
			ctx.Logger.level = logger.LogLevelDebug
			ctx.Logger.Debug("My debug message from before func")
			ctx.Logger.Info("My info message from before func")
			ctx.Logger.Warning("My warning message from before func")
			ctx.Logger.Error("My error message from before func")
			return nil
		},
		AfterFunc: func(ctx *AfterFuncCtx) error {
			ctx.Logger.Debug("My debug message from after func")
			ctx.Logger.Info("My info message from after func")
			ctx.Logger.Warning("My warning message from after func")
			ctx.Logger.Error("My error message from after func")
			return nil
		},
		Cmd: "scw log",
		Check: TestCheckCombine(
			TestCheckGolden(),
			TestCheckExitCode(0),
		),
	}))
}
