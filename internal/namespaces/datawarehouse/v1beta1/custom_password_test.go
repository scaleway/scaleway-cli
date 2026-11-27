package datawarehouse_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	datawarehouse "github.com/scaleway/scaleway-cli/v2/internal/namespaces/datawarehouse/v1beta1"
	datawarehouseSDK "github.com/scaleway/scaleway-sdk-go/api/datawarehouse/v1beta1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const dwhVersion = "v26"

func Test_PasswordGenerator(t *testing.T) {
	t.Run("Create deployment", core.Test(&core.TestConfig{
		Commands: datawarehouse.GetCommands(),
		Cmd: fmt.Sprintf(
			"scw datawarehouse deployment create name=cli-pwd-gen-create version=%s replica-count=1 cpu-min=2 cpu-max=4 ram-per-cpu=4 generate-password=true endpoints.0.public={}",
			dwhVersion,
		),
		Check: core.TestCheckCombine(
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				result, ok := ctx.Result.(datawarehouse.CreateDeploymentResult)
				if !ok {
					t.Fatalf("expected CreateDeploymentResult, got %T", ctx.Result)
				}
				if result.Password == "" {
					t.Fatal("expected generated password in result")
				}
			},
		),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				result, ok := ctx.CmdResult.(datawarehouse.CreateDeploymentResult)
				if !ok || result.Deployment == nil {
					return nil
				}
				api := datawarehouseSDK.NewAPI(ctx.Client)
				_, err := api.WaitForDeployment(&datawarehouseSDK.WaitForDeploymentRequest{
					DeploymentID:  result.ID,
					Region:        result.Region,
					Timeout:       scw.TimeDurationPtr(20 * time.Minute),
					RetryInterval: core.DefaultRetryInterval,
				})

				return err
			},
			core.ExecAfterCmd("scw datawarehouse deployment delete {{ .CmdResult.ID }}"),
		),
	}))
}
