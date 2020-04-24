package baremetal

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
)

func serverInstallBuilder(c *core.Command) *core.Command {

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))
		return api.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
			Zone:     argsI.(*baremetal.InstallServerRequest).Zone,
			ServerID: respI.(*baremetal.Server).ID,
			Timeout:  serverActionTimeout,
		})
	}

	return c
}
