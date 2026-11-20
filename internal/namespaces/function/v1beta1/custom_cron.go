package function

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

var cronStatusFormat = humanv2.EnumFormats[function.CronStatus]{
	function.CronStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
	function.CronStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
	function.CronStatusError:    humanv2.EnumFormat{Color: color.FgRed},
	function.CronStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
	function.CronStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
	function.CronStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
	function.CronStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
}
