package domain

import (
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	domain "github.com/scaleway/scaleway-sdk-go/api/domain/v2beta1"
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
	return cmds
}
