package instance

import (
	"context"
	"fmt"
	"net/http"

	"github.com/scaleway/scaleway-cli/internal/core"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Commands
//

func userDataDeleteBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true
	return c
}

func userDataSetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	*c.ArgSpecs.GetByName("content.name") = core.ArgSpec{
		Name:     "content",
		Short:    "Content of the user data",
		Required: true,
	}

	c.ArgSpecs.DeleteByName("content.content-type")
	c.ArgSpecs.DeleteByName("content.content")
	return c
}

func userDataGetBuilder(c *core.Command) *core.Command {
	c.ArgSpecs.GetByName("server-id").Positional = true

	originalRun := c.Run
	c.Run = func(ctx context.Context, argsI interface{}) (interface{}, error) {
		req := argsI.(*instance.GetServerUserDataRequest)
		res, err := originalRun(ctx, argsI)
		if err != nil {
			if resErr, ok := err.(*scw.ResponseError); ok {
				if resErr.StatusCode == http.StatusNotFound {
					return nil, fmt.Errorf("'%s' key does not exists", req.Key)
				}
			}
			return nil, err
		}

		return res, nil
	}

	return c
}
