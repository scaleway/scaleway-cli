package baremetal

import (
	"context"
	"errors"
	"net/http"

	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func serverDeleteBuilder(c *core.Command) *core.Command {
	c.WaitFunc = func(ctx context.Context, argsI, _ any) (any, error) {
		server, err := baremetal.NewAPI(core.ExtractClient(ctx)).
			WaitForServer(&baremetal.WaitForServerRequest{
				ServerID:      argsI.(*baremetal.DeleteServerRequest).ServerID,
				Zone:          argsI.(*baremetal.DeleteServerRequest).Zone,
				Timeout:       scw.TimeDurationPtr(ServerActionTimeout),
				RetryInterval: core.DefaultRetryInterval,
			})
		if err != nil {
			notFoundError := &scw.ResourceNotFoundError{}
			responseError := &scw.ResponseError{}
			if errors.As(err, &responseError) && responseError.StatusCode == http.StatusNotFound ||
				errors.As(err, &notFoundError) {
				return server, nil
			}

			return nil, err
		}

		return server, nil
	}

	return c
}
