package llm_inference_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	llm_inference "github.com/scaleway/scaleway-cli/v2/internal/namespaces/llm_inference/v1beta1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/vpc/v2"
)

func Test_DeploymentCreate(t *testing.T) {
	t.Skip("Out of stock")
	cmds := llm_inference.GetCommands()

	t.Run("Single public endpoint", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw llm-inference deployment create node-type=H100 model-name=wizardlm/wizardlm-70b-v1.0:fp8 accept-eula=true",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw llm-inference deployment delete {{ .CmdResult.ID }}"),
	}))

	t.Run("Deployment with wait flag", core.Test(&core.TestConfig{
		Commands:  cmds,
		Cmd:       "scw llm-inference deployment create model-name=meta/llama-2-7b-chat:fp16 node-type=L4 accept-eula=true -w",
		Check:     core.TestCheckGolden(),
		AfterFunc: core.ExecAfterCmd("scw llm-inference deployment delete {{ .CmdResult.ID }}"),
	}))
}

func Test_CreateDeploymentPrivateEndpoint(t *testing.T) {
	t.Skip("Out of stock")
	cmds := llm_inference.GetCommands()
	cmds.Merge(vpc.GetCommands())
	t.Run("Create Deployment Private Endpoint", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreatePN(),
		Cmd:        "scw llm-inference deployment create model-name=meta/llama-2-7b-chat:fp16 node-type=L4 accept-eula=true endpoints.0.private-network.private-network-id={{ .PN.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
		),
		AfterFunc: core.AfterFuncCombine(
			core.ExecAfterCmd("scw llm-inference endpoint delete {{ .CmdResult.ID }}"),
			DeletePrivateNetwork(),
			DeleteDeployment(),
		),
	}))
}

func Test_DeploymentDelete(t *testing.T) {
	t.Skip("No stock to run test")
	cmds := llm_inference.GetCommands()
	t.Run("Delete deployment with wait flag", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: CreateDeploymentPublicEndpoint(),
		Cmd:        "scw llm-inference deployment delete {{ .DEPLOYMENT.ID }} -w",
		Check:      core.TestCheckGolden(),
	}))
}
