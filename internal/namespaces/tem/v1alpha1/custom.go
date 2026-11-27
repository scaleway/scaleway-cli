package tem

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("tem").Groups = []string{"Domain & WebHosting"}

	// TODO
	// human.RegisterMarshalerFunc(
	// 	tem.DomainStatus(""),
	// 	human.EnumMarshalFunc(domainStatusMarshalSpecs),
	// )
	// human.RegisterMarshalerFunc(
	// 	tem.EmailStatus(""),
	// 	human.EnumMarshalFunc(emailStatusMarshalSpecs))

	cmds.MustFind("tem", "domain", "get").Override(domainGetBuilder)

	return cmds
}
