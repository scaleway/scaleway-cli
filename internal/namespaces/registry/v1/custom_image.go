package registry

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

//
// Marshalers
//

// imageStatusMarshalerFunc marshals a registry.ImageStatus.
var (
	imageStatusMarshalSpecs = human.EnumMarshalSpecs{
		registry.ImageStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		registry.ImageStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.ImageStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.ImageStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
