package instance

import (
	"testing"

	"github.com/alecthomas/assert"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

// addSSHKey add an ssh key to server stored in meta with given key
//
//nolint:unparam
func addSSHKey(serverKey string, sshKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		server := ctx.Meta[serverKey].(*instance.Server)
		tags := append(server.Tags, formatSSHKeyToTag(sshKey))

		resp, err := instance.NewAPI(ctx.Client).UpdateServer(&instance.UpdateServerRequest{
			Zone:     server.Zone,
			ServerID: server.ID,
			Tags:     &tags,
		})
		if err != nil {
			return err
		}

		ctx.Meta[serverKey] = resp.Server

		return nil
	}
}

func Test_SSHKey(t *testing.T) {
	sshKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICC/GwZzzlDeSKN6SliDqfRIUp7u9kDpArZ6Cj+BH1LH key1"
	sshKey2 := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOLfjymSzcwMG56dOUum91KzyuKlf4AI+S1fCmXI8P78 key2"

	t.Run("Add key", core.Test(&core.TestConfig{
		Commands:   GetCommands(),
		BeforeFunc: createServerBionic("Server"),
		Args:       []string{"scw", "instance", "ssh", "add-key", "server-id={{.Server.ID}}", "public-key=" + sshKey},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				server := ctx.Meta["Server"].(*instance.Server)
				resp, err := instance.NewAPI(ctx.Client).GetServer(&instance.GetServerRequest{
					Zone:     server.Zone,
					ServerID: server.ID,
				})
				assert.Nil(t, err)
				assert.Len(t, resp.Server.Tags, 1)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
	t.Run("List keys", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerBionic("Server"),
			addSSHKey("Server", sshKey),
			addSSHKey("Server", sshKey2),
		),
		Args: []string{"scw", "instance", "ssh", "list-keys", "{{.Server.ID}}"},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				out := string(ctx.Stdout)
				assert.Contains(t, out, "key1")
				assert.Contains(t, out, "key2")
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
	t.Run("Remove key", core.Test(&core.TestConfig{
		Commands: GetCommands(),
		BeforeFunc: core.BeforeFuncCombine(
			createServerBionic("Server"),
			addSSHKey("Server", sshKey),
			addSSHKey("Server", sshKey2),
		),
		Args: []string{"scw", "instance", "ssh", "remove-key", "server-id={{.Server.ID}}", "name=key2"},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				server := ctx.Meta["Server"].(*instance.Server)
				resp, err := instance.NewAPI(ctx.Client).GetServer(&instance.GetServerRequest{
					Zone:     server.Zone,
					ServerID: server.ID,
				})
				assert.Nil(t, err)
				assert.Len(t, resp.Server.Tags, 1)
			},
		),
		AfterFunc: deleteServer("Server"),
	}))
}
