package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

var (
	containerStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.ContainerStatusCreated:  &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.ContainerStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.ContainerStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.ContainerStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.ContainerStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
