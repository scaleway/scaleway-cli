package file

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("file").Groups = []string{"storage"}

	cmds.MustFind("file", "attachment", "list").Override(core.ListViewBuilder("ID", "FilesystemID", "ResourceID", "ResourceType", "Zone", "Region"))
	cmds.MustFind("file", "filesystem", "list").Override(core.ListViewBuilder("ID", "Name", "Size", "Status", "ProjectID", "OrganizationID", "Tags", "NumberOfAttachments", "Region", "CreatedAt", "UpdatedAt", "FilesystemTypeID"))
	cmds.MustFind("file", "filesystem-type", "list").Override(core.ListViewBuilder("Name", "FilesystemPriceGbPerHour", "SnapshotPriceGbPerHour"))

	return cmds
}
