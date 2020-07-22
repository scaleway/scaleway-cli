package lb

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
)

func GetCommands() *core.Commands {
	human.RegisterMarshalerFunc(lb.LBTypeStock(0), human.EnumMarshalFunc(lbTypeStockMarshalSpecs))
	human.RegisterMarshalerFunc(lb.ACLActionType(0), human.EnumMarshalFunc(aclMarshalSpecs))

	cmds := GetGeneratedCommands()

	cmds.Add(
		lbWaitCommand(),
	)

	cmds.MustFind("lb", "lb", "create").Override(lbCreateBuilder)
	cmds.MustFind("lb", "lb", "get").Override(lbGetBuilder)

	return cmds
}
