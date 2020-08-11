package instance

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"

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

	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		req := argsI.(*instance.GetServerUserDataRequest)
		res, err := runner(ctx, argsI)
		if err != nil {
			if resErr, ok := err.(*scw.ResponseError); ok {
				if resErr.StatusCode == http.StatusNotFound {
					return nil, fmt.Errorf("'%s' key does not exists", req.Key)
				}
			}
			return nil, err
		}

		return res, nil
	})

	return c
}

func userDataListBuilder(c *core.Command) *core.Command {
	type userDataRow struct {
		Key   string `json:"key"`
		Value string `json:"value"`
	}
	c.AddInterceptors(func(ctx context.Context, argsI interface{}, runner core.CommandRunner) (interface{}, error) {
		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		args := argsI.(*instance.ListServerUserDataRequest)
		res, err := api.GetAllServerUserData(&instance.GetAllServerUserDataRequest{
			Zone:     args.Zone,
			ServerID: args.ServerID,
		})
		if err != nil {
			return nil, err
		}
		var r []userDataRow
		for a, v := range res.UserData {
			buf := new(strings.Builder)
			_, err := io.Copy(buf, v)
			if err != nil {
				return nil, err
			}
			r = append(r, userDataRow{
				Key:   a,
				Value: buf.String(),
			})
		}
		sort.Slice(r, func(i, j int) bool {
			return r[i].Key < r[j].Key
		})

		return r, nil
	})

	return c
}
