package datalab

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("datalab", "cluster-version", "list").Override(core.ListViewBuilder("Name", "Description", "Versions"))
	cmds.MustFind("datalab", "datalab", "list").Override(core.ListViewBuilder("ID", "ProjectID", "Name", "Description", "Tags", "Status", "CreatedAt", "UpdatedAt", "Region", "HasNotebook", "NotebookURL", "SparkVersion", "PrivateNetworkID", "NotebookMasterURL"))
	cmds.MustFind("datalab", "node-type", "list").Override(core.ListViewBuilder("StockStatus", "Name", "Description", "Vcpus", "MemoryGigabytes", "VramGigabytes", "Gpus", "Disabled", "Beta", "CreatedAt", "UpdatedAt", "Targets"))
	cmds.MustFind("datalab", "notebook-version", "list").Override(core.ListViewBuilder("Name", "Description", "Versions"))

	return cmds
}
