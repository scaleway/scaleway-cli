package instance

import (
	"context"
	"errors"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	instance "github.com/scaleway/scaleway-sdk-go/api/instance/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

const serverActionTimeout = 10 * time.Minute

func serverStartBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = serverWaitFunc
	c.Examples = append(c.Examples, []*core.Example{
		{
			Short:    "Start a server in the default zone with a given id",
			ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
		{
			Short:    "Start a server in fr-par-1 zone with a given id",
			ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}...)

	return c
}

func serverStopBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = serverWaitFunc
	c.Examples = append(c.Examples, []*core.Example{
		{
			Short:    "Stop a server in the default zone with a given id",
			ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
		{
			Short:    "Stop a server in fr-par-1 zone with a given id",
			ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}...)

	return c
}

func serverRebootBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = serverWaitFunc
	c.Examples = append(c.Examples, []*core.Example{
		{
			Short:    "Reboot a server in the default zone with a given id",
			ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
		{
			Short:    "Reboot a server in fr-par-1 zone with a given id",
			ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}...)

	return c
}

func serverPauseBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = serverWaitFunc
	c.Examples = append(c.Examples, []*core.Example{
		{
			Short:    "Pause a server in the default zone with a given id",
			ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
		{
			Short:    "Pause a server in fr-par-1 zone with a given id",
			ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}...)

	return c
}

func serverStopAndDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	c.WaitFunc = func(ctx context.Context, argsI, _ any) (any, error) {
		req := argsI.(*instance.StopAndDeleteServerRequest)

		server, err := waitForServer(ctx, req.ServerID, req.Zone, nil, nil)
		if err != nil {
			if _, ok := errors.AsType[*scw.ResourceNotFoundError](err); !ok {
				return nil, err
			}
		}

		return server, nil
	}
	c.Examples = append(c.Examples, []*core.Example{
		{
			Short:    "Stop and delete a server in the default zone with a given id",
			ArgsJSON: `{"server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
		{
			Short:    "Stop and delete a server in fr-par-1 zone with a given id",
			ArgsJSON: `{"zone":"fr-par-1", "server_id": "11111111-1111-1111-1111-111111111111"}`,
		},
	}...)

	return c
}
