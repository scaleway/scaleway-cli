package secret_test

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func createSecret(name string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"Secret",
		"scw secret secret create name="+name,
	)
}

func createSecretVersion(content string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"SecretVersion",
		"scw secret version create {{ .Secret.ID }} data="+content,
	)
}

func deleteSecret() core.AfterFunc {
	return core.ExecAfterCmd("scw secret secret delete {{ .Secret.ID }}")
}
