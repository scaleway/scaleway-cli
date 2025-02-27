package webhosting

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	webhosting "github.com/scaleway/scaleway-sdk-go/api/webhosting/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(webhosting.HostingStatus(""), human.EnumMarshalFunc(hostingStatusMarshalSpecs))
	human.RegisterMarshalerFunc(webhosting.DNSRecordsStatus(""), human.EnumMarshalFunc(hostingDNSMarshalSpecs))
	human.RegisterMarshalerFunc(webhosting.NameserverStatus(""), human.EnumMarshalFunc(nameserverMarshalSpecs))

	cmds.MustFind("webhosting", "offer", "list").Override(webhostingOfferListBuilder)

	return cmds
}
