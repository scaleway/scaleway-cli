package billing

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	billing "github.com/scaleway/scaleway-sdk-go/api/billing/v2alpha1"
)

var (
	invoiceTypeMarshalSpecs = human.EnumMarshalSpecs{
		billing.DownloadInvoiceRequestFileTypePdf: &human.EnumMarshalSpec{Attribute: color.FgHiBlue, Value: "pdf"},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(billing.DownloadInvoiceRequestFileType("pdf"), human.EnumMarshalFunc(invoiceTypeMarshalSpecs))
	cmds.MustFind("billing", "invoice", "download").Override(buildDownloadCommand)
	return cmds
}
