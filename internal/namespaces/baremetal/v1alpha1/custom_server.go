package baremetal

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const (
	serverActionTimeout = 10 * time.Minute
)

const (
	serverActionReboot = iota
)

func waitForServerFunc(action int) core.WaitFunc {
	return func(ctx context.Context, _, respI interface{}) (interface{}, error) {
		server, err := baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			ServerID: respI.(*baremetal.Server).ID,
			Zone:     respI.(*baremetal.Server).Zone,
			Timeout:  scw.DurationPtr(serverActionTimeout),
		})
		switch action {
		case serverActionReboot:
			return server, err
		}
		return nil, err
	}
}
