package registry

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

//
// Marshalers
//

// namespaceStatusMarshalerFunc marshals a registry.NamespaceStatus.
var (
	namespaceStatusMarshalSpecs = human.EnumMarshalSpecs{
		registry.NamespaceStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		registry.NamespaceStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.NamespaceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		registry.NamespaceStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
	}
)
