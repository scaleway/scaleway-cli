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

	return cmds
}
