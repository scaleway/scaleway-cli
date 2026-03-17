package inference_test

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func CreateDeploymentPublicEndpoint() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"DEPLOYMENT",
		fmt.Sprintf(
			"scw inference deployment create node-type-name=%s model-id=%s -w",
			NodeTypeName,
			ModelID,
		),
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
