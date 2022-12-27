package iam

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/internal/core"
	iamsdk "github.com/scaleway/scaleway-sdk-go/api/iam/v1alpha1"
)

func Test_SSHKeyAddCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBBieay3nO9wViPkuvFVgGGaA1IRlkFrr946yqvg9LxZIRhsnZ61yLCPmIOhvUAZ/gTxZGmhgtMDxkenSUTsG3F0= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		Args: []string{
			"scw", "iam", "ssh-key", "create", "name=foobar", "public-key=" + key,
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
		AfterFunc: func(ctx *core.AfterFuncCtx) error {
			api := iamsdk.NewAPI(ctx.Client)
			return api.DeleteSSHKey(&iamsdk.DeleteSSHKeyRequest{
				SSHKeyID: ctx.CmdResult.(*iamsdk.SSHKey).ID,
			})
		},
	}))
}

func Test_SSHKeyRemoveCommand(t *testing.T) {
	key := `ecdsa-sha2-nistp256 AAAAE2VjZHNhLXNoYTItbmlzdHAyNTYAAAAIbmlzdHAyNTYAAABBBGh9rvkJKMu5ljnevB4oRu4i/EnxGS734/UJ6fSDvXGIvT08jIglahc7tM5dvo02abPVXsbiazO25avZZtL6fjo= foobar@foobar`
	t.Run("simple", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: addSSHKey("Key", key),
		Cmd:        "scw iam ssh-key remove {{ .Key.ID }}",
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
		),
	}))
}
