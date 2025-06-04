package inference_test

import "github.com/scaleway/scaleway-cli/v2/core"

func CreateDeploymentPublicEndpoint() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"DEPLOYMENT",
		"scw inference deployment create node-type=H100 accept-eula=true model-name=mistral/mistral-7b-instruct-v0.3:bf16 -w",
	)
}

func CreateDeploymentPrivateEndpoint() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"DEPLOYMENT",
		"scw inference deployment create node-type=H100 accept-eula=true model-name=mistral/mistral-7b-instruct-v0.3:bf16 endpoints.0.private-network.private-network-id={{ .PN.ID }} -w",
	)
}

func CreatePN() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"PN",
		"scw vpc private-network create",
	)
}

func DeletePrivateNetwork() core.AfterFunc {
	return core.ExecAfterCmd("scw vpc private-network delete {{ .PN.ID }}")
}

func DeleteDeployment() core.AfterFunc {
	return core.ExecAfterCmd("scw inference deployment delete {{ .DEPLOYMENT.ID }}")
}
