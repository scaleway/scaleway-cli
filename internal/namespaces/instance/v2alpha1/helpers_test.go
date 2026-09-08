package instance_test

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Server
//

// createServer creates a stopped ubuntu server without IP and
// register it in the context Meta at given metaKey
func createServerV1(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, testServerCommand("stopped=true image=ubuntu-jammy"))
}

// testServerCommand creates returns a create server command with the instance type and the given arguments
func testServerCommand(params string) string {
	baseCommand := "scw instance server create "
	if !strings.Contains(params, "ip=") {
		baseCommand += "ip=none "
	}
	if !strings.Contains(params, "image=") {
		baseCommand += "image=ubuntu_jammy "
	}
	if !strings.Contains(params, "type=") {
		baseCommand += "type=DEV1-S "
	}

	return baseCommand + params
}

func getServerFromMeta(meta core.TestMetadata, metaKey string) *instanceSDK.Server {
	switch resp := meta[metaKey].(type) {
	case *instanceSDK.Server:
		return resp
	case *instance.ServerWithWarningsResponse:
		return resp.Server
	default:
		return nil
	}
}

// deleteServer deletes a server and its attached IP and volumes
// previously registered in the context Meta at metaKey.
func deleteServerV1(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		server := getServerFromMeta(ctx.Meta, metaKey)
		if server.State == instanceSDK.ServerStateRunning {
			err := core.ExecAfterCmd("scw instance server stop -w {{ ." + metaKey + ".ID }}")(ctx)
			if err != nil {
				return err
			}
		}

		return core.ExecAfterCmd(
			"scw instance server delete {{ ." + metaKey + ".ID }} with-ip=true with-volumes=all",
		)(
			ctx,
		)
	}
}

//
// Private Network Interface
//

//nolint:unparam
func createPrivateNetworkInterface(serverMetaKey, privateNetworkMetaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"PNI",
		fmt.Sprintf(
			"scw instance private-network-interface create server-id={{ .%s.ID }} private-network-id={{ .%s.ID }}",
			serverMetaKey,
			privateNetworkMetaKey,
		),
	)
}
