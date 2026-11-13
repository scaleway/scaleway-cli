package instance_test

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	instanceV2 "github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v2alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Server
//

func instanceV1AndV2Commands() *core.Commands {
	cmds := instance.GetCommands()
	cmds.Merge(instanceV2.GetCommands())

	return cmds
}

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
			err := core.ExecAfterCmd(
				fmt.Sprintf("scw instance server stop -w {{ .%s.ID }}", metaKey),
			)(
				ctx,
			)
			if err != nil {
				return err
			}
		}

		return core.ExecAfterCmd(fmt.Sprintf(
			"scw instance server delete {{ .%s.ID }} with-ip=true with-volumes=all", metaKey),
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

//
// Snapshot
//

// createSnapshot create an instance snapshot from a local volume at volumeMetaKey, and stores it at metaKey.
func createSnapshot(metaKey, volumeMetaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		metaKey,
		fmt.Sprintf(
			"scw instance snapshot create base-volume-id={{ .%s.ID }}", volumeMetaKey,
		),
	)
}

// deleteSnapshot deletes a snapshot previously registered in the context Meta at metaKey.
func deleteSnapshot(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(
		fmt.Sprintf("scw instance snapshot delete {{ .%s.ID }}", metaKey),
	)
}

//
// Volume
//

// createNonEmptyLocalVolume creates a server with a local root volume of the given size and registers the volume in
// the context Meta at metaKey. The volume is then detached, and the server deleted, leaving a non-empty volume
// ready to be snapshot, or any other use case that requires a non-empty local volume.
func createNonEmptyLocalVolumeV1(metaKey string, sizeInGB int) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf(
			"scw instance server create type=DEV1-S root-volume=local:%dGB stopped=true",
			sizeInGB,
		)
		server := ctx.ExecuteCmd(strings.Split(cmd, " "))
		createServerResponse := server.(*instance.ServerWithWarningsResponse)
		serverID := createServerResponse.Server.ID
		volume := createServerResponse.Server.Volumes["0"]
		ctx.Meta[metaKey] = volume

		cmd = "scw instance server detach-volume volume-id=" + volume.ID + " server-id=" + serverID
		_ = ctx.ExecuteCmd(strings.Split(cmd, " "))

		cmd = "scw instance server delete " + serverID
		_ = ctx.ExecuteCmd(strings.Split(cmd, " "))

		return nil
	}
}

// deleteVolume deletes a volume previously registered in the context Meta at "Volume".
func deleteVolumeV1(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd(fmt.Sprintf("scw instance volume delete {{ .%s.ID }}", metaKey))
}
