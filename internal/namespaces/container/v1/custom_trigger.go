package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var triggerStatusMarshalSpecs = human.EnumMarshalSpecs{
	container.TriggerStatusCreating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.TriggerStatusDeleting:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.TriggerStatusError:         &human.EnumMarshalSpec{Attribute: color.FgRed},
	container.TriggerStatusLocked:        &human.EnumMarshalSpec{Attribute: color.FgRed},
	container.TriggerStatusLocking:       &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.TriggerStatusReady:         &human.EnumMarshalSpec{Attribute: color.FgGreen},
	container.TriggerStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
	container.TriggerStatusUpdating:      &human.EnumMarshalSpec{Attribute: color.FgBlue},
	container.TriggerStatusUpgrading:     &human.EnumMarshalSpec{Attribute: color.FgBlue},
}
