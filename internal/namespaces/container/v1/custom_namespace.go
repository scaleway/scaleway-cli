package container

import (
	"context"
	"time"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/container/v1"
)

var containerNamespaceActionTimeout = 5 * time.Minute

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
