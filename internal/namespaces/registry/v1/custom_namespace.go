package registry

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/registry/v1"
)

//
// Marshalers
//

// namespaceStatusMarshalerFunc marshals a registry.NamespaceStatus.
var (
	namespaceStatusFormat = humanv2.EnumFormats[registry.NamespaceStatus]{
		registry.NamespaceStatusReady:    humanv2.EnumFormat{Color: color.FgGreen},
		registry.NamespaceStatusLocked:   humanv2.EnumFormat{Color: color.FgRed},
		registry.NamespaceStatusError:    humanv2.EnumFormat{Color: color.FgRed},
		registry.NamespaceStatusDeleting: humanv2.EnumFormat{Color: color.FgBlue},
	}
)
