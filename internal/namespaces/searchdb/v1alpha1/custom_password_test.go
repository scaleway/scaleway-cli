package searchdb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	searchdb "github.com/scaleway/scaleway-cli/v2/internal/namespaces/searchdb/v1alpha1"
	searchdbSDK "github.com/scaleway/scaleway-sdk-go/api/searchdb/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	searchdbNodeType = "SEARCHDB-SHARED-2C-8G"
	searchdbVersion  = "2.0"
	searchdbUser     = "cliuser"
)

func Test_PasswordGenerator(t *testing.T) {
	t.Run("Create deployment", core.Test(&core.TestConfig{
		Commands: searchdb.GetCommands(),
		Cmd: fmt.Sprintf(
			"scw searchdb deployment create name=cli-pwd-gen-create node-type=%s node-count=1 version=%s user-name=%s generate-password=true volume.type=sbs_5k volume.size-bytes=5GB endpoints.0.public={}",
			searchdbNodeType,
			searchdbVersion,
			searchdbUser,
		),
		Check: core.TestCheckExitCode(0),
		AfterFunc: core.AfterFuncCombine(
			func(ctx *core.AfterFuncCtx) error {
				result, ok := ctx.CmdResult.(searchdb.CreateDeploymentResult)
				if !ok || result.Deployment == nil {
					return nil
				}
				api := searchdbSDK.NewAPI(ctx.Client)
				_, err := api.WaitForDeployment(&searchdbSDK.WaitForDeploymentRequest{
					DeploymentID:  result.ID,
					Region:        result.Region,
					Timeout:       scw.TimeDurationPtr(20 * time.Minute),
					RetryInterval: core.DefaultRetryInterval,
				})

				return err
			},
			core.ExecAfterCmd("scw searchdb deployment delete deployment-id={{ .CmdResult.ID }}"),
		),
	}))
}
