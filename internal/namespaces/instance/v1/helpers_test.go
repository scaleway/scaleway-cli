package instance

import (
	"fmt"
	"strings"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Server
//

// createServer creates a stopped ubuntu-bionic server and
// register it in the context Meta at metaKey.
// nolint:unparam
func createServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw instance server create stopped=true image=ubuntu-bionic")
}

// createServer creates a stopped ubuntu-bionic server and
// register it in the context Meta at metaKey.
func startServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw instance server start -w {{ ."+metaKey+".ID }}")
}

// deleteServer deletes a server and its attached IP and volumes
// previously registered in the context Meta at metaKey.
// nolint:unparam
func deleteServer(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		server := ctx.Meta[metaKey].(*instance.Server)
		if server.State == instance.ServerStateRunning {
			err := core.ExecAfterCmd("scw instance server stop -w {{ ." + metaKey + ".ID }}")(ctx)
			if err != nil {
				return err
			}
		}
		return core.ExecAfterCmd("scw instance server delete {{ ." + metaKey + ".ID }} with-ip=true with-volumes=all")(ctx)
	}
}

//
// Volume
//

// createVolume creates a volume of the given size and type and
// register it in the context Meta at metaKey.
// nolint:unparam
func createVolume(metaKey string, sizeInGb int, volumeType instance.VolumeVolumeType) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw instance volume create name=cli-test size=%dGB volume-type=%s", sizeInGb, volumeType)
		res := ctx.ExecuteCmd(strings.Split(cmd, " "))
		createVolumeResponse := res.(*instance.CreateVolumeResponse)
		ctx.Meta[metaKey] = createVolumeResponse.Volume
		return nil
	}
}

// deleteVolume deletes a volume previously registered in the context Meta at metaKey.
func deleteVolume(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance volume delete {{ ." + metaKey + ".ID }}")
}

//
// IP
//

// createIP creates an IP and register it in the context Meta at metaKey.
func createIP(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd(strings.Split("scw instance ip create", " "))
		createIPResponse := res.(*instance.CreateIPResponse)
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
		createPlacementGroupResponse := res.(*instance.CreatePlacementGroupResponse)
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
		createSecurityGroupResponse := res.(*instance.CreateSecurityGroupResponse)
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
