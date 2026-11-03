package tem

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	tem "github.com/scaleway/scaleway-sdk-go/api/tem/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("tem").Groups = []string{"Domain & WebHosting"}

	human.RegisterMarshalerFunc(
		tem.DomainStatus(""),
		human.EnumMarshalFunc(domainStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		tem.EmailStatus(""),
		human.EnumMarshalFunc(emailStatusMarshalSpecs))

	cmds.MustFind("tem", "domain", "get").Override(domainGetBuilder)

	cmds.MustFind("tem", "blocklists", "list").
		Override(core.ListViewBuilder("ID", "DomainID", "CreatedAt", "UpdatedAt", "EndsAt", "Email", "Type", "Reason", "Custom"))
	cmds.MustFind("tem", "domain", "list").
		Override(core.ListViewBuilder("ID", "OrganizationID", "ProjectID", "Name", "Status", "CreatedAt", "NextCheckAt", "LastValidAt", "RevokedAt", "LastError", "SpfConfig", "DkimConfig", "Autoconfig", "Region"))
	cmds.MustFind("tem", "email", "list").
		Override(core.ListViewBuilder("ID", "MessageID", "ProjectID", "MailFrom", "RcptTo", "MailRcpt", "RcptType", "Subject", "CreatedAt", "UpdatedAt", "Status", "StatusDetails", "TryCount", "LastTries", "Flags"))
	cmds.MustFind("tem", "offer-subscriptions", "list").
		Override(core.ListViewBuilder("OrganizationID", "ProjectID", "OfferName", "SubscribedAt", "CancellationAvailableAt", "SLA", "MaxDomains", "MaxDedicatedIPs", "MaxWebhooksPerDomain", "MaxCustomBlocklistsPerDomain", "IncludedMonthlyEmails"))
	cmds.MustFind("tem", "pools", "list").
		Override(core.ListViewBuilder("ProjectID", "Status", "Details", "Zone", "IPs", "Reverse"))
	cmds.MustFind("tem", "webhook", "list").
		Override(core.ListViewBuilder("ID", "DomainID", "OrganizationID", "ProjectID", "Name", "EventTypes", "SnsArn", "CreatedAt", "UpdatedAt"))

	return cmds
}
