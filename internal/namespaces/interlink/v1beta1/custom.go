package interlink

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	interlink "github.com/scaleway/scaleway-sdk-go/api/interlink/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("interlink").Groups = []string{"network"}

	human.RegisterMarshalerFunc(
		interlink.BgpStatus(""),
		human.EnumMarshalFunc(bgpStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		interlink.LinkStatus(""),
		human.EnumMarshalFunc(linkStatusMarshalSpecs),
	)

	cmds.MustFind("interlink", "connection", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Status", "Name", "Tags", "PopID", "BandwidthMbps", "AvailableLinkBandwidths", "CreatedAt", "UpdatedAt", "DemarcationInfo", "Region"))
	cmds.MustFind("interlink", "link", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "PopID", "BandwidthMbps", "Status", "BgpV4Status", "BgpV6Status", "VpcID", "RoutingPolicyID", "EnableRoutePropagation", "CreatedAt", "UpdatedAt", "Vlan", "RoutingPolicyV4ID", "RoutingPolicyV6ID", "Region"))
	cmds.MustFind("interlink", "partner", "list").Override(core.ListViewBuilder("ID", "Name", "ContactEmail", "LogoURL", "PortalURL", "CreatedAt", "UpdatedAt", "L3Connectivity", "Region"))
	cmds.MustFind("interlink", "pop", "list").Override(core.ListViewBuilder("ID", "Name", "HostingProviderName", "Address", "City", "LogoURL", "AvailableLinkBandwidthsMbps", "DisplayName", "Region"))
	cmds.MustFind("interlink", "routing-policy", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "PrefixFilterIn", "PrefixFilterOut", "CreatedAt", "UpdatedAt", "IsIPv6", "Region"))

	return cmds
}
