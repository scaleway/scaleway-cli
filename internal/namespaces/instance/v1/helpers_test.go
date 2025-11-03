package instance_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/internal/namespaces/instance/v1"
	"github.com/scaleway/scaleway-cli/v2/internal/testhelpers"
	block "github.com/scaleway/scaleway-sdk-go/api/block/v1alpha1"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
	"github.com/stretchr/testify/require"
)

//
// Server
//

// createServer creates a stopped ubuntu server without IP and
// register it in the context Meta at given metaKey
//
//nolint:unparam
func createServer(metaKey string) core.BeforeFunc {
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

// createServer creates a stopped ubuntu-bionic server and
// register it in the context Meta at metaKey.
func startServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw instance server start -w {{ ."+metaKey+
		".ID }}")
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
//
//nolint:unparam
func deleteServer(metaKey string) core.AfterFunc {
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
// Volume
//

// createVolume creates a volume of the given size and type and
// register it in the context Meta at metaKey.
//
//nolint:unparam
func createVolume(
	metaKey string,
	sizeInGb int,
	volumeType instanceSDK.VolumeVolumeType,
) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf(
			"scw instance volume create name=cli-test size=%dGB volume-type=%s",
			sizeInGb,
			volumeType,
		)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		createVolumeResponse := res.(*instanceSDK.CreateVolumeResponse)
		ctx.Meta[metaKey] = createVolumeResponse.Volume

		return nil
	}
}

// deleteVolume deletes a volume previously registered in the context Meta at metaKey.
func deleteVolume(metaKey string) core.AfterFunc { //nolint: unparam
	return core.ExecAfterCmd("scw instance volume delete {{ ." + metaKey + ".ID }}")
}

// createSbsVolume creates a volume of the given size and
// register it in the context Meta at metaKey
//
//nolint:unparam
func createSbsVolume(metaKey string, sizeInGb int) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf(
			"scw block volume create name=%s from-empty.size=%dGB perf-iops=5000 -w",
			ctx.T.Name(),
			sizeInGb,
		)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		volume := res.(*block.Volume)
		ctx.Meta[metaKey] = volume

		return nil
	}
}

// createNonEmptyLocalVolume creates a server with a local root volume of the given size and registers the volume in
// the context Meta at metaKey. The volume is then detached, and the server deleted, leaving a non-empty volume
// ready to be snapshot, or any other use case that requires a non-empty local volume.
func createNonEmptyLocalVolume(metaKey string, sizeInGB int) core.BeforeFunc {
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

//
// IP
//

// createIP creates an IP and register it in the context Meta at metaKey.
func createIP(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd(strings.Split("scw instance ip create", " "))
		createIPResponse := res.(*instanceSDK.CreateIPResponse)
		ctx.Meta[metaKey] = createIPResponse.IP

		return nil
	}
}

// deleteIP deletes an IP previously registered in the context Meta at metaKey.
func deleteIP(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance ip delete {{ ." + metaKey + ".Address }}")
}

//
// Placement Group
//

// createPlacementGroup creates a placement group and
// register it in the context Meta at metaKey.
func createPlacementGroup(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd([]string{"scw", "instance", "placement-group", "create"})
		createPlacementGroupResponse := res.(*instanceSDK.CreatePlacementGroupResponse)
		ctx.Meta[metaKey] = createPlacementGroupResponse.PlacementGroup

		return nil
	}
}

// deletePlacementGroup deletes a placement group
// previously registered in the context Meta at metaKey.
func deletePlacementGroup(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance placement-group delete {{ ." + metaKey + ".ID }}")
}

//
// Security Group
//

// createSecurityGroup creates a security group and
// register it in the context Meta at metaKey.
func createSecurityGroup(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd([]string{"scw", "instance", "security-group", "create"})
		createSecurityGroupResponse := res.(*instanceSDK.CreateSecurityGroupResponse)
		ctx.Meta[metaKey] = createSecurityGroupResponse.SecurityGroup

		return nil
	}
}

// deleteSecurityGroup deletes a security group
// previously registered in the context Meta at metaKey.
func deleteSecurityGroup(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance security-group delete {{ ." + metaKey + ".ID }}")
}

//
// Snapshot
//

// deleteSnapshot deletes a snapshot previously registered in the context Meta at metaKey.
func deleteSnapshot(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance snapshot delete {{ ." + metaKey + ".Snapshot.ID }}")
}

// deleteSnapshot deletes a snapshot previously registered in the context Meta at metaKey.
func deleteBlockSnapshot(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw block snapshot delete {{ ." + metaKey + ".ID }}")
}

func createNIC() core.BeforeFunc {
	return core.ExecStoreBeforeCmd(
		"NIC",
		"scw instance private-nic create server-id={{ .Server.ID }} private-network-id={{ .PN.ID }}",
	)
}

// testServerSBSVolumeSize checks the size of a volume in an instance server.
// The server must be returned by the given instanceFetcher function
func testServerFetcherSBSVolumeSize(
	volumeKey string,
	sizeInGB int,
	serverFetcher func(t *testing.T, ctx *core.CheckFuncCtx) *instanceSDK.Server,
) core.TestCheck {
	return func(t *testing.T, ctx *core.CheckFuncCtx) {
		t.Helper()
		server := serverFetcher(t, ctx)
		blockAPI := block.NewAPI(ctx.Client)
		serverVolume := testhelpers.MapTValue(t, server.Volumes, volumeKey)
		volume, err := blockAPI.GetVolume(&block.GetVolumeRequest{
			Zone:     server.Zone,
			VolumeID: serverVolume.ID,
		})
		require.NoError(t, err)

		require.Equal(
			t,
			scw.Size(sizeInGB)*scw.GB,
			volume.Size,
			"Size of volume should be %d GB",
			sizeInGB,
		)
	}
}

// testServerSBSVolumeSize checks the size of a volume in Result's server.
// The server must be returned as result of the test's Cmd
func testServerSBSVolumeSize(volumeKey string, sizeInGB int) core.TestCheck {
	return testServerFetcherSBSVolumeSize(
		volumeKey,
		sizeInGB,
		func(t *testing.T, ctx *core.CheckFuncCtx) *instanceSDK.Server {
			t.Helper()

			return testhelpers.Value[*instance.ServerWithWarningsResponse](t, ctx.Result).Server
		},
	)
}

// testAttachVolumeServerSBSVolumeSize is the same as testServerSBSVolumeSize but the test's Cmd must be "scw instance server attach-volume"
func testAttachVolumeServerSBSVolumeSize(volumeKey string, sizeInGB int) core.TestCheck {
	return testServerFetcherSBSVolumeSize(
		volumeKey,
		sizeInGB,
		func(t *testing.T, ctx *core.CheckFuncCtx) *instanceSDK.Server {
			t.Helper()

			return testhelpers.Value[*instanceSDK.AttachVolumeResponse](t, ctx.Result).Server
		},
	)
}

func testServerUpdateServerSBSVolumeSize(volumeKey string, sizeInGB int) core.TestCheck {
	return testServerFetcherSBSVolumeSize(
		volumeKey,
		sizeInGB,
		func(t *testing.T, ctx *core.CheckFuncCtx) *instanceSDK.Server {
			t.Helper()

			return testhelpers.Value[*instance.ServerWithWarningsResponse](t, ctx.Result).Server
		},
	)
}
