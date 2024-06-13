package flexibleip

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	fip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(fip.FlexibleIPStatus(""), human.EnumMarshalFunc(ipStatusMarshalSpecs))
	human.RegisterMarshalerFunc(fip.MACAddressStatus(""), human.EnumMarshalFunc(macAddressStatusMarshalSpecs))

	cmds.MustFind("fip", "ip", "create").Override(createIPBuilder)

	return cmds
}
