package instance

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
)

//
// Marshalers
//

var (
	volumeStateAttributes = human.Attributes{
		instance.VolumeStateError:     color.FgRed,
		instance.VolumeStateAvailable: color.FgGreen,
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
