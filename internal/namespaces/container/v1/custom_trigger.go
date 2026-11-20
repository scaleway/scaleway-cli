package container

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var triggerStatusFormat = humanv2.EnumFormats[container.TriggerStatus]{
	container.TriggerStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
	container.TriggerStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
	container.TriggerStatusError:         humanv2.EnumFormat{Color: color.FgRed},
	container.TriggerStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	container.TriggerStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
	container.TriggerStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
	container.TriggerStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
	container.TriggerStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
	container.TriggerStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
}
