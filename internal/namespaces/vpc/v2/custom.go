package vpc

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/vpc/v2"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("vpc").Groups = []string{"network"}

	cmds.Remove("vpc", "post")
	cmds.MustFind("vpc", "private-network", "get").Override(privateNetworkGetBuilder)
	cmds.MustFind("vpc", "private-network", "delete").Override(privateNetworkDeleteBuilder)
	human.RegisterMarshalerFunc(vpc.PrivateNetwork{}, privateNetworkMarshalerFunc)

	cmds.Merge(core.NewCommands(
		vpcACLEditCommand(),
	))

	cmds.MustFind("vpc", "ingress-rule", "list").
		Override(core.ListViewBuilder("ID", "VpcID", "CreatedAt", "UpdatedAt", "IsIPv6", "Source", "NexthopResourceIP", "NexthopPrivateNetworkID", "Description", "Tags", "OrganizationID", "ProjectID", "Region"))
	cmds.MustFind("vpc", "private-network", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Region", "Tags", "CreatedAt", "UpdatedAt", "Subnets", "VpcID", "DHCPEnabled", "DefaultRoutePropagationEnabled", "HasObjectStoragePrivateAccess"))
	cmds.MustFind("vpc", "route", "list").
		Override(core.ListViewBuilder("NexthopIP", "NexthopName", "NexthopResourceType"))
	cmds.MustFind("vpc", "vpc", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "Region", "Tags", "IsDefault", "CreatedAt", "UpdatedAt", "PrivateNetworkCount", "RoutingEnabled", "CustomRoutesPropagationEnabled", "TransitivityEnabled", "ObjectStoragePrivateAccessEnabled"))
	cmds.MustFind("vpc", "vpc-connector", "list").
		Override(core.ListViewBuilder("ID", "Name", "OrganizationID", "ProjectID", "VpcID", "TargetVpcID", "Status", "Region", "Tags", "CreatedAt", "UpdatedAt"))

	return cmds
}
