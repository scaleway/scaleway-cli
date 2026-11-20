package instance

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

var (
	oldOrganizationFieldName = "organization"
	newOrganizationFieldName = "organization-id"
	oldProjectFieldName      = "project"
	newProjectFieldName      = "project-id"
)

// helpers
func renameOrganizationIDArgSpec(argSpecs core.ArgSpecs) {
	argSpecs.GetByName(oldOrganizationFieldName).Name = newOrganizationFieldName
}

func renameProjectIDArgSpec(argSpecs core.ArgSpecs) {
	argSpecs.GetByName(oldProjectFieldName).Name = newProjectFieldName
}

// GetCommands returns instance commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run and Command.View)
// - Merge handwritten commands
func GetCommands() *core.Commands {
	cmds := core.NewCommands(
		instanceRoot(),
		instanceImage(),
		instanceIP(),
		instanceServer(),
		instanceServerList(),
		instanceServerGet(),
		instanceServerUpdate(),
		instanceServerListActions(),
		instanceServerAttachVolume(),
		instanceServerDetachVolume(),
		instanceServerAttachFilesystem(),
		instanceServerDetachFilesystem(),
		instanceImageList(),
		instanceImageGet(),
		instanceImageCreate(),
		instanceImageUpdate(),
		instanceImageDelete(),
		instanceIPList(),
		instanceIPCreate(),
		instanceIPGet(),
		instanceIPUpdate(),
		instanceIPDelete(),
		// PlacementGroup commands exclusive to v1
		instancePlacementGroupSet(),
		instancePlacementGroupGetServers(),
		instancePlacementGroupSetServers(),
		instancePlacementGroupUpdateServers(),
		// PrivateNIC commands exclusive to v1
		instancePrivateNic(),
		instancePrivateNicList(),
		instancePrivateNicCreate(),
		instancePrivateNicGet(),
		instancePrivateNicUpdate(),
		instancePrivateNicDelete(),
		// SecurityGroup commands exclusive to v1
		instanceSecurityGroupListDefaultRules(),
		instanceSecurityGroupListRules(),
		instanceSecurityGroupCreateRule(),
		instanceSecurityGroupDeleteRule(),
		instanceSecurityGroupGetRule(),
		// ServerType commands exclusive to v1
		instanceServerTypeGet(),
	)

	// Server
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateServerResponse](),
		humanv2.Layout{
			Sections: []humanv2.Section{
				{Path: "Server", Main: true},
			},
		},
	)
	humanv2.RegisterEnumFormat(serverStateFormat)
	humanv2.RegisterFormatFunc(serverLocationFormatFunc)
	humanv2.RegisterLayout(reflect.TypeFor[[]serverListEntry](), serverListLayout)
	humanv2.RegisterFormatFunc(bootscriptFormatFunc)

	cmds.MustFind("instance").Groups = []string{"compute"}

	cmds.MustFind("instance", "server", "list").Override(serverListBuilder)
	cmds.MustFind("instance", "server", "update").Override(serverUpdateBuilder)
	cmds.MustFind("instance", "server", "get").Override(serverGetBuilder)

	cmds.Merge(core.NewCommands(
		serverBackupCommand(),
		serverCreateCommand(),
		serverDeleteCommand(),
		serverTerminateCommand(),
		serverSSHCommand(),
		serverActionCommand(),
		serverStartCommand(),
		serverStopCommand(),
		serverStandbyCommand(),
		serverRebootCommand(),
		serverEnableRoutedIPCommand(),
		serverWaitCommand(),
		serverAttachIPCommand(),
		serverDetachIPCommand(),
	))

	if cmdConsole := serverConsoleCommand(); cmdConsole != nil {
		cmds.Add(cmdConsole)
	}

	//
	// Server-Type
	//
	cmds.MustFind("instance", "server-type", "get").Override(serverTypeGetBuilder)

	// IP
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateIPResponse](),
		humanv2.Layout{
			Sections: []humanv2.Section{
				{Path: "IP", Main: true},
			},
		},
	)

	cmds.MustFind("instance", "ip", "create").Override(ipCreateBuilder)
	cmds.MustFind("instance", "ip", "list").Override(ipListBuilder)
	cmds.Merge(core.NewCommands(
		ipAttachCommand(),
		ipDetachCommand(),
	))

	// Image
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateImageResponse](),
		humanv2.Layout{
			Sections: []humanv2.Section{
				{Path: "Image", Main: true},
			},
		},
	)
	humanv2.RegisterLayout(reflect.TypeFor[[]*imageListItem](), imageListLayout)
	humanv2.RegisterEnumFormat(imageStateFormat)
	cmds.MustFind("instance", "image", "create").Override(imageCreateBuilder)
	cmds.MustFind("instance", "image", "list").Override(imageListBuilder)
	cmds.MustFind("instance", "image", "delete").Override(imageDeleteBuilder)
	cmds.Merge(core.NewCommands(
		imageWaitCommand(),
	))

	// Security Group Rule
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateSecurityGroupRuleResponse](),
		humanv2.Layout{
			Sections: []humanv2.Section{
				{Path: "Rule", Main: true},
			},
		},
	)
	humanv2.RegisterEnumFormat(securityGroupRuleActionFormat)

	cmds.MustFind("instance", "security-group", "list-default-rules").
		Override(securityGroupListDefaultRulesBuilder)
	cmds.MustFind("instance", "security-group", "list-rules").
		Override(securityGroupListRulesBuilder)
	cmds.MustFind("instance", "security-group", "create-rule").
		Override(securityGroupCreateRuleBuilder)
	cmds.MustFind("instance", "security-group", "delete-rule").
		Override(securityGroupDeleteRuleBuilder)
	cmds.MustFind("instance", "security-group", "get-rule").
		Override(securityGroupGetRuleBuilder)

	//
	// Placement Group
	//

	cmds.MustFind("instance", "placement-group", "set").Override(placementGroupSetBuilder)
	cmds.MustFind("instance", "placement-group", "set-servers").
		Override(placementGroupSetServersBuilder)
	cmds.MustFind("instance", "placement-group", "get-servers").
		Override(placementGroupGetServersBuilder)
	cmds.MustFind("instance", "placement-group", "update-servers").
		Override(placementGroupUpdateServersBuilder)

	//
	// Private NICs
	//

	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "create"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "list"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "update"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "delete"))

	humanv2.RegisterEnumFormat(privateNICStateFormat)
	cmds.MustFind("instance", "private-nic", "get").Override(privateNicGetBuilder)

	// SSH Utilities

	cmds.Merge(core.NewCommands(
		instanceSSH(),
		sshAddKeyCommand(),
		sshConfigInstallCommand(),
		sshFetchKeysCommand(),
		sshListKeysCommand(),
		sshRemoveKeyCommand(),
		instanceServerGetRdpPassword(),
	))

	// Web URLs (--web)

	addWebUrls(cmds)

	cmds.Merge(core.NewCommands(serverCockpitCommand()))

	return cmds
}
