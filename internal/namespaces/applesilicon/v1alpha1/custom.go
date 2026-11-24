package applesilicon

import (
	"fmt"
	"reflect"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human/v2"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		serverSSHCommand(),
		serverWaitCommand(),
	))

	cmds.MustFind("apple-silicon").Groups = []string{"baremetal"}

	humanv2.RegisterLayout(reflect.TypeFor[*applesilicon.OS](), humanv2.Layout{
		Sections: []humanv2.Section{
			{
				Path:  "CompatibleServerTypes",
				Title: "Compatible Server Types",
			},
		},
	})
	humanv2.RegisterFormatFunc(func(cpu applesilicon.ServerTypeCPU) (string, error) {
		return fmt.Sprintf("%s (%d cores)", cpu.Name, cpu.CoreCount), nil
	})
	humanv2.RegisterFormatFunc(func(disk applesilicon.ServerTypeDisk) (string, error) {
		return humanv2.Format(disk.Capacity)
	})
	humanv2.RegisterFormatFunc(func(memory applesilicon.ServerTypeMemory) (string, error) {
		return humanv2.Format(memory.Capacity)
	})

	humanv2.RegisterEnumFormat(humanv2.EnumFormats[applesilicon.ServerTypeStock]{
		applesilicon.ServerTypeStockLowStock:  humanv2.EnumFormat{Color: color.FgYellow, Value: "low stock"},
		applesilicon.ServerTypeStockNoStock:   humanv2.EnumFormat{Color: color.FgRed, Value: "no stock"},
		applesilicon.ServerTypeStockHighStock: humanv2.EnumFormat{Color: color.FgGreen, Value: "high stock"},
	})
	humanv2.RegisterEnumFormat(humanv2.EnumFormats[applesilicon.ServerStatus]{
		applesilicon.ServerStatusError:     humanv2.EnumFormat{Color: color.FgRed, Value: "error"},
		applesilicon.ServerStatusReady:     humanv2.EnumFormat{Color: color.FgGreen, Value: "ready"},
		applesilicon.ServerStatusRebooting: humanv2.EnumFormat{Color: color.FgBlue, Value: "rebooting"},
		applesilicon.ServerStatusStarting:  humanv2.EnumFormat{Color: color.FgBlue, Value: "starting"},
		applesilicon.ServerStatusUpdating:  humanv2.EnumFormat{Color: color.FgBlue, Value: "updating"},
	})

	cmds.MustFind("apple-silicon", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("apple-silicon", "server", "reboot").Override(serverRebootBuilder)
	cmds.MustFind("apple-silicon", "server", "delete").Override(serverDeleteBuilder)

	cmds.MustFind("apple-silicon", "server-type", "list").Override(serverTypeBuilder)

	return cmds
}
