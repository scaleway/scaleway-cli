package instance

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
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
	cmds := GetGeneratedCommands()

	//
	// Server
	//
	human.RegisterMarshalerFunc(instance.CreateServerResponse{}, marshallNestedField("Server"))
	human.RegisterMarshalerFunc(
		instance.ServerState(""),
		human.EnumMarshalFunc(serverStateMarshalSpecs),
	)
	human.RegisterMarshalerFunc(instance.ServerLocation{}, serverLocationMarshalerFunc)
	human.RegisterMarshalerFunc([]*instance.Server{}, serversMarshalerFunc)
	human.RegisterMarshalerFunc(instance.Bootscript{}, bootscriptMarshalerFunc)

	cmds.MustFind("instance").Groups = []string{"compute"}

	cmds.MustFind("instance", "server", "list").Override(serverListBuilder)
	cmds.MustFind("instance", "server", "update").Override(serverUpdateBuilder)
	cmds.MustFind("instance", "server", "get").Override(serverGetBuilder)

	cmds.Merge(core.NewCommands(
		serverAttachVolumeCommand(),
		serverBackupCommand(),
		serverCreateCommand(),
		serverDeleteCommand(),
		serverTerminateCommand(),
		serverDetachVolumeCommand(),
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
	human.RegisterMarshalerFunc(
		instance.ServerTypesAvailability(""),
		human.EnumMarshalFunc(serverTypesAvailabilityMarshalSpecs),
	)

	cmds.MustFind("instance", "server-type", "list").Override(serverTypeListBuilder)

	//
	// Get-Compatible-Types
	//
	cmds.MustFind("instance", "server", "get-compatible-types").Override(getCompatibleTypesBuilder)

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
	human.RegisterMarshalerFunc(
		instance.ImageState(""),
		human.EnumMarshalFunc(imageStateMarshalSpecs),
	)

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
	human.RegisterMarshalerFunc(
		instance.VolumeState(""),
		human.EnumMarshalFunc(volumeStateMarshalSpecs),
	)
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
	human.RegisterMarshalerFunc(
		instance.SecurityGroupPolicy(""),
		human.EnumMarshalFunc(securityGroupPolicyMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		instance.SecurityGroupState(""),
		human.EnumMarshalFunc(securityGroupStateMarshalSpecs),
	)

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
	human.RegisterMarshalerFunc(
		instance.SecurityGroupRuleAction(""),
		human.EnumMarshalFunc(securityGroupRuleActionMarshalSpecs),
	)
	human.RegisterMarshalerFunc([]*instance.SecurityGroupRule{}, marshalSecurityGroupRules)

	//
	// Placement Group
	//
	human.RegisterMarshalerFunc(
		instance.CreatePlacementGroupResponse{},
		marshallNestedField("PlacementGroup"),
	)

	cmds.MustFind("instance", "placement-group", "create").Override(placementGroupCreateBuilder)
	cmds.MustFind("instance", "placement-group", "get").Override(placementGroupGetBuilder)
	cmds.MustFind("instance", "placement-group", "list").Override(placementGroupListBuilder)

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
	human.RegisterMarshalerFunc(
		instance.PrivateNICState(""),
		human.EnumMarshalFunc(privateNICStateMarshalSpecs),
	)

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
