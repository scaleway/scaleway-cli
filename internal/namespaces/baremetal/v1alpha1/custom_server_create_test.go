package baremetal

import (
	"fmt"
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

// All test below should succeed to create an instance.
func Test_CreateServer(t *testing.T) {
	// Simple use cases
	t.Run("Simple", func(t *testing.T) {
		t.Run("Default", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create",
			Check: core.TestCheckCombine(
				core.TestCheckGolden(),
				core.TestCheckExitCode(0),
			),
			AfterFunc: func(ctx *core.AfterFuncCtx) error {
				_, err := baremetal.NewAPI(ctx.Client).WaitForServer(&baremetal.WaitForServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}
				_, err = baremetal.NewAPI(ctx.Client).DeleteServer(&baremetal.DeleteServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}
				return nil
			},
			DefaultZone: scw.ZoneFrPar2,
		}))

		t.Run("With name", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create name=test-create-server-with-name",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "test-create-server-with-name", ctx.Result.(*baremetal.Server).Name)
				},
				core.TestCheckExitCode(0),
			),
			DefaultZone: scw.ZoneFrPar2,
			AfterFunc: func(ctx *core.AfterFuncCtx) error {
				_, err := baremetal.NewAPI(ctx.Client).WaitForServer(&baremetal.WaitForServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}

				_, err = baremetal.NewAPI(ctx.Client).DeleteServer(&baremetal.DeleteServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}

				return nil
			},
		}))

		t.Run("Tags", core.Test(&core.TestConfig{
			Commands: GetCommands(),
			Cmd:      "scw baremetal server create tags.0=prod tags.1=blue",
			Check: core.TestCheckCombine(
				func(t *testing.T, ctx *core.CheckFuncCtx) {
					assert.Equal(t, "prod", ctx.Result.(*baremetal.Server).Tags[0])
					assert.Equal(t, "blue", ctx.Result.(*baremetal.Server).Tags[1])
				},
				core.TestCheckExitCode(0),
			),
			DefaultZone: scw.ZoneFrPar2,
			AfterFunc: func(ctx *core.AfterFuncCtx) error {
				_, err := baremetal.NewAPI(ctx.Client).WaitForServer(&baremetal.WaitForServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}

				_, err = baremetal.NewAPI(ctx.Client).DeleteServer(&baremetal.DeleteServerRequest{
					ServerID: ctx.CmdResult.(*baremetal.Server).ID,
				})
				if err != nil {
					return err
				}

				return nil
			},
		}))
	})
}

// None of the tests below should succeed to create an instance.
func Test_CreateServerErrors(t *testing.T) {
	////
	// Instance type errors
	////
	t.Run("Error: invalid instance type", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Cmd:      "scw baremetal server create type=foobar",
		Check: core.TestCheckCombine(
			core.TestCheckError(&core.CliError{
				Err:  fmt.Errorf("invalid value 'foobar' for arg 'type'"),
				Hint: "Accepted values for 'type' are [GP-BM1-L GP-BM1-M GP-BM1-S HC-BM1-L HC-BM1-S HM-BM1-XL HM-BM1-M]",
			}),
			core.TestCheckExitCode(1),
		),
	}))
}
