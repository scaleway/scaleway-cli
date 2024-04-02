package llm_inference

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	llm_inference "github.com/scaleway/scaleway-sdk-go/api/llm_inference/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(llm_inference.DeploymentStatus(""), human.EnumMarshalFunc(deployementStateMarshalSpecs))

	human.RegisterMarshalerFunc(llm_inference.Deployment{}, DeploymentMarshalerFunc)
	human.RegisterMarshalerFunc([]*llm_inference.Model{}, ListModelMarshalerFunc)

	cmds.MustFind("llm-inference", "deployment", "create").Override(deploymentCreateBuilder)
	cmds.MustFind("llm-inference", "deployment", "delete").Override(deploymentDeleteBuilder)
	cmds.MustFind("llm-inference", "endpoint", "create").Override(endpointCreateBuilder)

	return cmds
}
