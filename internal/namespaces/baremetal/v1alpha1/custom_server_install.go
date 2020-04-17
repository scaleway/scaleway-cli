package baremetal

import (
	"context"

	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
)

func serverInstallBuilder(c *core.Command) *core.Command {
	c.Examples = []*core.Example{
		{
			Short:   "Install an os on a given server",
			Request: `{"os-id":"11111111-1111-1111-1111-111111111111","server-id":"11111111-1111-1111-1111-111111111111"}`,
		},
	}

	c.SeeAlsos = []*core.SeeAlso{
		{
			Short:   "List os (Useful to get all os-id)",
			Command: "scw baremetal os list",
		},
		{
			Short:   "Create a server",
			Command: "scw baremetal server create",
		},
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))
		return api.WaitForServerInstall(&baremetal.WaitForServerInstallRequest{
			Zone:     argsI.(*baremetal.Server).Zone,
			ServerID: respI.(*baremetal.Server).ID,
			Timeout:  serverActionTimeout,
		})
	}

	return c
}
