package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

var (
	namespaceStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.NamespaceStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.NamespaceStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.NamespaceStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.NamespaceStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.NamespaceStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
