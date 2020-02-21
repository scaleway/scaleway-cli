package instance

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Server
//

// createServer creates a stopped ubuntu-bionic server and
// register it in the context Meta at metaKey.
func createServer(metaKey string) core.BeforeFunc {
	return core.ExecStoreBeforeCmd(metaKey, "scw instance server create stopped=true image=ubuntu-bionic")
}

// deleteServer deletes a server and its attached IP and volumes
// previously registered in the context Meta at metaKey.
func deleteServer(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance server delete server-id={{ ." + metaKey + ".ID }} delete-ip=true delete-volumes=true")
}

//
// Volume
//

// createVolume creates a volume of the given size and type and
// register it in the context Meta at metaKey.
func createVolume(metaKey string, sizeInGb int, volumeType instance.VolumeType) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		cmd := fmt.Sprintf("scw instance volume create name=cli-test size=%dGB volume-type=%s", sizeInGb, volumeType)
		res := ctx.ExecuteCmd(cmd)
		createVolumeResponse := res.(*instance.CreateVolumeResponse)
		ctx.Meta[metaKey] = createVolumeResponse.Volume
		return nil
	}
}

// deleteVolume deletes a volume previously registered in the context Meta at metaKey.
func deleteVolume(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance volume delete volume-id={{ ." + metaKey + ".ID }}")
}

//
// IP
//

// createIP creates an IP and register it in the context Meta at metaKey.
func createIP(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd("scw instance ip create")
		createIPResponse := res.(*instance.CreateIPResponse)
		ctx.Meta[metaKey] = createIPResponse.IP
		return nil
	}
}

// deleteIP deletes an IP previously registered in the context Meta at metaKey.
func deleteIP(metaKey string) core.AfterFunc {
	return func(ctx *core.AfterFuncCtx) error {
		ctx.ExecuteCmd("scw instance ip delete ip={{ ." + metaKey + ".Address }}")
		return nil
	}
}

//
// Placement Group
//

// createPlacementGroup creates a placement group and
// register it in the context Meta at metaKey.
func createPlacementGroup(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd("scw instance placement-group create")
		createPlacementGroupResponse := res.(*instance.CreatePlacementGroupResponse)
		ctx.Meta[metaKey] = createPlacementGroupResponse.PlacementGroup
		return nil
	}
}

// deletePlacementGroup deletes a placement group
// previously registered in the context Meta at metaKey.
func deletePlacementGroup(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance placement-group delete placement-group-id={{ ." + metaKey + ".ID }}")
}

//
// Security Group
//

// createSecurityGroup creates a security group and
// register it in the context Meta at metaKey.
func createSecurityGroup(metaKey string) core.BeforeFunc {
	return func(ctx *core.BeforeFuncCtx) error {
		res := ctx.ExecuteCmd("scw instance security-group create")
		createSecurityGroupResponse := res.(*instance.CreateSecurityGroupResponse)
		ctx.Meta[metaKey] = createSecurityGroupResponse.SecurityGroup
		return nil
	}
}

// deleteSecurityGroup deletes a security group
// previously registered in the context Meta at metaKey.
func deleteSecurityGroup(metaKey string) core.AfterFunc {
	return core.ExecAfterCmd("scw instance security-group delete security-group-id={{ ." + metaKey + ".ID }}")
}
