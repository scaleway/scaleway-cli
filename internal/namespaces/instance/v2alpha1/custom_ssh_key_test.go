package instance_test

import (
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	instanceV1 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDKV1 "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// addSSHKey add an ssh key to server stored in meta with given key
//
//nolint:unparam
func addSSHKey(serverKey string, sshKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		server := getServerFromMeta(ctx.Meta, serverKey)
		tags := server.Tags
		tags = append(tags, instance.FormatSSHKeyToTag(sshKey))

		_, err := instanceSDK.NewAPI(ctx.Client).
			UpdateServer(&instanceSDK.UpdateServerRequest{
				Zone:     server.Zone,
				ServerID: server.ID,
				Tags:     &tags,
			}, scw.WithContext(ctx.T.Context()))
		if err != nil {
			return err
		}

		// TODO: remove this as soon as the `server create` command is migrated to v2.
		// In the tests, the result that addSSHKey set into the meta will be read by getServerFromMeta, which only handles v1 servers.
		updatedServerV1, err := instanceSDKV1.NewAPI(ctx.Client).
			GetServer(&instanceSDKV1.GetServerRequest{
				Zone:     server.Zone,
				ServerID: server.ID,
			}, scw.WithContext(ctx.T.Context()))
		if err != nil {
			return err
		}

		ctx.Meta[serverKey] = updatedServerV1.Server

		return nil
	}
}

func Test_SSHKey(t *testing.T) {
	cmds := instanceV1AndV2Commands()
	sshKey := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAICC/GwZzzlDeSKN6SliDqfRIUp7u9kDpArZ6Cj+BH1LH key1"
	sshKey2 := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOLfjymSzcwMG56dOUum91KzyuKlf4AI+S1fCmXI8P78 key2"

	t.Run("Add key", core.Test(&core.TestConfig{
		Commands:   cmds,
		BeforeFunc: createServerV1("Server"),
		Args: []string{
			"scw",
			"instance",
			"ssh",
			"add-key",
			"server-id={{.Server.ID}}",
			"public-key=" + sshKey,
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				server := ctx.Meta["Server"].(*instanceV1.ServerWithWarningsResponse).Server
				resp, err := instanceSDK.NewAPI(ctx.Client).GetServer(&instanceSDK.GetServerRequest{
					Zone:     server.Zone,
					ServerID: server.ID,
				}, scw.WithContext(t.Context()))
				require.NoError(t, err)
				assert.Len(t, resp.Tags, 1)
			},
		),
		AfterFunc: deleteServerV1("Server"),
	}))
	t.Run("List keys", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createServerV1("Server"),
			addSSHKey("Server", sshKey),
			addSSHKey("Server", sshKey2),
		),
		Args: []string{"scw", "instance", "ssh", "list-keys", "{{.Server.ID}}"},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				out := string(ctx.Stdout)
				assert.Contains(t, out, "key1")
				assert.Contains(t, out, "key2")
			},
		),
		AfterFunc: deleteServerV1("Server"),
	}))
	t.Run("Remove key", core.Test(&core.TestConfig{
		Commands: cmds,
		BeforeFunc: core.BeforeFuncCombine(
			createServerV1("Server"),
			addSSHKey("Server", sshKey),
			addSSHKey("Server", sshKey2),
		),
		Args: []string{
			"scw",
			"instance",
			"ssh",
			"remove-key",
			"server-id={{.Server.ID}}",
			"name=key2",
		},
		Check: core.TestCheckCombine(
			core.TestCheckGolden(),
			core.TestCheckExitCode(0),
			func(t *testing.T, ctx *core.CheckFuncCtx) {
				t.Helper()
				server := getServerFromMeta(ctx.Meta, "Server")
				resp, err := instanceSDK.NewAPI(ctx.Client).GetServer(&instanceSDK.GetServerRequest{
					Zone:     server.Zone,
					ServerID: server.ID,
				}, scw.WithContext(t.Context()))
				require.NoError(t, err)
				assert.Len(t, resp.Tags, 1)
			},
		),
		AfterFunc: deleteServerV1("Server"),
	}))
}
