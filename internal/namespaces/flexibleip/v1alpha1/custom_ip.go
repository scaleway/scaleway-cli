package flexibleip

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	flexibleip "github.com/scaleway/scaleway-sdk-go/api/flexibleip/v1alpha1"
)

var ipStatusFormat = humanv2.EnumFormats[flexibleip.FlexibleIPStatus]{
	flexibleip.FlexibleIPStatusAttached:  humanv2.EnumFormat{Color: color.FgGreen},
	flexibleip.FlexibleIPStatusDetaching: humanv2.EnumFormat{Color: color.FgBlue},
	flexibleip.FlexibleIPStatusError:     humanv2.EnumFormat{Color: color.FgRed},
	flexibleip.FlexibleIPStatusLocked:    humanv2.EnumFormat{Color: color.FgRed},
	flexibleip.FlexibleIPStatusReady:     humanv2.EnumFormat{Color: color.FgGreen},
	flexibleip.FlexibleIPStatusUnknown:   humanv2.EnumFormat{Color: color.Faint},
	flexibleip.FlexibleIPStatusUpdating:  humanv2.EnumFormat{Color: color.FgBlue},
}

const (
	FlexibleIPTimeout = 60 * time.Second
)

func createIPBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, _, respI any) (any, error) {
		getResp := respI.(*flexibleip.FlexibleIP)
		api := flexibleip.NewAPI(core.ExtractClient(ctx))

		return api.WaitForFlexibleIP(&flexibleip.WaitForFlexibleIPRequest{
			FipID:         getResp.ID,
			Zone:          getResp.Zone,
			Timeout:       new(FlexibleIPTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
