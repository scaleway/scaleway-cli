package webhosting

import (
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	webhosting "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(webhosting.HostingStatus(""), human.EnumMarshalFunc(hostingStatusMarshalSpecs))
	human.RegisterMarshalerFunc(webhosting.DNSRecordsStatus(""), human.EnumMarshalFunc(hostingDNSMarshalSpecs))
	human.RegisterMarshalerFunc(webhosting.NameserverStatus(""), human.EnumMarshalFunc(nameserverMarshalSpecs))

	return cmds
}
