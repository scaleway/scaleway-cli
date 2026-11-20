package lb

import (
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
)

// lbPrivateNetworkListLayout reproduces v1's projected list view, dropping
// the raw DHCP/Static config structs in favour of their flat ID columns.
var lbPrivateNetworkListLayout = humanv2.Layout{
	Sections: []humanv2.Section{
		{
			Path: "",
			Fields: []humanv2.Field{
				{Path: "IpamIDs"},
				{Path: "DHCPConfig.IPID"},
				{Path: "StaticConfig.IPAddress"},
				{Path: "PrivateNetworkID"},
				{Path: "Status"},
				{Path: "CreatedAt"},
				{Path: "UpdatedAt"},
			},
		},
	},
}
