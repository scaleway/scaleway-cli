package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var domainStatusMarshalSpecs = human.EnumMarshalSpecs{
	container.DomainStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.DomainStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.DomainStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
	container.DomainStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	container.DomainStatusLocking:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.DomainStatusReady:         &human.EnumMarshalSpec{Attribute: color.FgGreen},
	container.DomainStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
	container.DomainStatusUpdating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.DomainStatusUpgrading:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
}
