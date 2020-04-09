package instance

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

var (
	volumeStateMarshalSpecs = human.EnumMarshalSpecs{
		instance.VolumeStateError:     &human.EnumMarshalSpec{Attribute: color.FgRed},
		instance.VolumeStateAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	}
)

// serversMarshalerFunc marshals a VolumeSummary.
func volumeSummaryMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumeSummary := i.(instance.VolumeSummary)
	return human.Marshal(volumeSummary.ID, opt)
}

// volumeMapMarshalerFunc returns the length of the map.
func volumeMapMarshalerFunc(i interface{}, opt *human.MarshalOpt) (string, error) {
	volumes := i.(map[string]*instance.Volume)
	return fmt.Sprintf("%v", len(volumes)), nil
}

// Builders

func volumeCreateBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}

func volumeListBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("organization").Name = "organization-id"
	return c
}
