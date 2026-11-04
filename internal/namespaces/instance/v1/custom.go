package instance

import (
	"fmt"
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
		instanceSecurityGroup(),
		instanceServer(),
		instanceVolumeType(),
		instanceSnapshot(),
		instanceUserData(),
		instanceVolume(),
		instancePrivateNic(),
		instanceVolumeTypeList(),
		instanceServerList(),
		instanceServerGet(),
		instanceServerUpdate(),
		instanceServerListActions(),
		instanceUserDataList(),
		instanceUserDataDelete(),
		instanceUserDataSet(),
		instanceUserDataGet(),
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
		instanceSecurityGroupList(),
		instanceSecurityGroupCreate(),
		instanceSecurityGroupGet(),
		instanceSecurityGroupDelete(),
		instanceSecurityGroupUpdate(),
		instanceSecurityGroupListDefaultRules(),
		instanceSecurityGroupListRules(),
		instanceSecurityGroupCreateRule(),
		instanceSecurityGroupSetRules(),
		instanceSecurityGroupDeleteRule(),
		instanceSecurityGroupGetRule(),
		instanceSecurityGroupUpdateRule(),
		instanceIPList(),
		instanceIPCreate(),
		instanceIPGet(),
		instanceIPUpdate(),
		instanceIPDelete(),
		instancePrivateNicList(),
		instancePrivateNicCreate(),
		instancePrivateNicGet(),
		instancePrivateNicUpdate(),
		instancePrivateNicDelete(),
		instanceVolumePlanMigration(),
		instanceVolumeApplyMigration(),
		// PlacementGroup commands exclusive to v1
		instancePlacementGroupSet(),
		instancePlacementGroupGetServers(),
		instancePlacementGroupSetServers(),
		instancePlacementGroupUpdateServers(),
		// ServerType commands exclusive to v1
		instanceServerTypeGet(),
	)

	//
	// Server
	//
	human.RegisterMarshalerFunc(instance.CreateServerResponse{}, marshallNestedField("Server"))
	humanv2.RegisterEnumFormat(serverStateFormat)
	human.RegisterMarshalerFunc(instance.ServerLocation{}, serverLocationMarshalerFunc)
	human.RegisterMarshalerFunc([]*instance.Server{}, serversMarshalerFunc)
	human.RegisterMarshalerFunc(instance.Bootscript{}, bootscriptMarshalerFunc)

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

	//
	// IP
	//
	human.RegisterMarshalerFunc(instance.CreateIPResponse{}, marshallNestedField("IP"))

	cmds.MustFind("instance", "ip", "create").Override(ipCreateBuilder)
	cmds.MustFind("instance", "ip", "list").Override(ipListBuilder)
	cmds.Merge(core.NewCommands(
		ipAttachCommand(),
		ipDetachCommand(),
	))

	//
	// Image
	//
	human.RegisterMarshalerFunc(instance.CreateImageResponse{}, marshallNestedField("Image"))
	human.RegisterMarshalerFunc([]*imageListItem{}, imagesMarshalerFunc)
	humanv2.RegisterEnumFormat(imageStateFormat)
	cmds.MustFind("instance", "image", "create").Override(imageCreateBuilder)
	cmds.MustFind("instance", "image", "list").Override(imageListBuilder)
	cmds.MustFind("instance", "image", "delete").Override(imageDeleteBuilder)
	cmds.Merge(core.NewCommands(
		imageWaitCommand(),
	))

	//
	// Snapshot
	//
	human.RegisterMarshalerFunc(instance.CreateSnapshotResponse{}, marshallNestedField("Snapshot"))

	cmds.MustFind("instance", "snapshot", "create").Override(snapshotCreateBuilder)
	cmds.MustFind("instance", "snapshot", "list").Override(snapshotListBuilder)
	cmds.MustFind("instance", "snapshot", "update").Override(snapshotUpdateBuilder)
	cmds.Merge(core.NewCommands(
		snapshotWaitCommand(),
		snapshotPlanMigrationCommand(),
		snapshotApplyMigrationCommand(),
	))

	//
	// Volume
	//
	human.RegisterMarshalerFunc(instance.CreateVolumeResponse{}, marshallNestedField("Volume"))
	humanv2.RegisterEnumFormat(volumeStateFormat)
	human.RegisterMarshalerFunc(instance.VolumeSummary{}, volumeSummaryMarshalerFunc)
	human.RegisterMarshalerFunc(map[string]*instance.Volume{}, volumeMapMarshalerFunc)

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

	//
	// Security Group
	//
	human.RegisterMarshalerFunc(
		instance.CreateSecurityGroupResponse{},
		marshallNestedField("SecurityGroup"),
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

	//
	// Security Group Rule
	//
	human.RegisterMarshalerFunc(
		instance.CreateSecurityGroupRuleResponse{},
		marshallNestedField("Rule"),
	)
	humanv2.RegisterEnumFormat(securityGroupRuleActionFormat)
	human.RegisterMarshalerFunc([]*instance.SecurityGroupRule{}, marshalSecurityGroupRules)

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
	// User Data
	//
	cmds.MustFind("instance", "user-data", "delete").Override(userDataDeleteBuilder)
	cmds.MustFind("instance", "user-data", "set").Override(userDataSetBuilder)
	cmds.MustFind("instance", "user-data", "get").Override(userDataGetBuilder)
	cmds.MustFind("instance", "user-data", "list").Override(userDataListBuilder)

	//
	// Private NICs
	//
	humanv2.RegisterEnumFormat(privateNICStateFormat)
	cmds.MustFind("instance", "private-nic", "get").Override(privateNicGetBuilder)

	// SSH Utilities

	human.RegisterMarshalerFunc([]*SSHKeyFormat(nil), marshalSSHKeys)

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

	return cmds
}

// marshallNestedField will marshal only the given field of a struct.
func marshallNestedField(nestedKey string) human.MarshalerFunc {
	return func(i any, opt *human.MarshalOpt) (s string, err error) {
		if reflect.TypeOf(i).Kind() != reflect.Struct {
			return "", fmt.Errorf("%T must be a struct", i)
		}
		nestedValue := reflect.ValueOf(i).FieldByName(nestedKey)

		return human.Marshal(nestedValue.Interface(), opt)
	}
}
