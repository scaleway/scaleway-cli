package container

import (
	"github.com/fatih/color"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var domainStatusFormat = humanv2.EnumFormats[container.DomainStatus]{
	container.DomainStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
	container.DomainStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
	container.DomainStatusError:         humanv2.EnumFormat{Color: color.FgRed},
	container.DomainStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
	container.DomainStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
	container.DomainStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
	container.DomainStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
	container.DomainStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
	container.DomainStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
}
