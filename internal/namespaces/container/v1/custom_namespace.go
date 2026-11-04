package container

import (
	"context"
	"time"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	humanv2 "github.com/scaleway/scaleway-cli/v2/core/human/v2"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var (
	containerNamespaceActionTimeout = 5 * time.Minute

	namespaceStatusFormat = humanv2.EnumFormats[container.NamespaceStatus]{
		container.NamespaceStatusCreating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusDeleting:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusError:         humanv2.EnumFormat{Color: color.FgRed},
		container.NamespaceStatusLocked:        humanv2.EnumFormat{Color: color.FgRed},
		container.NamespaceStatusLocking:       humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusReady:         humanv2.EnumFormat{Color: color.FgGreen},
		container.NamespaceStatusUpdating:      humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusUpgrading:     humanv2.EnumFormat{Color: color.FgBlue},
		container.NamespaceStatusUnknownStatus: humanv2.EnumFormat{Color: color.Faint},
	}
)

func waitForContainerNamespace(ctx context.Context, _, respI any) (any, error) {
	ns := respI.(*container.Namespace)

	client := core.ExtractClient(ctx)
	api := container.NewAPI(client)

	return api.WaitForNamespace(&container.WaitForNamespaceRequest{
		NamespaceID:   ns.ID,
		Region:        ns.Region,
		Timeout:       new(containerNamespaceActionTimeout),
		RetryInterval: core.DefaultRetryInterval,
	})
}

func containerNamespaceCreateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForContainerNamespace

	return c
}

func containerNamespaceUpdateBuilder(c *core.Command) *core.Command {
	c.WaitFunc = waitForContainerNamespace

	return c
}

func containerNamespaceDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, _ any) (any, error) {
		req := argsI.(*container.DeleteNamespaceRequest)

		client := core.ExtractClient(ctx)
		api := container.NewAPI(client)
		_, err := api.WaitForNamespace(&container.WaitForNamespaceRequest{
			NamespaceID:   req.NamespaceID,
			Region:        req.Region,
			Timeout:       new(containerNamespaceActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
		if err != nil {
			if core.IsNotFoundError(err) {
				return nil, nil
			}

			return nil, err
		}

		return nil, nil
	}

	c.ArgSpecs.GetByName("namespace-id").Positional = true

	return c
}
