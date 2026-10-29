package instance

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	instanceSDK "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
)

// GetCommands returns instance commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run and Command.View)
// - Merge handwritten commands
func GetCommands() *core.Commands {
	cmds := core.NewCommands(
		instancePlacementGroup(),
		instanceServerType(),
		instanceTemplate(),
		instanceUserData(),
		instancePlacementGroupList(),
		instancePlacementGroupCreate(),
		instancePlacementGroupGet(),
		instancePlacementGroupUpdate(),
		instancePlacementGroupDelete(),
		instanceServerGetCompatibleTypes(),
		instanceServerTypeList(),
		instanceTemplateList(),
		instanceTemplateCreate(),
		instanceTemplateGet(),
		instanceTemplateUpdate(),
		instanceTemplateDelete(),
		instanceTemplateListUserDataKeys(),
		instanceTemplateGetUserData(),
		instanceTemplateSetUserData(),
		instanceTemplateDeleteUserData(),
		instanceTemplateGetCloudInit(),
		instanceTemplateSetCloudInit(),
		instanceTemplateCheck(),
		instanceTemplateCreateServer(),
		instanceUserDataList(),
		instanceUserDataSet(),
		instanceUserDataGet(),
		instanceUserDataDelete(),
		instanceUserDataGetCloudInit(),
		instanceUserDataSetCloudInit(),
	)

	//
	// Placement Group
	//

	cmds.MustFind("instance", "placement-group", "create").Override(placementGroupCreateBuilder)
	cmds.MustFind("instance", "placement-group", "get").Override(placementGroupGetBuilder)
	cmds.MustFind("instance", "placement-group", "list").Override(placementGroupListBuilder)

	//
	// Server-Type
	//
	human.RegisterMarshalerFunc(
		instanceSDK.ServerTypeAvailability(""),
		human.EnumMarshalFunc(serverTypesAvailabilityMarshalSpecs),
	)
	human.RegisterMarshalerFunc([]*customServerType{}, serverTypesListMarshalerFunc)

	cmds.MustFind("instance", "server-type", "list").Override(serverTypeListBuilder)

	//
	// Get-Compatible-Types
	//
	cmds.MustFind("instance", "server", "get-compatible-types").Override(getCompatibleTypesBuilder)

	//
	// Templates
	//
	human.RegisterMarshalerFunc(instanceSDK.Template{}, templateMarshalerFunc)

	cmds.MustFind("instance", "template", "create").Override(TemplateCreateBuilder)
	cmds.MustFind("instance", "template", "get").Override(TemplateGetBuilder)
	cmds.MustFind("instance", "template", "list").Override(TemplateListBuilder)
	cmds.MustFind("instance", "template", "update").Override(TemplateUpdateBuilder)
	cmds.MustFind("instance", "template", "delete").Override(TemplateDeleteBuilder)
	cmds.MustFind("instance", "template", "check").Override(TemplateCheckBuilder)
	cmds.MustFind("instance", "template", "create-server").Override(TemplateCreateServerBuilder)

	cmds.MustFind("instance", "template", "set-user-data").Override(TemplateSetUserDataBuilder)
	cmds.MustFind("instance", "template", "get-user-data").Override(TemplateGetUserDataBuilder)
	cmds.MustFind("instance", "template", "list-user-data-keys").
		Override(TemplateListUserDataKeysBuilder)
	cmds.MustFind("instance", "template", "delete-user-data").
		Override(TemplateDeleteUserDataBuilder)

	cmds.MustFind("instance", "template", "set-cloud-init").Override(TemplateSetCloudInitBuilder)
	cmds.MustFind("instance", "template", "get-cloud-init").Override(TemplateGetCloudInitBuilder)

	//
	// User Data
	//
	cmds.MustFind("instance", "user-data", "set").Override(userDataSetBuilder)
	cmds.MustFind("instance", "user-data", "get").Override(userDataGetBuilder)
	cmds.MustFind("instance", "user-data", "list").Override(userDataListBuilder)
	cmds.MustFind("instance", "user-data", "delete").Override(userDataDeleteBuilder)

	cmds.MustFind("instance", "user-data", "set-cloud-init").Override(userDataSetCloudInitBuilder)
	cmds.MustFind("instance", "user-data", "get-cloud-init").Override(userDataGetCloudInitBuilder)

	// Web URLs (--web)

	addWebUrls(cmds)

	return cmds
}
