package instance

import "github.com/scaleway/scaleway-cli/v2/core"

func CreateSecurityGroup(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw instance security-group create")
}

func DeleteSecurityGroup(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance security-group delete {{ ." + metaKey + ".ID }}")
}
