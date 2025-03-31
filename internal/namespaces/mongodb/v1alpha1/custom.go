package mongodb

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	mongodb "github.com/scaleway/scaleway-sdk-go/api/mongodb/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("mongodb").Groups = []string{"database"}

	human.RegisterMarshalerFunc(
		mongodb.SnapshotStatus(""),
		human.EnumMarshalFunc(snapshotStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		mongodb.InstanceStatus(""),
		human.EnumMarshalFunc(instanceStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		mongodb.NodeTypeStock(""),
		human.EnumMarshalFunc(nodeTypeStockMarshalSpecs),
	)

	cmds.MustFind("mongodb", "instance", "create").Override(instanceCreateBuilder)

	return cmds
}
