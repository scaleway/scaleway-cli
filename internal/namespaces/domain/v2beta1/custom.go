package domain

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
)

var (
	domainTypes = []string{"A", "AAAA", "CNAME", "TXT", "SRV", "TLSA", "MX", "NS", "PTR", "CAA", "ALIAS", "LOC", "SSHFP", "HINFO", "RP", "URI", "DS", "NAPTR"}
)

// GetCommands returns dns commands.
//
// This function:
// - Gets the generated commands
// - Register handwritten marshalers
// - Apply handwritten overrides (of Command.Run)
func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		dnsRecordAddCommand(),
		dnsRecordSetCommand(),
		dnsRecordDeleteCommand(),
	))

	human.RegisterMarshalerFunc(domain.DNSZoneStatus(""), human.EnumMarshalFunc(zoneStatusMarshalSpecs))
	human.RegisterMarshalerFunc(domain.SSLCertificateStatus(""), human.EnumMarshalFunc(certificateStatusMarshalSpecs))
	return cmds
}
