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
		instanceTemplate(),
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

	return cmds
}
