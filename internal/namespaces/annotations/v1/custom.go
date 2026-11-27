package annotations

import (
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	annotationsSDK "github.com/scaleway/scaleway-sdk-go/api/annotations/v1"
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	humanv2.RegisterFormatFunc(func(key annotationsSDK.BindingKey) (string, error) {
		return key.Name, nil
	})
	humanv2.RegisterFormatFunc(func(value annotationsSDK.BindingValue) (string, error) {
		return value.Name, nil
	})

	cmds.MustFind("annotations", "key-value", "list").Override(keyValueListBuilder)

	return cmds
}
