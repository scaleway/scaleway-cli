package applesilicon

import (
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	applesilicon "github.com/scaleway/scaleway-sdk-go/api/applesilicon/v1alpha1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	cmds.Merge(core.NewCommands(
		serverSSHCommand(),
		serverWaitCommand(),
	))

	cmds.MustFind("apple-silicon").Groups = []string{"baremetal"}

	cmds.MustFind("apple-silicon", "server", "create").Override(serverCreateBuilder)
	cmds.MustFind("apple-silicon", "server", "reboot").Override(serverRebootBuilder)
	cmds.MustFind("apple-silicon", "server", "delete").Override(serverDeleteBuilder)

	cmds.MustFind("apple-silicon", "server-type", "list").Override(serverTypeBuilder)

	humanv2.RegisterLayout(reflect.TypeFor[applesilicon.OS](), osLayout)

	humanv2.RegisterFormatFunc(func(cpu applesilicon.ServerTypeCPU) (string, error) {
		return fmt.Sprintf("%s (%d cores)", cpu.Name, cpu.CoreCount), nil
	})
	humanv2.RegisterFormatFunc(func(disk applesilicon.ServerTypeDisk) (string, error) {
		return humanv2.Format(disk.Capacity)
	})
	humanv2.RegisterFormatFunc(func(memory applesilicon.ServerTypeMemory) (string, error) {
		return humanv2.Format(memory.Capacity)
	})

	humanv2.RegisterEnumFormat(serverStatusFormat)
	humanv2.RegisterEnumFormat(serverTypeStockFormat)

	return cmds
}
