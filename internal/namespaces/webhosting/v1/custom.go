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
	cmds.MustFind("webhosting", "database-user", "create").Override(databaseUserCreateBuilder)
	cmds.MustFind("webhosting", "database-user", "change-password").
		Override(databaseUserChangePasswordBuilder)
	cmds.MustFind("webhosting", "ftp-account", "create").Override(ftpAccountCreateBuilder)
	cmds.MustFind("webhosting", "mail-account", "create").Override(mailAccountCreateBuilder)
	cmds.MustFind("webhosting", "mail-account", "change-password").
		Override(mailAccountChangePasswordBuilder)

	return cmds
}
