package flexibleip

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(flexibleip.FlexibleIPStatus(""), human.EnumMarshalFunc(ipStatusMarshalSpecs))
	human.RegisterMarshalerFunc(flexibleip.MACAddressStatus(""), human.EnumMarshalFunc(macAddressStatusMarshalSpecs))

	return cmds
}
