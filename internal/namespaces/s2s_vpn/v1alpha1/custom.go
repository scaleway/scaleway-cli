package s2s_vpn

import "github.com/scaleway/scaleway-cli/v2/core"

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("s2s-vpn").Groups = []string{"network"}

	cmds.MustFind("s2s-vpn", "connection", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "CreatedAt", "UpdatedAt", "Status", "IsIPv6", "InitiationPolicy", "SecretID", "SecretRevision", "Ikev2Ciphers", "EspCiphers", "RoutePropagationEnabled", "VpnGatewayID", "CustomerGatewayID", "TunnelStatus", "TunnelStatusIPv4", "TunnelStatusIPv6", "BgpStatusIPv4", "BgpStatusIPv6", "Region"))
	cmds.MustFind("s2s-vpn", "customer-gateway", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "CreatedAt", "UpdatedAt", "PublicIPv4", "PublicIPv6", "Asn", "ConnectionIDs", "Region"))
	cmds.MustFind("s2s-vpn", "routing-policy", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "CreatedAt", "UpdatedAt", "IsIPv6", "PrefixFilterIn", "PrefixFilterOut", "Region"))
	cmds.MustFind("s2s-vpn", "vpn-gateway", "list").Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Name", "Tags", "CreatedAt", "UpdatedAt", "Status", "GatewayType", "PrivateNetworkID", "IpamPrivateIPv4ID", "IpamPrivateIPv6ID", "Asn", "ConnectionIDs", "Zone", "Region"))
	cmds.MustFind("s2s-vpn", "vpn-gateway-type", "list").Override(core.ListViewBuilder("Name", "Bandwidth", "AllowedConnections", "Zones", "Region"))

	return cmds
}
