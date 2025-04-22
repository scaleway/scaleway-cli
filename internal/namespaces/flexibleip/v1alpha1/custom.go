package flexibleip

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	fip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("fip").Groups = []string{"baremetal"}

	human.RegisterMarshalerFunc(
		fip.FlexibleIPStatus(""),
		human.EnumMarshalFunc(ipStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		fip.MACAddressStatus(""),
		human.EnumMarshalFunc(macAddressStatusMarshalSpecs),
	)

	cmds.MustFind("fip", "ip", "create").Override(createIPBuilder)

	return cmds
}
