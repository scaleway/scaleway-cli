package registry

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

//
// Marshalers
//

// tagStatusMarshalerFunc marshals a registry.TagStatus.
var (
	tagStatusMarshalSpecs = human.EnumMarshalSpecs{
		registry.TagStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		registry.TagStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.TagStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.TagStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
