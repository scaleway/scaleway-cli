package flexibleip

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("fip").Groups = []string{"baremetal"}

	cmds.MustFind("fip", "ip", "create").Override(createIPBuilder)

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[flexibleip.FlexibleIPStatus]{
		flexibleip.FlexibleIPStatusAttached:  humanv2.EnumFormat{Color: color.FgGreen},
		flexibleip.FlexibleIPStatusDetaching: humanv2.EnumFormat{Color: color.FgBlue},
		flexibleip.FlexibleIPStatusError:     humanv2.EnumFormat{Color: color.FgRed},
		flexibleip.FlexibleIPStatusLocked:    humanv2.EnumFormat{Color: color.FgRed},
		flexibleip.FlexibleIPStatusReady:     humanv2.EnumFormat{Color: color.FgGreen},
		flexibleip.FlexibleIPStatusUnknown:   humanv2.EnumFormat{Color: color.Faint},
		flexibleip.FlexibleIPStatusUpdating:  humanv2.EnumFormat{Color: color.FgBlue},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[flexibleip.MACAddressStatus]{
		flexibleip.MACAddressStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
		flexibleip.MACAddressStatusReady:    humanv2.EnumFormat{Color: color.FgCyan},
		flexibleip.MACAddressStatusUpdating: humanv2.EnumFormat{Color: color.FgBlue},
		flexibleip.MACAddressStatusUsed:     humanv2.EnumFormat{Color: color.FgGreen},
		flexibleip.MACAddressStatusError:    humanv2.EnumFormat{Color: color.FgRed},
		flexibleip.MACAddressStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
	})

	return cmds
}
