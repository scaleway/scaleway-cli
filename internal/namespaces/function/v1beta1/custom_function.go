package function

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var functionStatusFormat = humanv2.EnumFormats[function.FunctionStatus]{
	function.FunctionStatusCreated:  humanv2.EnumFormat{Color: color.FgGreen},
	function.FunctionStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
	function.FunctionStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
	function.FunctionStatusError:    humanv2.EnumFormat{Color: color.FgRed},
	function.FunctionStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
	function.FunctionStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
	function.FunctionStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
	function.FunctionStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
}
