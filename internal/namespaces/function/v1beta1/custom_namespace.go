package function

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var namespaceStatusFormat = humanv2.EnumFormats[function.NamespaceStatus]{
	function.NamespaceStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
	function.NamespaceStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
	function.NamespaceStatusError:    humanv2.EnumFormat{Color: color.FgRed},
	function.NamespaceStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
	function.NamespaceStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
	function.NamespaceStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
	function.NamespaceStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
}
