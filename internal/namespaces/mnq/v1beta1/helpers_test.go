package mnq_test

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func createNATSAccount(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		"scw mnq nats create-account")
}

func deleteNATSAccount(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw mnq nats delete-account {{ ." + metaKey + ".ID }}")
}
