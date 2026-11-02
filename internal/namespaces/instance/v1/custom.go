package instance

import (
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
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
		instanceVolumeType(),
		instanceSnapshot(),
		instanceVolume(),
		instanceVolumeTypeList(),
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
		instanceSnapshotList(),
		instanceSnapshotCreate(),
		instanceSnapshotGet(),
		instanceSnapshotUpdate(),
		instanceSnapshotDelete(),
		instanceSnapshotExport(),
		instanceVolumeList(),
		instanceVolumeCreate(),
		instanceVolumeGet(),
		instanceVolumeUpdate(),
		instanceVolumeDelete(),
		instanceVolumePlanMigration(),
		instanceVolumeApplyMigration(),
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
		nestedResponseLayout("Server"),
	)
	humanv2.RegisterEnumFormat(serverStateFormat)
	humanv2.RegisterFormatFunc(func(v instance.ServerLocation) (string, error) {
		return serverLocationMarshalerFunc(v, nil)
	})
	humanv2.RegisterLayout(reflect.TypeFor[[]serverListEntry](), serverListLayout)
	humanv2.RegisterFormatFunc(func(v instance.Bootscript) (string, error) {
		return bootscriptMarshalerFunc(v, nil)
	})

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
	humanv2.RegisterLayout(reflect.TypeFor[instance.CreateIPResponse](), nestedResponseLayout("IP"))

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
		nestedResponseLayout("Image"),
	)
	humanv2.RegisterLayout(reflect.TypeFor[[]*imageListItem](), imageListLayout)
	humanv2.RegisterEnumFormat(imageStateFormat)
	cmds.MustFind("instance", "image", "create").Override(imageCreateBuilder)
	cmds.MustFind("instance", "image", "list").Override(imageListBuilder)
	cmds.MustFind("instance", "image", "delete").Override(imageDeleteBuilder)
	cmds.Merge(core.NewCommands(
		imageWaitCommand(),
	))

	// Snapshot
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateSnapshotResponse](),
		nestedResponseLayout("Snapshot"),
	)

	cmds.MustFind("instance", "snapshot", "create").Override(snapshotCreateBuilder)
	cmds.MustFind("instance", "snapshot", "list").Override(snapshotListBuilder)
	cmds.MustFind("instance", "snapshot", "update").Override(snapshotUpdateBuilder)
	cmds.Merge(core.NewCommands(
		snapshotWaitCommand(),
		snapshotPlanMigrationCommand(),
		snapshotApplyMigrationCommand(),
	))

	// Volume
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateVolumeResponse](),
		nestedResponseLayout("Volume"),
	)
	humanv2.RegisterEnumFormat(volumeStateFormat)
	humanv2.RegisterFormatFunc(func(v instance.VolumeSummary) (string, error) {
		return volumeSummaryMarshalerFunc(v, nil)
	})
	humanv2.RegisterFormatFunc(func(v map[string]*instance.Volume) (string, error) {
		return volumeMapMarshalerFunc(v, nil)
	})

	cmds.MustFind("instance", "volume", "create").Override(volumeCreateBuilder)
	cmds.MustFind("instance", "volume", "list").Override(volumeListBuilder)
	cmds.MustFind("instance", "volume", "plan-migration").Override(volumeMigrationBuilder)
	cmds.MustFind("instance", "volume", "apply-migration").Override(volumeMigrationBuilder)
	cmds.Merge(core.NewCommands(
		volumeWaitCommand(),
	))

	//
	// Volume-Type
	//
	cmds.MustFind("instance", "volume-type", "list").Override(volumeTypeListBuilder)

	// Security Group
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateSecurityGroupResponse](),
		nestedResponseLayout("SecurityGroup"),
	)
	humanv2.RegisterEnumFormat(securityGroupPolicyFormat)
	humanv2.RegisterEnumFormat(securityGroupStateFormat)
	cmds.MustFind("instance", "security-group", "create").Override(securityGroupCreateBuilder)
	cmds.MustFind("instance", "security-group", "get").Override(securityGroupGetBuilder)
	cmds.MustFind("instance", "security-group", "list").Override(securityGroupListBuilder)
	cmds.MustFind("instance", "security-group", "delete").Override(securityGroupDeleteBuilder)

	cmds.Merge(core.NewCommands(
		securityGroupClearCommand(),
		securityGroupEditCommand(),
	))

	// Security Group Rule
	//
	humanv2.RegisterLayout(
		reflect.TypeFor[instance.CreateSecurityGroupRuleResponse](),
		nestedResponseLayout("Rule"),
	)
	humanv2.RegisterEnumFormat(securityGroupRuleActionFormat)
	human.RegisterMarshalerFunc([]*instance.SecurityGroupRule{}, marshalSecurityGroupRules)

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
	human.RegisterMarshalerFunc(
		instance.PrivateNICState(""),
		human.EnumMarshalFunc(privateNICStateMarshalSpecs),
	)

	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "create"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "list"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "update"))
	deprecatePrivateNicCommand(cmds.MustFind("instance", "private-nic", "delete"))

	humanv2.RegisterEnumFormat(privateNICStateFormat)
	cmds.MustFind("instance", "private-nic", "get").Override(privateNicGetBuilder)

	// SSH Utilities

	// The v1 registration feeds the "Removed keys" success message, which is
	// built through human.Marshal in sshRemoveKeyCommand.
	human.RegisterMarshalerFunc([]*SSHKeyFormat(nil), marshalSSHKeys)
	humanv2.RegisterFormatFunc(func(v []*SSHKeyFormat) (string, error) {
		return marshalSSHKeys(v, nil)
	})

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

// nestedResponseLayout renders a Create*Response wrapper flat by resolving
// its single nested field as the payload: the v2 port of the v1
// marshallNestedField registrations below.
func nestedResponseLayout(path string) humanv2.Layout {
	return humanv2.Layout{Sections: []humanv2.Section{{Path: path, Main: true}}}
}
