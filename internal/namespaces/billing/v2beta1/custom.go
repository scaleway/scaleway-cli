package billing

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-cli/v2/internal/color"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2beta1"
)

var invoiceTypeMarshalSpecs = human.EnumMarshalSpecs{
	billing.DownloadInvoiceRequestFileTypePdf: &human.EnumMarshalSpec{
		Attribute: color.FgHiBlue,
		Value:     "pdf",
	},
}

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("billing").Groups = []string{"cost"}

	human.RegisterMarshalerFunc(
		billing.DownloadInvoiceRequestFileType("pdf"),
		human.EnumMarshalFunc(invoiceTypeMarshalSpecs),
	)

	cmds.MustFind("billing", "invoice", "download").Override(invoiceDownloadBuilder)
	cmds.MustFind("billing", "invoice", "export").Override(invoiceExportBuilder)
	cmds.MustFind("billing", "charge", "list").Override(chargeListBuilder)

	return cmds
}
