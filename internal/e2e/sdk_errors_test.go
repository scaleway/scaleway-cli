package e2e_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/test/v1"
	sdktest "github.com/scaleway/scaleway-sdk-go/api/test/v1"
	"github.com/stretchr/testify/assert"
)

// TestSdkStandardErrors tests standard errors
//
// Some errors ar not tested on purpose:
// - InvalidField: this error is deprecated
// - PermissionsDenied: this error cannot be triggered using the SDK
func TestSdkStandardErrors(t *testing.T) {
	t.Skip("Test API not available")

	t.Run("invalid-arguments", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human create altitude-in-meter=-7000000",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))

	t.Run("quotas-exceeded", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			for range 10 {
				ctx.ExecuteCmd([]string{"scw", "test", "human", "create"})
			}

			return nil
		},
		Cmd: "scw test human create",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))

	t.Run("transient-state", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		BeforeFunc: func(ctx *core.BeforeFuncCtx) error {
			ctx.ExecuteCmd([]string{"scw", "test", "human", "create"})
			api := sdktest.NewAPI(ctx.Client)
			_, err := api.RunHuman(&sdktest.RunHumanRequest{
				HumanID: "0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
			})
			assert.NoError(t, err)

			return nil
		},
		Cmd: "scw test human update human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8 eyes-color=red",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))

	t.Run("resource-not-found", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human get human-id=0194fdc2-fa2f-fcc0-41d3-ff12045b73c8",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))

	t.Run("out-of-stock", core.Test(&core.TestConfig{
		Commands:        test.GetCommands(),
		UseE2EClient:    true,
		DisableParallel: true, // because e2e client is used
		Cmd:             "scw test human create shoe-size=60",
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(1),
			core.TestCheckGolden(),
		),
	}))
}
