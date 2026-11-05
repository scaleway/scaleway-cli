package iot

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/iot/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("iot").Groups = []string{"integration"}

	humanv2.RegisterEnumFormat(hubStatusFormat)
	humanv2.RegisterEnumFormat(deviceMessageFiltersRulePolicyFormat)
	humanv2.RegisterEnumFormat(deviceStatusFormat)
	humanv2.RegisterLayout(reflect.TypeFor[iot.CreateNetworkResponse](), iotNetworkCreateLayout)
	humanv2.RegisterLayout(reflect.TypeFor[iot.CreateDeviceResponse](), iotDeviceCreateLayout)

	cmds.MustFind("iot", "hub", "create").Override(hubCreateBuilder)

	return cmds
}
