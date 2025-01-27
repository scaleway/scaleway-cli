package editor

import (
	"github.com/scaleway/scaleway-cli/v2/core"
)

var LongDescription = `This command starts your default editor to edit a marshaled version of your resource
Default editor will be taken from $VISUAL, then $EDITOR or an editor based on your system`

func MarshalModeArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "mode",
		Short:      "marshaling used when editing data",
		Required:   false,
		Default:    core.DefaultValueSetter(MarshalModeDefault),
		EnumValues: MarshalModeEnum,
	}
}
