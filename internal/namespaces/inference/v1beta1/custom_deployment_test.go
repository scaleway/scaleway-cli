package inference_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	inference "github.com/scaleway/scaleway-cli/v2/internal/namespaces/inference/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_DeploymentCreate(t *testing.T) {
	cmds := inference.GetCommands()

	t.Run("Simple deployment", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw inference deployment create node-type=L4 model-name=meta/llama-3-8b-instruct:bf16 accept-eula=true",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw inference deployment delete {{ .CmdResult.ID }}"),
	}))

	t.Run("Deployment with wait flag", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw inference deployment create node-type=L4 model-name=meta/llama-3-8b-instruct:bf16 accept-eula=true --wait",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw inference deployment delete {{ .CmdResult.ID }}"),
	}))
}

func Test_CreateDeploymentPrivateEndpoint(t *testing.T) {
	t.Skip("Out of stock")
	cmds := inference.GetCommands()
	cmds.Merge(vpc.GetCommands())
	t.Run("Create Deployment Private Endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreatePN(),
		Cmd:        "scw inference deployment create model-name=meta/llama-2-7b-chat:fp16 node-type=L4 accept-eula=true endpoints.0.private-network.private-network-id={{ .PN.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw inference endpoint delete {{ .CmdResult.ID }}"),
			DeletePrivateNetwork(),
			DeleteDeployment(),
		),
	}))
}

func Test_DeploymentDelete(t *testing.T) {
	t.Skip("No stock to run test")
	cmds := inference.GetCommands()
	t.Run("Delete deployment with wait flag", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreateDeploymentPublicEndpoint(),
		Cmd:        "scw inference deployment delete {{ .DEPLOYMENT.ID }} -w",
		Check:      core.TestCheckGolden(),
	}))
}
