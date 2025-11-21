package inference_test

import (
	"fmt"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	inference "github.com/scaleway/scaleway-cli/v2/internal/namespaces/inference/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

const (
	ModelID      = "739d51ae-4f1e-4193-a4bf-f7380c090d46"
	NodeTypeName = "H100-2"
)

func Test_DeploymentCreate(t *testing.T) {
	cmds := inference.GetCommands()

	t.Run("Simple deployment", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd: fmt.Sprintf(
			"scw inference deployment create node-type-name=%s model-id=%s",
			NodeTypeName,
			ModelID,
		),
		Check: core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd(
			"scw inference deployment delete {{ .CmdResult.ID }}",
		),
	}))

	t.Run("Deployment with wait flag", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd: fmt.Sprintf(
			"scw inference deployment create node-type-name=%s model-id=%s accept-eula=true --wait",
			NodeTypeName, ModelID,
		),
		Check: core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd(
			"scw inference deployment delete {{ .CmdResult.ID }}",
		),
	}))

	t.Run("Deployment with no endpoints must fail", core.Test(&core.TestConfig{
		Commands: cmds,
		Cmd: fmt.Sprintf(
			"scw inference deployment create node-type-name=%s model-id=%s endpoints.0.is-public=false",
			NodeTypeName,
			ModelID,
		),
		Check: core.TestCheckGolden(),
	}))
}

func Test_CreateDeploymentPrivateEndpoint(t *testing.T) {
	cmds := inference.GetCommands()
	cmds.Merge(vpc.GetCommands())

	t.Run("Create Deployment Private Endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreatePN(),
		Cmd: fmt.Sprintf(
			"scw inference deployment create model-id=%s node-type-name=H100-SXM-2 accept-eula=true endpoints.0.private-network.private-network-id={{ .PN.ID }}",
			ModelID,
		),
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw inference deployment delete {{ .CmdResult.ID }} --wait"),
			DeletePrivateNetwork(),
		),
	}))
}

func Test_DeploymentDelete(t *testing.T) {
	cmds := inference.GetCommands()

	t.Run("Delete deployment with wait flag", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreateDeploymentPublicEndpoint(),
		Cmd:        "scw inference deployment delete {{ .DEPLOYMENT.ID }} --wait",
		Check:      core.TestCheckGolden(),
	}))
}
