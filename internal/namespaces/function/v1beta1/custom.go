package function

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	function "github.com/scaleway/scaleway-sdk-go/api/function/v1beta1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.MustFind("function").Groups = []string{"serverless"}

	if cmdDeploy := functionDeploy(); cmdDeploy != nil {
		cmds.Add(cmdDeploy)
	}

	cmds.Add(functionCockpitCommand())

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[function.NamespaceStatus]{
		function.NamespaceStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
		function.NamespaceStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
		function.NamespaceStatusError:    humanv2.EnumFormat{Color: color.FgRed},
		function.NamespaceStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
		function.NamespaceStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
		function.NamespaceStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
		function.NamespaceStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[function.CronStatus]{
		function.CronStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
		function.CronStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
		function.CronStatusError:    humanv2.EnumFormat{Color: color.FgRed},
		function.CronStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
		function.CronStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
		function.CronStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
		function.CronStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[function.FunctionStatus]{
		function.FunctionStatusCreated:  humanv2.EnumFormat{Color: color.FgGreen},
		function.FunctionStatusCreating: humanv2.EnumFormat{Color: color.FgBlue},
		function.FunctionStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
		function.FunctionStatusError:    humanv2.EnumFormat{Color: color.FgRed},
		function.FunctionStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
		function.FunctionStatusPending:  humanv2.EnumFormat{Color: color.FgBlue},
		function.FunctionStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
		function.FunctionStatusUnknown:  humanv2.EnumFormat{Color: color.Faint},
	})

	return cmds
}
