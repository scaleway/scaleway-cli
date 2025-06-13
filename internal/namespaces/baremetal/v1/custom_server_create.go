package baremetal

import (
	"context"
	"fmt"
	"reflect"

	"github.com/scaleway/scaleway-cli/v2/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

func serverCreateBuilder(c *core.Command) *core.Command {
	type baremetalCreateServerRequestCustom struct {
		Zone scw.Zone `json:"-"`
		// OrganizationID with which the server will be created
		OrganizationID *string `json:"organization_id"`
		// ProjectID with which the server will be created
		ProjectID *string `json:"project_id"`
		// Name of the server (≠hostname)
		Name string `json:"name"`
		// Description associated to the server, max 255 characters
		Description string `json:"description"`
		// Tags associated with the server
		Tags []string `json:"tags"`
		// Type of the server
		Type string
		// Installation configuration
		Install *baremetal.CreateServerRequestInstall
	}

	c.ArgsType = reflect.TypeOf(baremetalCreateServerRequestCustom{})

	c.ArgSpecs.DeleteByName("offer-id")

	c.ArgSpecs.GetByName("name").Default = core.RandomValueGenerator("bm")
	c.ArgSpecs.GetByName("description").Required = false

	c.ArgSpecs.AddBefore("tags.{index}", &core.ArgSpec{
		Name:  "type",
		Short: "Server commercial type",
	})

	c.Run = func(ctx context.Context, argsI any) (i any, e error) {
		client := core.ExtractClient(ctx)
		api := baremetal.NewAPI(client)

		tmpRequest := argsI.(*baremetalCreateServerRequestCustom)
		request := &baremetal.CreateServerRequest{
			Zone:           tmpRequest.Zone,
			OrganizationID: tmpRequest.OrganizationID,
			ProjectID:      tmpRequest.ProjectID,
			Name:           tmpRequest.Name,
			Description:    tmpRequest.Description,
			Tags:           tmpRequest.Tags,
		}

		if tmpRequest.Install != nil {
			request.Install = tmpRequest.Install
		}

		// We need to find the offer ID.
		// While baremetal does not have list offer name filter we are forced to iterate
		// on the list of offers provided.
		offer, err := api.GetOfferByName(&baremetal.GetOfferByNameRequest{
			OfferName: tmpRequest.Type,
			Zone:      tmpRequest.Zone,
		})
		if err != nil {
			return nil, err
		}
		if offer == nil {
			return nil, fmt.Errorf("could not match an offer with the type: %s", tmpRequest.Type)
		}
		request.OfferID = offer.ID

		return api.CreateServer(request)
	}

	c.SeeAlsos = []*core.SeeAlso{
		{
			Short:   "List os",
			Command: "scw baremetal os list",
		},
		{
			Short:   "Install an OS on your server",
			Command: "scw baremetal server install",
		},
	}

	c.Examples = []*core.Example{
		{
			Short:    "Create instance",
			ArgsJSON: `{}`,
		},
	}

	c.WaitFunc = func(ctx context.Context, argsI, respI any) (any, error) {
		api := baremetal.NewAPI(core.ExtractClient(ctx))

		return api.WaitForServer(&baremetal.WaitForServerRequest{
			Zone:          argsI.(*baremetalCreateServerRequestCustom).Zone,
			ServerID:      respI.(*baremetal.Server).ID,
			Timeout:       scw.TimeDurationPtr(ServerActionTimeout),
			RetryInterval: core.DefaultRetryInterval,
		})
	}

	return c
}
