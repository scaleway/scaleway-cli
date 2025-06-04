package interlink

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("interlink").Groups = []string{"network"}

	human.RegisterMarshalerFunc(
		interlink.BgpStatus(""),
		human.EnumMarshalFunc(bgpStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		interlink.LinkStatus(""),
		human.EnumMarshalFunc(linkStatusMarshalSpecs),
	)

	return cmds
}
