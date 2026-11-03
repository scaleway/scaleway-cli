package webhosting

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/webhosting/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("webhosting").Groups = []string{"Domain & WebHosting"}

	human.RegisterMarshalerFunc(
		webhosting.HostingStatus(""),
		human.EnumMarshalFunc(hostingStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		webhosting.DNSRecordsStatus(""),
		human.EnumMarshalFunc(hostingDNSMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		webhosting.NameserverStatus(""),
		human.EnumMarshalFunc(nameserverMarshalSpecs),
	)

	cmds.MustFind("webhosting", "offer", "list").Override(webhostingOfferListBuilder)

	cmds.MustFind("webhosting", "control-panel", "list").Override(core.ListViewBuilder("Name", "Available", "LogoURL", "AvailableLanguages"))
	cmds.MustFind("webhosting", "hosting", "list").Override(core.ListViewBuilder("ID", "ProjectID", "CreatedAt", "UpdatedAt", "Status", "Domain", "Protected", "DNSStatus", "OfferName", "DomainStatus", "Region"))

	cmds.MustFind("webhosting", "backup", "list").Override(core.ListViewBuilder("ID", "Size", "CreatedAt", "Status", "TotalItems"))
	cmds.MustFind("webhosting", "database", "list").Override(core.ListViewBuilder("DatabaseName", "Users"))
	cmds.MustFind("webhosting", "database-user", "list").Override(core.ListViewBuilder("Username", "Databases"))
	cmds.MustFind("webhosting", "ftp-account", "list").Override(core.ListViewBuilder("Username", "Path"))
	cmds.MustFind("webhosting", "mail-account", "list").Override(core.ListViewBuilder("Domain", "Username"))
	cmds.MustFind("webhosting", "website", "list").Override(core.ListViewBuilder("Domain", "Path", "SslStatus"))

	return cmds
}
