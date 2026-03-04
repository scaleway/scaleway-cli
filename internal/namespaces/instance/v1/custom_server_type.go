package instance

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/scaleway/scaleway-cli/v2/core"
	"github.com/scaleway/scaleway-cli/v2/core/human"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	product_catalog "github.com/scaleway/scaleway-sdk-go/api/product_catalog/v2alpha1"
	"github.com/scaleway/scaleway-sdk-go/scw"
)

//
// Marshalers
//

var serverTypesAvailabilityMarshalSpecs = human.EnumMarshalSpecs{
	instance.ServerTypesAvailabilityAvailable: &human.EnumMarshalSpec{Attribute: color.FgGreen},
	instance.ServerTypesAvailabilityScarce: &human.EnumMarshalSpec{
		Attribute: color.FgYellow,
		Value:     "low stock",
	},
	instance.ServerTypesAvailabilityShortage: &human.EnumMarshalSpec{
		Attribute: color.FgRed,
		Value:     "out of stock",
	},
}

//
// Builders
//

type customServerType struct {
	Name             string     `json:"name"`
	HourlyPrice      *scw.Money `json:"hourly_price"`
	SupportedStorage string
	CPU              uint32                           `json:"cpu"`
	GPU              uint32                           `json:"gpu"`
	RAM              scw.Size                         `json:"ram"`
	Arch             string                           `json:"arch"`
	Bandwidth        uint64                           `json:"bandwidth"`
	Availability     instance.ServerTypesAvailability `json:"availability"`
	MaxFileSystems   uint32                           `json:"max_file_systems"`
}

// serverTypeListBuilder transforms the server map into a list to display a
// table of server types instead of a flat key/value list.
// We need it for:
// - [APIGW-1932] hide deprecated instance for scw instance server-type list
func serverTypeListBuilder(c *core.Command) *core.Command {
	c.Run = func(ctx context.Context, argsI any) (any, error) {
		pcuAPI := product_catalog.NewPublicCatalogAPI(core.ExtractClient(ctx))
		instanceAPI := instance.NewAPI(core.ExtractClient(ctx))

		// Get server types from Product Catalog API
		request := argsI.(*instance.ListServersTypesRequest)
		instanceProductType := product_catalog.ListPublicCatalogProductsRequestProductTypeInstance
		listServersTypesResponse, err := pcuAPI.ListPublicCatalogProducts(
			&product_catalog.PublicCatalogAPIListPublicCatalogProductsRequest{
				ProductTypes: []product_catalog.ListPublicCatalogProductsRequestProductType{
					instanceProductType,
				},
				Zone:   &request.Zone,
				Status: nil,
			},
			scw.WithAllPages(),
			scw.WithContext(ctx),
		)
		if err != nil {
			return nil, err
		}

		// Get server types from Instance API (still needed for the number of file systems)
		computeServerTypes, err := instanceAPI.ListServersTypes(request, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		// Get server availabilities.
		availabilitiesResponse, err := instanceAPI.GetServerTypesAvailability(
			&instance.GetServerTypesAvailabilityRequest{
				Zone: request.Zone,
			},
			scw.WithAllPages(),
		)
		if err != nil {
			return nil, err
		}

		serverTypes := []*customServerType(nil)

		for _, pcuServerType := range listServersTypesResponse.Products {
			switch pcuServerType.Status {
			case product_catalog.PublicCatalogProductStatusUnknownStatus:
				continue
			case product_catalog.PublicCatalogProductStatusPublicBeta:
			case product_catalog.PublicCatalogProductStatusPreview:
			case product_catalog.PublicCatalogProductStatusGeneralAvailability:
			case product_catalog.PublicCatalogProductStatusEndOfNewFeatures:
			case product_catalog.PublicCatalogProductStatusEndOfGrowth:
				continue
			case product_catalog.PublicCatalogProductStatusEndOfDeployment:
				continue
			case product_catalog.PublicCatalogProductStatusEndOfSupport:
				continue
			case product_catalog.PublicCatalogProductStatusEndOfSale:
				continue
			case product_catalog.PublicCatalogProductStatusEndOfLife:
				continue
			case product_catalog.PublicCatalogProductStatusRetired:
				continue
			}

			name := pcuServerType.Properties.Instance.OfferID
			computeServerType := computeServerTypes.Servers[name]
			serverType := &customServerType{
				Name:           name,
				HourlyPrice:    pcuServerType.Price.RetailPrice,
				MaxFileSystems: computeServerType.Capabilities.MaxFileSystems,
			}

			if availability, exists := availabilitiesResponse.Servers[name]; exists {
				serverType.Availability = availability.Availability
			}

			if pcuServerType.Properties.Hardware != nil {
				if pcuServerType.Properties.Hardware.CPU != nil {
					serverType.CPU = pcuServerType.Properties.Hardware.CPU.Virtual.Count
					serverType.Arch = pcuServerType.Properties.Hardware.CPU.Arch.String()
				}

				if pcuServerType.Properties.Hardware.Gpu != nil {
					serverType.GPU = pcuServerType.Properties.Hardware.Gpu.Count
				}

				if pcuServerType.Properties.Hardware.RAM != nil {
					serverType.RAM = pcuServerType.Properties.Hardware.RAM.Size
				}

				if pcuServerType.Properties.Hardware.Storage != nil {
					serverType.SupportedStorage = strings.Replace(
						pcuServerType.Properties.Hardware.Storage.Description,
						"Dynamic local: 1 x SSD",
						"Local",
						1,
					)
				}

				if pcuServerType.Properties.Hardware.Network != nil {
					serverType.Bandwidth = pcuServerType.Properties.Hardware.Network.MaxPublicBandwidth
				}
			}

			serverTypes = append(serverTypes, serverType)
		}

		return serverTypes, nil
	}

	return c
}

func getCompatibleTypesBuilder(c *core.Command) *core.Command {
	c.Interceptor = func(ctx context.Context, argsI any, runner core.CommandRunner) (any, error) {
		rawResp, err := runner(ctx, argsI)
		if err != nil {
			return rawResp, err
		}
		getCompatibleTypesResp := rawResp.(*instance.ServerCompatibleTypes)

		request := argsI.(*instance.GetServerCompatibleTypesRequest)
		client := core.ExtractClient(ctx)
		api := instance.NewAPI(client)
		warnings := []error(nil)

		// Get all server types to fill in the details in the response.
		listServersTypesResponse, err := api.ListServersTypes(&instance.ListServersTypesRequest{
			Zone: request.Zone,
		}, scw.WithAllPages())
		if err != nil {
			return nil, err
		}

		// Build compatible types list with details
		compatibleServerTypesCustom := []*customServerType(nil)
		for _, compatibleType := range getCompatibleTypesResp.CompatibleTypes {
			serverType, ok := listServersTypesResponse.Servers[compatibleType]
			if !ok {
				warnings = append(
					warnings,
					fmt.Errorf("could not find details on compatible type %q", compatibleType),
				)
			}

			compatibleServerTypesCustom = append(compatibleServerTypesCustom, &customServerType{
				Name: compatibleType,
				HourlyPrice: scw.NewMoneyFromFloat(
					float64(serverType.HourlyPrice),
					"EUR",
					3,
				),
				CPU:  serverType.Ncpus,
				GPU:  uint32(*serverType.Gpu),
				RAM:  scw.Size(serverType.RAM),
				Arch: serverType.Arch.String(),
			})
		}

		// Get server's current type to fill in the details in the response
		server, err := api.GetServer(&instance.GetServerRequest{
			Zone:     request.Zone,
			ServerID: request.ServerID,
		})
		if err != nil {
			return nil, err
		}
		currentServerType, ok := listServersTypesResponse.Servers[server.Server.CommercialType]
		if !ok {
			warnings = append(
				warnings,
				fmt.Errorf(
					"could not find details on current type %q",
					server.Server.CommercialType,
				),
			)
		}
		currentServerTypeCustom := []*customServerType{
			{
				Name: server.Server.CommercialType,
				HourlyPrice: scw.NewMoneyFromFloat(
					float64(currentServerType.HourlyPrice),
					"EUR",
					3,
				),
				CPU:  currentServerType.Ncpus,
				GPU:  uint32(*currentServerType.Gpu),
				RAM:  scw.Size(currentServerType.RAM),
				Arch: currentServerType.Arch.String(),
			},
		}

		return &struct {
			CurrentServerType     []*customServerType
			CompatibleServerTypes []*customServerType
		}{
			currentServerTypeCustom,
			compatibleServerTypesCustom,
		}, errors.Join(warnings...)
	}

	c.View = &core.View{
		Sections: []*core.ViewSection{
			{
				Title:     "Current Server Type",
				FieldName: "CurrentServerType",
			},
			{
				Title:     "Compatible Server Types",
				FieldName: "CompatibleServerTypes",
			},
		},
	}

	return c
}
