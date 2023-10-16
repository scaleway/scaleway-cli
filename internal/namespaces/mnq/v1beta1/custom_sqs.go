package mnq

import (
	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/internal/human"
	mnq "github.com/scaleway/scaleway-sdk-go/api/mnq/v1beta1"
)

var (
	mnqSqsInfoStatusMarshalSpecs = human.EnumMarshalSpecs{
		mnq.SqsInfoStatusUnknownStatus: &human.EnumMarshalSpec{Attribute: color.Faint},
		mnq.SqsInfoStatusEnabled:       &human.EnumMarshalSpec{Attribute: color.FgGreen},
		mnq.SqsInfoStatusDisabled:      &human.EnumMarshalSpec{Attribute: color.FgRed},
	}
)
