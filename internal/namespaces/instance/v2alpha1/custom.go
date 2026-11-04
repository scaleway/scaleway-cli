package instance

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
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
	humanv2.RegisterEnumFormat(serverTypesAvailabilityFormat)
	human.RegisterMarshalerFunc([]*customServerType{}, serverTypesListMarshalerFunc)

	cmds.MustFind("instance", "server-type", "list").Override(serverTypeListBuilder)

	//
	// Get-Compatible-Types
	//
	cmds.MustFind("instance", "server", "get-compatible-types").Override(getCompatibleTypesBuilder)

	//
	// Templates
	//
	humanv2.RegisterLayout(reflect.TypeFor[instanceSDK.Template](), templateLayout)

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

	// Web URLs (--web)

	addWebUrls(cmds)

	return cmds
}
