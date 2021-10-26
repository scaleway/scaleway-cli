package container

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/internal/human"
	container "github.com/scaleway/scaleway-sdk-go/api/container/v1beta1"
)

var (
	cronStatusMarshalSpecs = human.EnumMarshalSpecs{
		container.CronStatusCreating: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.CronStatusDeleting: &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.CronStatusError:    &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.CronStatusLocked:   &human.EnumMarshalSpec{Attribute: color.FgRed},
		container.CronStatusPending:  &human.EnumMarshalSpec{Attribute: color.FgBlue},
		container.CronStatusReady:    &human.EnumMarshalSpec{Attribute: color.FgGreen},
		container.CronStatusUnknown:  &human.EnumMarshalSpec{Attribute: color.Faint},
	}
)
