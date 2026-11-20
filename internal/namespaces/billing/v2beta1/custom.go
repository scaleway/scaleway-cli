package billing

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
)

var invoiceTypeFormat = humanv2.EnumFormats[billing.DownloadInvoiceRequestFileType]{
	billing.DownloadInvoiceRequestFileTypePdf: humanv2.EnumFormat{
		Color: color.FgHiBlue,
		Value: "pdf",
	},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("billing").Groups = []string{"cost"}

	humanv2.RegisterEnumFormat(invoiceTypeFormat)
	cmds.MustFind("billing", "invoice", "download").Override(invoiceDownloadBuilder)
	cmds.MustFind("billing", "invoice", "export").Override(invoiceExportBuilder)
	cmds.MustFind("billing", "charge", "list").Override(chargeListBuilder)

	return cmds
}
