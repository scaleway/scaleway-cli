package baremetal

import (
	"context"
	"log"
	"reflect"

	"github.com/scaleway/scaleway-cli/internal/core"
	baremetal "github.com/scaleway/scaleway-sdk-go/api/baremetal/v1alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

type createServerRequest struct {
	Zone scw.Zone `json:"-"`
	// OfferID offer ID of the new server
	//OfferID string `json:"offer_id"`
	// OrganizationID with which the server will be created
	OrganizationID string `json:"organization_id"`
	// Name of the server (≠hostname)
	Name string `json:"name"`
	// Description associated to the server, max 255 characters
	Description string `json:"description"`
	// Tags associated with the server
	Tags []string `json:"tags"`
	// Type of the server
	Type string
}

func serverCreateCommand() *core.Command {
	return &core.Command{
		Short:     `Create server`,
		Long:      `Create a new server. Once the server is created, you probably want to install an OS.`,
		Namespace: "baremetal",
		Verb:      "create",
		Resource:  "server",
		ArgsType:  reflect.TypeOf(createServerRequest{}),
		ArgSpecs: core.ArgSpecs{
			core.ZoneArgSpec(scw.ZoneFrPar2),
			core.OrganizationIDArgSpec(),
			{
				Name:    "type",
				Short:   "Server commercial type",
				Default: core.DefaultValueSetter("GP-BM1-M"),

				EnumValues: []string{
					// General Purpose offers
					"GP-BM1-L",
					"GP-BM1-M",

					// High-computing offers
					"HC-BM1-L",
					"HC-BM1-S",

					// High-Memory offers
					"HM-BM1-XL",
					"HM-BM1-M",
				},
			},
			{
				Name:    "name",
				Short:   `Name of the server (≠hostname)`,
				Default: core.RandomValueGenerator("bm"),
			},
			{
				Name:  "description",
				Short: `Description associated to the server, max 255 characters`,
			},
			{
				Name:     "tags.{index}",
				Short:    `Tags to associate to the server`,
				Required: false,
			},
		},
		Run:      baremetalServerCreateRun,
		WaitFunc: baremetalWaitServerCreateRun(),
		SeeAlsos: []*core.SeeAlso{{
			Short:   "List os",
			Command: "scw baremetal os list",
		}},
		Examples: []*core.Example{
			{
				Short:   "Create instance",
				Request: `{}`,
			},
			{
				Short:   "Create a GP-BM1-M instance, give it a name and add tags",
				Request: `{"type":"GP-BM1-M","name":"foo","tags":["prod","blue"]}`,
			},
		},
	}
}

func baremetalWaitServerCreateRun() core.WaitFunc {
	return func(ctx context.Context, argsI, respI interface{}) (interface{}, error) {
		return baremetal.NewAPI(core.ExtractClient(ctx)).WaitForServer(&baremetal.WaitForServerRequest{
			Zone:     argsI.(*baremetal.CreateServerRequest).Zone,
			ServerID: respI.(*baremetal.Server).ID,
			Timeout:  serverActionTimeout,
		})
	}
}

func baremetalServerCreateRun(ctx context.Context, argsI interface{}) (i interface{}, e error) {
	client := core.ExtractClient(ctx)
	api := baremetal.NewAPI(client)

	tmpRequest := argsI.(*createServerRequest)
	request := &baremetal.CreateServerRequest{
		Zone:           tmpRequest.Zone,
		OrganizationID: tmpRequest.OrganizationID,
		Name:           tmpRequest.Name,
		Description:    tmpRequest.Description,
		Tags:           tmpRequest.Tags,
	}

	// We need to find the offer id.
	// while baremetal does not have listoffer name filter we are force to iterate
	// on the list of offer provided
	requestedType := tmpRequest.Type
	offerID := findOfferID(api, tmpRequest.Zone, requestedType)
	if offerID == "" {
		log.Fatal("Could not match")
	}
	request.OfferID = offerID

	return api.CreateServer(request)
}

func findOfferID(api *baremetal.API, zone scw.Zone, requestedType string) string {
	res, err := api.ListOffers(
		&baremetal.ListOffersRequest{
			Zone: zone},
		scw.WithAllPages())

	if err != nil {
		log.Fatal("Could not fetch list of offers.")
	}

	for _, v := range res.Offers {
		offerName := v.Name
		if requestedType == offerName {
			return v.ID
		}
	}
	return ""
}
