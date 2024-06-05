package key_manager

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	key_manager "github.com/scaleway/scaleway-sdk-go/api/key_manager/v1alpha1"
)

var (
	keymanagerKeyStateMarshalSpecs = human.EnumMarshalSpecs{
		key_manager.KeyStateUnknownState:       &human.EnumMarshalSpec{Attribute: color.Faint},
		key_manager.KeyStateEnabled:            &human.EnumMarshalSpec{Attribute: color.FgGreen},
		key_manager.KeyStateDisabled:           &human.EnumMarshalSpec{Attribute: color.FgRed},
		key_manager.KeyStatePendingKeyMaterial: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)

func GetCommands() *core.Commands {
	cmds := GetGeneratedCommands()

	human.RegisterMarshalerFunc(key_manager.KeyState(""), human.EnumMarshalFunc(keymanagerKeyStateMarshalSpecs))

	return cmds
}
