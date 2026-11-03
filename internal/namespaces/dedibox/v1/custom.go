package dedibox

import (
	"context"

	"github.com/scaleway/scaleway-cli/v2/core"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("dedibox").Groups = []string{"baremetal"}

	for _, commandPath := range [][]string{
		{"dedibox", "server", "list"},
		{"dedibox", "service", "list"},
		{"dedibox", "offer", "list"},
		{"dedibox", "offer", "get"},
		{"dedibox", "os", "list"},
		{"dedibox", "fip", "list"},
		{"dedibox", "fip", "get-quota"},
		{"dedibox", "billing", "list-invoice"},
		{"dedibox", "billing", "list-refund"},
		{"dedibox", "ipv6-block", "get-quota"},
		{"dedibox", "ipv6-block", "create"},
		{"dedibox", "ipv6-block", "get"},
		{"dedibox", "rpn-info", "list"},
		{"dedibox", "rpn-info", "get"},
		{"dedibox", "san", "list"},
		{"dedibox", "rpn-v1", "list"},
		{"dedibox", "rpn-v1", "list-members"},
		{"dedibox", "rpn-v1", "list-capable-server"},
		{"dedibox", "rpn-v1", "list-capable-san-server"},
		{"dedibox", "rpn-v2", "list"},
		{"dedibox", "rpn-v2", "list-capable-resources"},
	} {
		cmds.MustFind(commandPath...).Override(setProjectDefaultValue)
	}

	cmds.MustFind("dedibox", "fip", "list").
		Override(core.ListViewBuilder("ID", "Address", "Reverse", "IPVersion", "Cidr", "Netmask", "GatewayIP", "Mac", "ServerID", "Status", "Type", "ServerZone"))
	cmds.MustFind("dedibox", "offer", "list").
		Override(core.ListViewBuilder("ID", "Name", "Catalog", "PaymentFrequency", "Pricing", "OfferID", "Commitments", "AvailableOs"))
	cmds.MustFind("dedibox", "os", "list").
		Override(core.ListViewBuilder("ID", "Name", "Type", "Version", "Arch", "AllowCustomPartitioning", "AllowSSHKeys", "RequiresUser", "RequiresAdminPassword", "RequiresPanelPassword", "AllowedFilesystems", "RequiresLicense", "LicenseOffers", "MaxPartitions", "DisplayName", "PasswordRegex", "PanelPasswordRegex", "RequiresValidHostname", "HostnameRegex", "HostnameMaxLength", "ReleasedAt", "Beta", "Description", "OfficialWebsite"))
	cmds.MustFind("dedibox", "rpn-info", "list").
		Override(core.ListViewBuilder("ID", "Hostname", "DatacenterName", "Zone", "IPAddress", "RpnVersion", "CompatibleQinq", "CanJoinQinqGroup", "Rpnv1GroupCount", "Rpnv2GroupCount", "CanJoinRpnv2Group"))
	cmds.MustFind("dedibox", "rpn-v1", "list").
		Override(core.ListViewBuilder("ID", "Name", "Type", "Active", "CreatedAt", "Owner", "MembersCount", "OrganizationID", "ProjectID"))
	cmds.MustFind("dedibox", "rpn-v2", "list").
		Override(core.ListViewBuilder("ID", "Name", "CompatibleRpnv1", "OrganizationID", "ProjectID", "Type", "Status", "Owner", "MembersCount", "Gateway"))
	cmds.MustFind("dedibox", "san", "list").
		Override(core.ListViewBuilder("ID", "DatacenterName", "OrganizationID", "ProjectID", "ServerHostname", "IqnSuffix", "CreatedAt", "OfferID", "OfferName", "Status", "StorageSize", "Rpnv1Compatible", "Rpnv1Implicit", "DeliveredAt", "TerminatedAt", "ExpiresAt"))
	cmds.MustFind("dedibox", "server", "list").
		Override(core.ListViewBuilder("ID", "DatacenterName", "OrganizationID", "ProjectID", "Hostname", "CreatedAt", "UpdatedAt", "ExpiredAt", "OfferID", "OfferName", "Status", "OsID", "Interfaces", "Zone", "IsOutsourced", "Qinq", "RpnVersion", "IsHds"))
	cmds.MustFind("dedibox", "service", "list").
		Override(core.ListViewBuilder("ID", "ResourceID", "ProvisioningStatus", "CreatedAt", "DeliveredAt", "TerminatedAt", "ExpiresAt", "Type"))

	return cmds
}

func setProjectDefaultValue(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("project-id").Default = func(ctx context.Context) (value string, doc string) {
		projectID := core.GetProjectIDFromContext(ctx)

		return projectID, "<retrieved from config>"
	}

	return c
}
