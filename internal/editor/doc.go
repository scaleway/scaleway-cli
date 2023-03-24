package editor

import (
	"fmt"

	"github.com/scaleway/scaleway-cli/v2/internal/config"
	"github.com/scaleway/scaleway-cli/v2/internal/core"
)

var LongDescription = fmt.Sprintf(`This command starts your default editor to edit a marshaled version of your resource
Default editor will be taken from $VISUAL, then $EDITOR or will be %q`, config.GetSystemDefaultEditor())

func MarshalModeArgSpec() *core.ArgSpec {
	return &core.ArgSpec{
		Name:       "mode",
		Short:      "marshaling used when editing data",
		Required:   false,
		Default:    core.DefaultValueSetter(MarshalModeDefault),
		EnumValues: MarshalModeEnum,
	}
}
