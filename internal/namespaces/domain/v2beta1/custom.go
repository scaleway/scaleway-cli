package domain

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

var domainTypes = []string{
	"A",
	"AAAA",
	"CNAME",
	"TXT",
	"SRV",
	"TLSA",
	"MX",
	"NS",
	"PTR",
	"CAA",
	"ALIAS",
	"LOC",
	"SSHFP",
	"HINFO",
	"RP",
	"URI",
	"DS",
	"NAPTR",
}

const defaultTTL = "3600"

// GetCommands returns dns commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("dns").Groups = []string{"Domain & WebHosting"}

	cmds.Merge(core.NewCommands(
		dnsRecordAddCommand(),
		dnsRecordSetCommand(),
		dnsRecordDeleteCommand(),
		dnsRecordImportCommand(),
	))

	cmds.MustFind("dns", "zone", "import").ArgSpecs.GetByName("bind-source.content").CanLoadFile = true

	human.RegisterMarshalerFunc(
		domain.DNSZoneStatus(""),
		human.EnumMarshalFunc(zoneStatusMarshalSpecs),
	)
	human.RegisterMarshalerFunc(
		domain.SSLCertificateStatus(""),
		human.EnumMarshalFunc(certificateStatusMarshalSpecs),
	)

	cmds.MustFind("domain", "contact", "list").
		Override(core.ListViewBuilder("Contact.ID", "Contact.Firstname", "Contact.Lastname", "Contact.CompanyName", "Contact.Email", "Contact.EmailStatus", "Contact.Status"))
	cmds.MustFind("domain", "domain", "list").
		Override(core.ListViewBuilder("Domain", "ProjectID", "OrganizationID", "AutoRenewStatus", "DnssecStatus", "EppCode", "ExpiredAt", "UpdatedAt", "CreatedAt", "Registrar", "IsExternal", "Status", "PendingTrade"))
	cmds.MustFind("domain", "task", "list").
		Override(core.ListViewBuilder("ID", "ProjectID", "OrganizationID", "Domain", "Type", "Status", "StartedAt", "UpdatedAt", "Message", "ContactIDentifier"))
	cmds.MustFind("domain", "tld", "list").
		Override(core.ListViewBuilder("Name", "DnssecSupport", "DurationInYearsMin", "DurationInYearsMax", "IDnSupport"))

	cmds.MustFind("dns", "certificate", "list").
		Override(core.ListViewBuilder("DNSZone", "AlternativeDNSZones", "Status", "PrivateKey", "CertificateChain", "CreatedAt", "ExpiredAt"))
	cmds.MustFind("domain", "host", "list").
		Override(core.ListViewBuilder("Domain", "Name", "IPs", "Status"))

	return cmds
}
